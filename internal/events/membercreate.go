package events

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di"
	"github.com/z4vr/subayai/internal/services/database"
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
}
