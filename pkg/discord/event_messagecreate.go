package discord

import (
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"github.com/z4vr/subayai/pkg/database/dberr"
	"github.com/z4vr/subayai/pkg/leveling"
)

func (d *Discord) MessageCreate(s *discordgo.Session, e *discordgo.MessageCreate) {

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

	lastMessageTimestamp, err := d.db.GetLastMessageTimestamp(e.GuildID, e.Author.ID)
	if err != nil && err != dberr.ErrNotFound {
		logrus.WithError(err).Error("Failed to get last message timestamp")
		return
	}

	if time.Now().Unix()-lastMessageTimestamp < 30 {
		return
	}

	levelData, err = d.lp.FetchFromDB(e.GuildID, e.Author.ID)
	if err != nil && err == dberr.ErrNotFound {
		levelData = &leveling.LevelData{
			UserID:    e.Author.ID,
			GuildID:   e.GuildID,
			Level:     0,
			CurrentXP: 0,
			TotalXP:   0,
		}
		d.lp.SaveToDB(levelData)
	}

	earnedXP := rand.Intn(60) + 25
	levelup := levelData.LevelUp(earnedXP, false)

	err = d.lp.SaveToDB(levelData)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"guildID": levelData.GuildID,
			"userID":  levelData.UserID,
		}).WithError(err).Error("Failed to save level data to DB")
		return
	}

	// If the user leveled up, we need to send a message to the channel
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
