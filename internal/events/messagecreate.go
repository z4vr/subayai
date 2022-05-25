package events

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di"
	"github.com/sirupsen/logrus"
	"github.com/z4vr/subayai/internal/models"
	"github.com/z4vr/subayai/internal/services/database"
	"github.com/z4vr/subayai/internal/util/static"
	"github.com/z4vr/subayai/internal/util/xp"
	"math/rand"
	"strconv"
	"strings"
)

type MessageCreateEvent struct {
	db database.Database
}

func NewMessageCreateEvent(ctn di.Container) *MessageCreateEvent {
	return &MessageCreateEvent{
		db: ctn.Get(static.DiDatabase).(database.Database),
	}
}

func (m *MessageCreateEvent) HandlerXP(session *discordgo.Session, event *discordgo.MessageCreate) {

	var (
		channelID      string
		levelUpMessage string
		xpData         *models.XPData
		err            error
	)

	if event.Author.ID == session.State.User.ID {
		return
	}

	if event.Author.Bot {
		return
	}

	xpData, err = xp.GetUserXP(event.Author.ID, event.GuildID, m.db)
	if err != nil && err == database.ErrValueNotFound {
		logrus.WithFields(logrus.Fields{
			"gid": event.GuildID,
			"uid": event.Author.ID,
		}).WithError(err).Error("Failed to get user xp")
		xpData = &models.XPData{
			UserID:    event.Author.ID,
			GuildID:   event.GuildID,
			CurrentXP: 0,
			TotalXP:   0,
			Level:     0,
		}
	}

	earnedXP := rand.Intn(60) + 25
	levelup := xp.AddXP(xpData, earnedXP, false)

	err = xp.UpdateUserXP(xpData, m.db)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"gid": event.GuildID,
			"uid": event.Author.ID,
		}).WithError(err).Error("Failed to update user xp")
	}

	if levelup {
		channelID, err = m.db.GetGuildBotMessageChannelID(event.GuildID)
		if err != nil && err == database.ErrValueNotFound {
			logrus.WithFields(logrus.Fields{
				"gid": event.GuildID,
				"uid": event.Author.ID,
			}).WithError(err).Error("Failed to get bot message channel id")
			err = m.db.SetGuildBotMessageChannelID(event.GuildID, event.ChannelID)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"gid": event.GuildID,
					"uid": event.Author.ID,
				}).WithError(err).Error("Failed to set bot message channel id")
			}
			logrus.WithFields(logrus.Fields{
				"gid": event.GuildID,
				"cid": event.ChannelID,
			}).Warn("Set bot message channel id -> setup yourself")
			channelID = event.ChannelID
		} else if channelID == "" {
			channelID = event.ChannelID
		}
		levelUpMessage, err = m.db.GetGuildLevelUpMessage(event.GuildID)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"gid": event.GuildID,
				"uid": event.Author.ID,
			}).WithError(err).Error("Failed to get level up message")
			levelUpMessage = "Well done {user}, your Level of wasting time just advanced to {level}!"
		} else if levelUpMessage == "" {
			return
		}

		levelUpMessage = strings.Replace(levelUpMessage, "{user}", event.Author.Mention(), -1)
		levelUpMessage = strings.Replace(levelUpMessage, "{level}", strconv.Itoa(xpData.Level), -1)

		_, err = session.ChannelMessageSend(channelID, levelUpMessage)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"gid": event.GuildID,
				"uid": event.Author.ID,
				"cid": channelID,
			}).WithError(err).Error("Failed to send level up message")
		}

		return

	}

	return

}
