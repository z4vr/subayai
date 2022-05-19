package events

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di"
	"github.com/z4vr/subayai/internal/services/database"
	"github.com/z4vr/subayai/internal/util/static"
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
}
