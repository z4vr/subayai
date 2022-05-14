package events

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di"
	"github.com/sirupsen/logrus"
	"github.com/z4vr/subayai/internal/database"
	"github.com/z4vr/subayai/internal/util/static"
)

type GuildDeleteEvent struct {
	db database.Database
}

func NewGuildDeleteEvent(ctn di.Container) *GuildDeleteEvent {
	return &GuildDeleteEvent{
		db: ctn.Get(static.DiDatabase).(database.Database),
	}
}

func (g *GuildDeleteEvent) Handler(session *discordgo.Session, event *discordgo.GuildDelete) {
	exists, err := g.db.GuildEntryExists(event.Guild.ID)
	if err != nil {
		logrus.WithError(err).Error("Failed checking if guild exists")
	} else if !exists {
		logrus.WithField("guildID", event.Guild.ID).Info("Guild does not exist")
	}
	ok, err := g.db.GetGuildAutoDelete(event.Guild.ID)
	if err != nil {
		logrus.WithError(err).Error("Failed checking if guild exists")
	} else if ok {
		err = g.db.DeleteGuildEntry(event.Guild.ID)
		if err != nil {
			logrus.WithError(err).Error("Failed deleting guild from auto delete list")
		}
	}
}
