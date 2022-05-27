package events

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di"
	"github.com/sirupsen/logrus"
	"github.com/z4vr/subayai/internal/services/database"
	"github.com/z4vr/subayai/internal/services/level"
	"github.com/z4vr/subayai/internal/util/static"
	"github.com/z4vr/subayai/pkg/embedify"
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

func (m *MessageCreateEvent) HandlerXP(s *discordgo.Session, e *discordgo.MessageCreate) {

	var (
		channelID      string
		levelUpMessage string
		xpData         level.Data
		err            error
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

	xpData, err = xp2.GetUserXP(e.Author.ID, e.GuildID, m.db)
	if err != nil && err == database.ErrValueNotFound {
		logrus.WithFields(logrus.Fields{
			"gid": e.GuildID,
			"uid": e.Author.ID,
		}).WithError(err).Error("Failed to get user level")
		xpData = &xp2.XPData{
			UserID:    e.Author.ID,
			GuildID:   e.GuildID,
			CurrentXP: 0,
			TotalXP:   0,
			Level:     0,
		}
	}

	earnedXP := rand.Intn(60) + 25
	levelup := xpData.AddXP(earnedXP, false)

	err = xp2.UpdateUserXP(xpData, m.db)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"gid": e.GuildID,
			"uid": e.Author.ID,
		}).WithError(err).Error("Failed to update user level")
	}

	if levelup {
		channelID, err = m.db.GetGuildBotMessageChannelID(e.GuildID)
		if err != nil && err == database.ErrValueNotFound {
			logrus.WithFields(logrus.Fields{
				"gid": e.GuildID,
				"uid": e.Author.ID,
			}).WithError(err).Error("Failed to get bot message channel id")
			err = m.db.SetGuildBotMessageChannelID(e.GuildID, e.ChannelID)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"gid": e.GuildID,
					"uid": e.Author.ID,
				}).WithError(err).Error("Failed to set bot message channel id")
			}
			logrus.WithFields(logrus.Fields{
				"gid": e.GuildID,
				"cid": e.ChannelID,
			}).Warn("Set bot message channel id -> setup yourself")
			channelID = e.ChannelID
		} else if channelID == "" {
			channelID = e.ChannelID
		}
		levelUpMessage, err = m.db.GetGuildLevelUpMessage(e.GuildID)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"gid": e.GuildID,
				"uid": e.Author.ID,
			}).WithError(err).Error("Failed to get level up message")
			levelUpMessage = "Well done {user}, your Level of wasting time just advanced to {level}!"
		} else if levelUpMessage == "" {
			return
		}

		levelUpMessage = strings.Replace(levelUpMessage, "{user}", e.Author.Mention(), -1)
		levelUpMessage = strings.Replace(levelUpMessage, "{level}", strconv.Itoa(xpData.Level), -1)

		emb := embedify.New().
			SetAuthor(e.Author.Username, static.AppIcon).
			SetDescription(levelUpMessage).Build()

		_, err = s.ChannelMessageSendEmbed(channelID, emb)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"gid": e.GuildID,
				"uid": e.Author.ID,
				"cid": channelID,
			}).WithError(err).Error("Failed to send level up message")
		}

		return

	}

	return

}
