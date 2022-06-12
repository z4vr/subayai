package discord

import (
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"github.com/z4vr/subayai/internal/services/database/dberr"
)

func (h *EventHandler) VoiceLeveling(s *discordgo.Session, e *discordgo.VoiceStateUpdate) {
	// check if the member is a bot
	member, err := h.d.GetMember(e.VoiceState.UserID, e.VoiceState.GuildID)
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

	afkChannelID, err = h.db.GetGuildAFKChannelID(e.VoiceState.GuildID)
	if err != nil && err != dberr.ErrNotFound {
		logrus.WithFields(logrus.Fields{
			"guildID": e.VoiceState.GuildID,
		}).WithError(err).Error("Failed to get afk channel id")
	}
	lastSessionID, err = h.db.GetLastVoiceSessionID(member.User.ID, e.VoiceState.GuildID)
	if err != nil && err != dberr.ErrNotFound {
		logrus.WithFields(logrus.Fields{
			"guildID": e.VoiceState.GuildID,
			"userID":  member.User.ID,
		}).WithError(err).Error("Failed to get last session id")
	}
	lastSessionTimestamp, err = h.db.GetLastVoiceSessionTimestamp(member.User.ID, e.VoiceState.GuildID)
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
		err = h.db.SetLastVoiceSessionID(member.User.ID, e.VoiceState.GuildID, lastSessionID)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"guildID": e.VoiceState.GuildID,
				"userID":  member.User.ID,
			}).WithError(err).Error("Failed to set last session id")
		}
		err = h.db.SetLastVoiceSessionTimestamp(member.User.ID, e.VoiceState.GuildID, lastSessionTimestamp)
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
		levelData, err := h.lp.FetchFromDB(member.User.ID, e.VoiceState.GuildID)
		if err != nil && err == dberr.ErrNotFound {
			err = h.lp.SaveToDB(levelData)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"guildID": e.VoiceState.GuildID,
					"userID":  member.User.ID,
				}).WithError(err).Error("Failed to save level data")
			}
		}

		// reward the user 1 xp per minute spent in voice
		earnedXP := int(nowTimestamp-lastSessionTimestamp) / 60
		levelup := levelData.LevelUp(earnedXP, false)

		err = h.lp.SaveToDB(levelData)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"guildID": e.VoiceState.GuildID,
				"userID":  member.User.ID,
			}).WithError(err).Error("Failed to update user leveling")
		}

		if levelup {
			levelUpMessage, err := h.db.GetGuildLevelUpMessage(e.GuildID)
			if err != nil && err == dberr.ErrNotFound {
				logrus.WithError(err).Warn("Failed to get level up message")
				err = h.db.SetGuildLevelUpMessage(e.GuildID,
					"Well done {user}, your Level of wasting time just advanced to {leveling}!")
				if err != nil {
					logrus.WithError(err).Error("Failed to set level up message")
				}
			} else if levelUpMessage == "" {
				return
			}
			botMessageChannelID, err := h.db.GetGuildBotMessageChannelID(e.GuildID)
			if err != nil && err == dberr.ErrNotFound {
				logrus.WithFields(logrus.Fields{
					"guildID": e.GuildID,
				}).WithError(err).Error("Failed to get bot message channel id")

			} else if botMessageChannelID == "" {
				botMessageChannel := h.d.FindGuildTextChannel(e.GuildID)
				if botMessageChannel == nil {
					logrus.WithFields(logrus.Fields{
						"guildID": e.GuildID,
					}).Error("No bot message channel found")
					return
				}
				botMessageChannelID = botMessageChannel.ID
			}

			levelUpMessage = strings.Replace(levelUpMessage, "{user}", member.Mention(), -1)
			levelUpMessage = strings.Replace(levelUpMessage, "{leveling}", strconv.Itoa(levelData.Level), -1)

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
