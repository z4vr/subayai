package discord

import (
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"github.com/z4vr/subayai/pkg/database/dberr"
	"github.com/z4vr/subayai/pkg/leveling"
)

func (d *Discord) VoiceUpdate(s *discordgo.Session, e *discordgo.VoiceStateUpdate) {
	// check if the member is a bot
	member, err := d.GetMember(e.GuildID, e.UserID)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"guildID": e.GuildID,
				"userID":  e.UserID,
			}).
			Error("Failed to get member")
		return
	}

	if member.User.Bot {
		return
	}

	var (
		afkChannelID         string = "empty"
		lastSessionID        string = ""
		lastSessionTimestamp int64  = 0
		nowTimestamp         int64  = time.Now().Unix()
	)

	afkChannelID, err = d.db.GetGuildAFKChannelID(e.GuildID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"guildID": e.GuildID,
		}).WithError(err).Error("Failed to get afk channel id")
	}
	lastSessionID, err = d.db.GetLastVoiceSessionID(e.GuildID, e.UserID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"guildID": e.GuildID,
			"userID":  e.UserID,
		}).WithError(err).Error("Failed to get last session id")
	}
	lastSessionTimestamp, err = d.db.GetLastVoiceSessionTimestamp(e.GuildID, e.UserID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"guildID": e.GuildID,
			"userID":  e.UserID,
		}).WithError(err).Error("Failed to get last session timestamp")
	}

	// scenario: user freshly joined voice channel, or rejoined from afk channel
	if e.BeforeUpdate == nil || e.BeforeUpdate.ChannelID == afkChannelID {
		// update our last records
		lastSessionID = e.SessionID
		lastSessionTimestamp = nowTimestamp

		// save them to db
		err = d.db.SetLastVoiceSessionID(e.GuildID, e.UserID, lastSessionID)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"guildID": e.GuildID,
				"userID":  e.UserID,
			}).WithError(err).Error("Failed to set last session id")
		}
		err = d.db.SetLastVoiceSessionTimestamp(e.GuildID, e.UserID, lastSessionTimestamp)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"guildID": e.GuildID,
				"userID":  e.UserID,
			}).WithError(err).Error("Failed to set last session timestamp")
		}

	}

	// scenario: user left voice channel
	if e.VoiceState.ChannelID == "" && e.VoiceState.SessionID == lastSessionID ||
		e.VoiceState.ChannelID == afkChannelID && e.VoiceState.SessionID == lastSessionID {
		// reward the leveling for the time spent in voice
		levelData, err := d.lp.FetchFromDB(e.GuildID, e.UserID)
		if err != nil && err == dberr.ErrNotFound {
			levelData := &leveling.LevelData{
				UserID:    e.UserID,
				GuildID:   e.GuildID,
				Level:     0,
				CurrentXP: 0,
				TotalXP:   0,
			}
			d.lp.SaveToDB(levelData)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"userID":  e.UserID,
					"guildID": e.GuildID,
				}).WithError(err).Error("Failed to add user to level map")
			}
		}

		// reward the user 1 xp per minute spent in voice
		earnedXP := int(nowTimestamp-lastSessionTimestamp) / 60
		levelup := levelData.LevelUp(earnedXP, false)

		err = d.lp.SaveToDB(levelData)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"guildID": e.GuildID,
				"userID":  member.User.ID,
			}).WithError(err).Error("Failed to update user leveling")
		}

		if levelup {
			levelUpMessage, err := d.db.GetGuildLevelUpMessage(e.GuildID)
			if err != nil && err == dberr.ErrNotFound {
				logrus.WithError(err).Warn("Failed to get level up message")
				err = d.db.SetGuildLevelUpMessage(e.GuildID,
					"Well done {user}, your Level of wasting time just advanced to {leveling}!")
				if err != nil {
					logrus.WithError(err).Warn("Failed to set level up message")
				}
			} else if levelUpMessage == "" {
				levelUpMessage = "Well done {user}, your Level of wasting time just advanced to {leveling}!"
			}
			botMessageChannelID, err := d.db.GetGuildBotMessageChannelID(e.GuildID)
			if err != nil {
				return
			} else if botMessageChannelID == "" {
				botMessageChannelID = e.ChannelID
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
