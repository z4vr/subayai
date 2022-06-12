package discord

import (
	"github.com/z4vr/subayai/internal/util/math"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"github.com/z4vr/subayai/internal/services/database/dberr"
)

func (d *Discord) VoiceLeveling(s *discordgo.Session, e *discordgo.VoiceStateUpdate) {
	// check if the member is a bot
	member, err := d.GetMember(e.VoiceState.UserID, e.VoiceState.GuildID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"guildID": e.VoiceState.GuildID,
			"userID":  e.VoiceState.UserID,
		}).WithError(err).Error("Failed to get guild member")
	}

	if member.User.Bot {
		return
	}

	var (
		afkChannelID               = "empty"
		lastSessionID              = ""
		lastSessionTimestamp int64 = 0
		nowTimestamp               = time.Now().Unix()
	)

	afkChannelID, err = d.db.GetGuildAFKChannelID(e.VoiceState.GuildID)
	if err != nil && err != dberr.ErrNotFound {
		logrus.WithFields(logrus.Fields{
			"guildID": e.VoiceState.GuildID,
		}).WithError(err).Error("Failed to get afk channel id")
	}
	lastSessionID, err = d.db.GetLastVoiceSessionID(member.User.ID, e.VoiceState.GuildID)
	if err != nil && err != dberr.ErrNotFound {
		logrus.WithFields(logrus.Fields{
			"guildID": e.VoiceState.GuildID,
			"userID":  member.User.ID,
		}).WithError(err).Error("Failed to get last session id")
	}
	lastSessionTimestamp, err = d.db.GetLastVoiceSessionTimestamp(member.User.ID, e.VoiceState.GuildID)
	if err != nil && err != dberr.ErrNotFound {
		logrus.WithFields(logrus.Fields{
			"guildID": e.VoiceState.GuildID,
			"userID":  member.User.ID,
		}).WithError(err).Error("Failed to get last session timestamp")
	}

	// scenario: user freshly joined voice channel, or rejoined from afk channel
	if e.BeforeUpdate == nil || e.BeforeUpdate.ChannelID == afkChannelID {
		// update our last records
		lastSessionID = e.VoiceState.SessionID
		lastSessionTimestamp = nowTimestamp

		// save them to db
		err = d.db.SetLastVoiceSessionID(member.User.ID, e.VoiceState.GuildID, lastSessionID)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"guildID": e.VoiceState.GuildID,
				"userID":  member.User.ID,
			}).WithError(err).Error("Failed to set last session id")
		}
		err = d.db.SetLastVoiceSessionTimestamp(member.User.ID, e.VoiceState.GuildID, lastSessionTimestamp)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"guildID": e.VoiceState.GuildID,
				"userID":  member.User.ID,
			}).WithError(err).Error("Failed to set last session timestamp")
		}

	}

	// scenario: user left voice channel
	if e.VoiceState.ChannelID == "" && e.VoiceState.SessionID == lastSessionID ||
		e.VoiceState.ChannelID == afkChannelID && e.VoiceState.SessionID == lastSessionID {
		// reward the leveling for the time spent in voice
		currentLevel, err := d.db.GetUserLevel(member.User.ID, e.GuildID)
		if err != nil && err != dberr.ErrNotFound {
			logrus.WithFields(logrus.Fields{
				"guildID": e.GuildID,
				"userID":  member.User.ID,
			}).WithError(err).Error("Failed to get user level")
			return
		} else if err == dberr.ErrNotFound {
			currentLevel = 0
		}
		currentXP, err := d.db.GetUserCurrentXP(member.User.ID, e.GuildID)
		if err != nil && err != dberr.ErrNotFound {
			logrus.WithFields(logrus.Fields{
				"guildID": e.GuildID,
				"userID":  member.User.ID,
			}).WithError(err).Error("Failed to get user current xp")
			return
		} else if err == dberr.ErrNotFound {
			currentXP = 0
		}
		totalXP, err := d.db.GetUserTotalXP(member.User.ID, e.GuildID)
		if err != nil && err != dberr.ErrNotFound {
			logrus.WithFields(logrus.Fields{
				"guildID": e.GuildID,
				"userID":  member.User.ID,
			}).WithError(err).Error("Failed to get user total xp")
			return
		} else if err == dberr.ErrNotFound {
			totalXP = 0
		}

		levelMap := map[string]int{
			"level":     currentLevel,
			"currentXP": currentXP,
			"totalXP":   totalXP,
		}

		// reward the user 1 xp per minute spent in voice
		earnedXP := int(nowTimestamp-lastSessionTimestamp) / 60
		levelMap["currentXP"] += earnedXP
		levelMap["totalXP"] += earnedXP
		newLevel := math.CurrentLevel(levelMap["currentXP"], levelMap["level"])

		if newLevel > levelMap["level"] {
			levelUpMessage, err := d.db.GetGuildLevelUpMessage(e.GuildID)
			if err != nil && err == dberr.ErrNotFound {
				logrus.WithError(err).Warn("Failed to get level up message")
				err = d.db.SetGuildLevelUpMessage(e.GuildID,
					"Well done {user}, your Level of wasting time just advanced to {leveling}!")
				if err != nil {
					logrus.WithError(err).Error("Failed to set level up message")
				}
			} else if levelUpMessage == "" {
				return
			}

			levelUpMessage = strings.Replace(levelUpMessage, "{user}", member.Mention(), -1)
			levelUpMessage = strings.Replace(levelUpMessage, "{leveling}", strconv.Itoa(newLevel), -1)

			botMessageChannelID, err := d.db.GetGuildBotMessageChannelID(e.GuildID)
			if err != nil && err == dberr.ErrNotFound {
				logrus.WithFields(logrus.Fields{
					"guildID": e.GuildID,
				}).WithError(err).Error("Failed to get bot message channel id")

			} else if botMessageChannelID == "" {
				botMessageChannel := d.FindGuildTextChannel(e.GuildID)
				if botMessageChannel.ID == "" {
					_, err := d.SendMessageDM(member.User.ID, levelUpMessage)
					if err != nil {
						logrus.WithFields(logrus.Fields{
							"guildID": e.GuildID,
							"userID":  member.User.ID,
						}).WithError(err).Error("Failed to send level up message")
						return
					}
				} else {
					botMessageChannelID = botMessageChannel.ID
					_, err = s.ChannelMessageSend(botMessageChannelID, levelUpMessage)
					if err != nil {
						logrus.WithFields(logrus.Fields{
							"guildID":   e.GuildID,
							"userID":    member.User.ID,
							"channelID": botMessageChannelID,
						}).WithError(err).Error("Failed to send level up message")
						return
					}
				}

			}
		}
	}
}
