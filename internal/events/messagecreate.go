package events

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di"
	"github.com/sirupsen/logrus"
	"github.com/z4vr/subayai/internal/database"
	"github.com/z4vr/subayai/internal/util/static"
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

func (m *MessageCreateEvent) Handler(session *discordgo.Session, event *discordgo.MessageCreate) {
	exists, err := m.db.UserXPEntryExists(event.Author.ID, event.GuildID)
	if err != nil {
		logrus.WithError(err).Error("Failed checking if user xp exists")
	} else if !exists {
		err = m.db.CreateUserXPEntry(event.Author.ID, event.GuildID)
		if err != nil {
			logrus.WithError(err).Error("Failed to create user xp entry")
		}
	}

	// Leveling starts here

	// Get the user's xp
	xp, err := m.db.GetUserXPEntry(event.Author.ID, event.GuildID)
	if err != nil {
		logrus.WithError(err).Error("Failed to get user xp entry")
		return
	}

	xpEarned := rand.Intn(75) + 25
	xp.LevelUp(xpEarned)

	// Save the user's new xp
	err = m.db.SetUserXPEntry(event.Author.ID, event.GuildID, xp)

}
