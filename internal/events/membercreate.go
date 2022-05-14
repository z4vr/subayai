package events

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di"
	"github.com/sirupsen/logrus"
	"github.com/z4vr/subayai/internal/database"
	"github.com/z4vr/subayai/internal/util/static"
)

type GuildMemberAddEvent struct {
	db database.Database
}

func NewGuildMemberAddEvent(ctn di.Container) *GuildMemberAddEvent {
	return &GuildMemberAddEvent{
		db: ctn.Get(static.DiDatabase).(database.Database),
	}
}

func (g *GuildMemberAddEvent) Handler(session *discordgo.Session, event *discordgo.GuildMemberAdd) {
	exists, err := g.db.UserEntryExists(event.User.ID)
	if err != nil {
		logrus.WithError(err).Error("Failed checking if user exists")
	} else if !exists {
		err = g.db.CreateUserEntry(event.User.ID)
		if err != nil {
			logrus.WithError(err).Error("Failed to create user entry")
		}
	}
	exists, err = g.db.UserXPEntryExists(event.User.ID, event.Member.GuildID)
	if err != nil {
		logrus.WithError(err).Error("Failed checking if user xp exists")
	} else if !exists {
		err = g.db.CreateUserXPEntry(event.User.ID, event.Member.GuildID)
		if err != nil {
			logrus.WithError(err).Error("Failed to create user xp entry")
		}
	}
}
