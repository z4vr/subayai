package events

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di"
	"github.com/z4vr/subayai/internal/services/database"
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
}
