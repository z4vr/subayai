package events

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di"
	"github.com/sirupsen/logrus"
	"github.com/z4vr/subayai/internal/services/database"
	"github.com/z4vr/subayai/internal/util/static"
	"github.com/z4vr/subayai/internal/util/xp"
	"math/rand"
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
	if event.Author.ID == session.State.User.ID {
		return
	}

	if event.Author.Bot {
		return
	}

	xpData, err := xp.GetUserXP(event.Author.ID, event.GuildID, m.db)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"gid": event.GuildID,
			"uid": event.Author.ID,
		}).WithError(err).Error("Failed to get user xp")
		xpData, err = xp.GenerateUserXP(event.Author.ID, event.GuildID, m.db)
	}

	// generate random number between 25 and 85
	earnedXP := rand.Intn(60) + 25
	xpData.AddXP(earnedXP)

	err = xp.UpdateUserXP(xpData, m.db)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"gid": event.GuildID,
			"uid": event.Author.ID,
		}).WithError(err).Error("Failed to update user xp")
	}

	return

}
