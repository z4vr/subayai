package discord

import (
	"github.com/z4vr/subayai/internal/util/math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"github.com/z4vr/subayai/internal/services/database/dberr"
)

func (d *Discord) MessageLeveling(s *discordgo.Session, e *discordgo.MessageCreate) {

	var (
		err error
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

	lastMessageTimestamp, err := d.db.GetLastMessageTimestamp(e.Author.ID, e.GuildID)
	if err != nil && err != dberr.ErrNotFound {
		logrus.WithFields(logrus.Fields{
			"guildID": e.GuildID,
		}).WithError(err).Error("Failed to get last message timestamp")
		return
	}

	if time.Now().Unix()-lastMessageTimestamp < 30 {
		return
	}

	currentLevel, err := d.db.GetUserLevel(e.Author.ID, e.GuildID)
	if err != nil && err != dberr.ErrNotFound {
		logrus.WithFields(logrus.Fields{
			"guildID": e.GuildID,
			"userID":  e.Author.ID,
		}).WithError(err).Error("Failed to get user level")
		return
	} else if err == dberr.ErrNotFound {
		currentLevel = 0
	}
	currentXP, err := d.db.GetUserCurrentXP(e.Author.ID, e.GuildID)
	if err != nil && err != dberr.ErrNotFound {
		logrus.WithFields(logrus.Fields{
			"guildID": e.GuildID,
			"userID":  e.Author.ID,
		}).WithError(err).Error("Failed to get user current xp")
		return
	} else if err == dberr.ErrNotFound {
		currentXP = 0
	}
	totalXP, err := d.db.GetUserTotalXP(e.Author.ID, e.GuildID)
	if err != nil && err != dberr.ErrNotFound {
		logrus.WithFields(logrus.Fields{
			"guildID": e.GuildID,
			"userID":  e.Author.ID,
		}).WithError(err).Error("Failed to get user total xp")
		return
	} else if err == dberr.ErrNotFound {
		totalXP = 0
	}

	earnedXP := rand.Intn(60) + 25
	totalXP += earnedXP
	currentXP, newLevel := math.CurrentLevel(earnedXP, currentLevel, currentXP)

	err = d.db.SetUserLevel(e.Author.ID, e.GuildID, newLevel)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"guildID": e.GuildID,
			"userID":  e.Author.ID,
		}).WithError(err).Error("Failed to set user level")
		return
	}
	err = d.db.SetUserCurrentXP(e.Author.ID, e.GuildID, currentXP)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"guildID": e.GuildID,
			"userID":  e.Author.ID,
		}).WithError(err).Error("Failed to set user current xp")
		return
	}
	err = d.db.SetUserTotalXP(e.Author.ID, e.GuildID, totalXP)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"guildID": e.GuildID,
			"userID":  e.Author.ID,
		}).WithError(err).Error("Failed to set user total xp")
		return
	}

	err = d.db.SetLastMessageTimestamp(e.Author.ID, e.GuildID, time.Now().Unix())
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"guildID": e.GuildID,
			"userID":  e.Author.ID,
		}).WithError(err).Error("Failed to save last message timestamp to DB")
		return
	}

	if newLevel > currentLevel {
		levelUpMessage, err := d.db.GetGuildLevelUpMessage(e.GuildID)
		if err != nil && err == dberr.ErrNotFound {
			logrus.WithError(err).Warn("Failed to get level up message")
			err = d.db.SetGuildLevelUpMessage(e.GuildID,
				"Well done {user}, your Level of wasting time just advanced to {leveling}!")
			if err != nil {
				logrus.WithError(err).Warn("Failed to set level up message")
			}
		} else if levelUpMessage == "" {
			return
		}
		botMessageChannelID, err := d.db.GetGuildBotMessageChannelID(e.GuildID)
		if err != nil && err == dberr.ErrNotFound {
			logrus.WithError(err).Warn("Failed to get bot message channel ID")
			err = d.db.SetGuildBotMessageChannelID(e.GuildID, e.ChannelID)
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
		levelUpMessage = strings.Replace(levelUpMessage, "{leveling}", strconv.Itoa(newLevel), -1)

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
