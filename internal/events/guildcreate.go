package events

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di"
	"github.com/sirupsen/logrus"
	"github.com/z4vr/subayai/internal/database"
	"github.com/z4vr/subayai/internal/util/static"
)

type GuildCreateEvent struct {
	db database.Database
}

func NewGuildCreateEvent(ctn di.Container) *GuildCreateEvent {
	return &GuildCreateEvent{
		db: ctn.Get(static.DiDatabase).(database.Database),
	}
}

func (g *GuildCreateEvent) Handler(session *discordgo.Session, event *discordgo.GuildCreate) {
	exists, err := g.db.GuildEntryExists(event.Guild.ID)
	if err != nil {
		logrus.WithError(err).Error("Failed checking if guild exists")
	} else if !exists {
		err = g.db.CreateGuildEntry(event.Guild.ID)
		if err != nil {
			logrus.WithError(err).Error("Failed to create guild entry")
		}
	}
}
