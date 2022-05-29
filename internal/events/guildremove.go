package events

import (
	"github.com/bwmarrin/discordgo"
	"github.com/z4vr/subayai/pkg/database"
)

type GuildDeleteEvent struct {
	db database.Database
}

func NewGuildDeleteEvent(db database.Database) *GuildDeleteEvent {
	return &GuildDeleteEvent{
		db: db,
	}
}

func (g *GuildDeleteEvent) Handler(s *discordgo.Session, e *discordgo.GuildDelete) {
}
