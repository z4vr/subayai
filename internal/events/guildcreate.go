package events

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di"
	"github.com/z4vr/subayai/internal/services/database"
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
}
