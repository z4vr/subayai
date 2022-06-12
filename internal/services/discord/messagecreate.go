package discord

import (
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"github.com/z4vr/subayai/internal/services/database/dberr"
	"github.com/z4vr/subayai/internal/services/leveling"
)

func (h *EventHandler) MessageLeveling(s *discordgo.Session, e *discordgo.MessageCreate) {

	var (
		levelData *leveling.LevelData
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

	lastMessageTimestamp, err := h.db.GetLastMessageTimestamp(e.Author.ID, e.GuildID)
	if err != nil && err != dberr.ErrNotFound {
		logrus.WithFields(logrus.Fields{
			"guildID": e.GuildID,
		}).WithError(err).Error("Failed to get last message timestamp")
		return
	}

	if time.Now().Unix()-lastMessageTimestamp < 30 {
		return
	}

	levelData, err = h.lp.FetchFromDB(e.Author.ID, e.GuildID)
	if err != nil && err == dberr.ErrNotFound {
		err := h.lp.SaveToDB(levelData)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"guildID": e.GuildID,
				"userID":  e.Author.ID,
			}).WithError(err).Error("Failed to save level data")
		}
	}

	earnedXP := rand.Intn(60) + 25
	levelup := levelData.LevelUp(earnedXP, false)

	err = h.lp.SaveToDB(levelData)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"guildID": levelData.GuildID,
			"userID":  levelData.UserID,
		}).WithError(err).Error("Failed to save level data to DB")
		return
	}

	err = h.db.SetLastMessageTimestamp(e.Author.ID, e.GuildID, time.Now().Unix())
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"guildID": levelData.GuildID,
			"userID":  levelData.UserID,
		}).WithError(err).Error("Failed to save last message timestamp to DB")
		return
	}

	// If the user leveled up, we need to send a message to the channel
	if levelup {
		levelUpMessage, err := h.db.GetGuildLevelUpMessage(e.GuildID)
		if err != nil && err == dberr.ErrNotFound {
			logrus.WithError(err).Warn("Failed to get level up message")
			err = h.db.SetGuildLevelUpMessage(e.GuildID,
				"Well done {user}, your Level of wasting time just advanced to {leveling}!")
			if err != nil {
				logrus.WithError(err).Warn("Failed to set level up message")
			}
		} else if levelUpMessage == "" {
			return
		}
		botMessageChannelID, err := h.db.GetGuildBotMessageChannelID(e.GuildID)
		if err != nil && err == dberr.ErrNotFound {
			logrus.WithError(err).Warn("Failed to get bot message channel ID")
			err = h.db.SetGuildBotMessageChannelID(e.GuildID, e.ChannelID)
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
