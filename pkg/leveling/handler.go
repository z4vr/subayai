package leveling

import (
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/z4vr/subayai/pkg/database/dberr"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

func (p *Provider) MessageCreate(s *discordgo.Session, e *discordgo.MessageCreate) {

	var (
		levelData *LevelData
		err       error
	)

	if e.Author.ID == s.State.User.ID {
		return
	}

	if e.Author.Bot {
		return
	}

	if e.GuildID == "" {
		return
	}

	lastMessageTimestamp, err := p.db.GetLastMessageTimestamp(e.GuildID, e.Author.ID)
	if err != nil && err != dberr.ErrNotFound {
		logrus.WithError(err).Error("Failed to get last message timestamp")
		return
	}

	if time.Now().Unix()-lastMessageTimestamp < 30 {
		return
	}

	levelData, err = p.Get(e.Author.ID, e.GuildID)
	if err != nil && err == dberr.ErrNotFound {
		levelData = &LevelData{
			UserID:    e.Author.ID,
			GuildID:   e.GuildID,
			Level:     0,
			CurrentXP: 0,
			TotalXP:   0,
		}
		p.Set(e.Author.ID, e.GuildID, levelData)
		p.SaveToDB(levelData)
	}

	earnedXP := rand.Intn(60) + 25
	levelup := levelData.AddXP(earnedXP, false)

	err = p.SaveToDB(levelData)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"guildID": levelData.GuildID,
			"userID":  levelData.UserID,
		}).WithError(err).Error("Failed to save level data to DB")
		return
	}

	// If the user leveled up, we need to send a message to the channel
	if levelup {
		levelUpMessage, err := p.db.GetGuildLevelUpMessage(e.GuildID)
		if err != nil && err == dberr.ErrNotFound {
			logrus.WithError(err).Warn("Failed to get level up message")
			err = p.db.SetGuildLevelUpMessage(e.GuildID,
				"Well done {user}, your Level of wasting time just advanced to {leveling}!")
			if err != nil {
				logrus.WithError(err).Warn("Failed to set level up message")
			}
		} else if levelUpMessage == "" {
			levelUpMessage = "Well done {user}, your Level of wasting time just advanced to {leveling}!"
		}
		botMessageChannelID, err := p.db.GetGuildBotMessageChannelID(e.GuildID)
		if err != nil && err == dberr.ErrNotFound {
			logrus.WithError(err).Warn("Failed to get bot message channel ID")
			err = p.db.SetGuildBotMessageChannelID(e.GuildID, e.ChannelID)
			if err != nil {
				logrus.WithError(err).Error("Failed to set bot message channel ID")
			}
			logrus.WithFields(
				logrus.Fields{
					"guildID": e.GuildID,
				}).Info("Set bot message channel ID to ", e.ChannelID)
			return
		} else if botMessageChannelID == "" {
			botMessageChannelID = e.ChannelID
		}

		levelUpMessage = strings.Replace(levelUpMessage, "{user}", e.Author.Mention(), -1)
		levelUpMessage = strings.Replace(levelUpMessage, "{leveling}", strconv.Itoa(levelData.Level), -1)

		_, err = s.ChannelMessageSend(botMessageChannelID, levelUpMessage)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"guildID":   e.GuildID,
				"userID":    e.Author.ID,
				"channelID": botMessageChannelID,
			}).WithError(err).Error("Failed to send level up message")
			return
		}
	}

}

func (p *Provider) VoiceUpdate(s *discordgo.Session, e *discordgo.VoiceStateUpdate) {
	// check if the member is a bot
	member, err := p.dc.GetMember(e.GuildID, e.UserID)
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

	afkChannelID, err = p.db.GetGuildAFKChannelID(e.GuildID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"guildID": e.GuildID,
		}).WithError(err).Error("Failed to get afk channel id")
	}
	lastSessionID, err = p.db.GetLastVoiceSessionID(e.GuildID, e.UserID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"guildID": e.GuildID,
			"userID":  e.UserID,
		}).WithError(err).Error("Failed to get last session id")
	}
	lastSessionTimestamp, err = p.db.GetLastVoiceSessionTimestamp(e.GuildID, e.UserID)
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
		err = p.db.SetLastVoiceSessionID(e.GuildID, e.UserID, lastSessionID)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"guildID": e.GuildID,
				"userID":  e.UserID,
			}).WithError(err).Error("Failed to set last session id")
		}
		err = p.db.SetLastVoiceSessionTimestamp(e.GuildID, e.UserID, lastSessionTimestamp)
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
		levelData, err := p.Get(e.UserID, e.GuildID)
		if err != nil && err == dberr.ErrNotFound {
			levelData := &LevelData{
				UserID:    e.UserID,
				GuildID:   e.GuildID,
				Level:     0,
				CurrentXP: 0,
				TotalXP:   0,
			}
			p.Set(e.UserID, e.GuildID, levelData)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"userID":  e.UserID,
					"guildID": e.GuildID,
				}).WithError(err).Error("Failed to add user to level map")
			}
		}

		// reward the user 1 xp per minute spent in voice
		earnedXP := int(nowTimestamp-lastSessionTimestamp) / 60
		levelup := levelData.AddXP(earnedXP, false)

		err = p.SaveToDB(levelData)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"guildID": e.GuildID,
				"userID":  member.User.ID,
			}).WithError(err).Error("Failed to update user leveling")
		}

		if levelup {
			levelUpMessage, err := p.db.GetGuildLevelUpMessage(e.GuildID)
			if err != nil && err == dberr.ErrNotFound {
				logrus.WithError(err).Warn("Failed to get level up message")
				err = p.db.SetGuildLevelUpMessage(e.GuildID,
					"Well done {user}, your Level of wasting time just advanced to {leveling}!")
				if err != nil {
					logrus.WithError(err).Warn("Failed to set level up message")
				}
			} else if levelUpMessage == "" {
				levelUpMessage = "Well done {user}, your Level of wasting time just advanced to {leveling}!"
			}
			botMessageChannelID, err := p.db.GetGuildBotMessageChannelID(e.GuildID)
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
