package events

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di"
	"github.com/z4vr/subayai/internal/static"
	"github.com/z4vr/subayai/pkg/database"
)

type GuildDeleteEvent struct {
	db database.Database
}

func NewGuildDeleteEvent(ctn di.Container) *GuildDeleteEvent {
	return &GuildDeleteEvent{
		db: ctn.Get(static.DiDatabase).(database.Database),
	}
}

func (g *GuildDeleteEvent) Handler(s *discordgo.Session, e *discordgo.GuildDelete) {
}
