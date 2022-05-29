package events

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"github.com/z4vr/subayai/internal/static"
	"github.com/z4vr/subayai/pkg/database"
	"github.com/z4vr/subayai/pkg/database/dberr"
	"github.com/z4vr/subayai/pkg/embedify"
	"github.com/z4vr/subayai/pkg/leveling"
	"math/rand"
	"strconv"
	"strings"
)

type MessageCreateEvent struct {
	db database.Database
	lp *leveling.Provider
}

func NewMessageCreateEvent(db database.Database, lp *leveling.Provider) *MessageCreateEvent {
	return &MessageCreateEvent{
		db: db,
		lp: lp,
	}
}

func (m *MessageCreateEvent) HandlerXP(s *discordgo.Session, e *discordgo.MessageCreate) {

	var (
		channelID      string
		levelUpMessage string
		xpData         leveling.Data
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

	levelData := m.lp.Get(e.Author.ID, e.GuildID)
	if levelData == nil {
		levelData := &leveling.Data{
			UserID:    e.Author.ID,
			GuildID:   e.GuildID,
			Level:     0,
			CurrentXP: 0,
			TotalXP:   0,
		}
		err := m.lp.Add(levelData)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"uid": e.Author.ID,
				"gid": e.GuildID,
			}).WithError(err).Error("Failed to add user to level map")
		}
	}

	earnedXP := rand.Intn(60) + 25
	levelup := xpData.AddXP(earnedXP, false)

	err = m.lp.SaveLevelData(xpData)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"gid": e.GuildID,
			"uid": e.Author.ID,
		}).WithError(err).Error("Failed to update user leveling")
	}

	if levelup {
		channelID, err = m.db.GetGuildBotMessageChannelID(e.GuildID)
		if err != nil && err == dberr.ErrNotFound {
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
			}).WithError(err).Error("Failed to get leveling up message")
			levelUpMessage = "Well done {user}, your Level of wasting time just advanced to {leveling}!"
		} else if levelUpMessage == "" {
			return
		}

		levelUpMessage = strings.Replace(levelUpMessage, "{user}", e.Author.Mention(), -1)
		levelUpMessage = strings.Replace(levelUpMessage, "{leveling}", strconv.Itoa(xpData.Level), -1)

		emb := embedify.New().
			SetAuthor(e.Author.Username, static.AppIcon).
			SetDescription(levelUpMessage).Build()

		_, err = s.ChannelMessageSendEmbed(channelID, emb)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"gid": e.GuildID,
				"uid": e.Author.ID,
				"cid": channelID,
			}).WithError(err).Error("Failed to send leveling up message")
		}

		return

	}

	return

}
