package inits

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di"
	"github.com/z4vr/subayai/internal/events"
	"github.com/z4vr/subayai/internal/services/config"
	static2 "github.com/z4vr/subayai/internal/util/static"
)

func NewDiscordSession(ctn di.Container) (session *discordgo.Session, err error) {
	cfg := ctn.Get(static2.DiConfigProvider).(config.Provider)

	session, err = discordgo.New("Bot " + cfg.Instance().Discord.Token)
	if err != nil {
		return
	}

	session.Identify.Intents = discordgo.MakeIntent(static2.Intents)

	session.AddHandler(events.NewReadyEvent().Handler)
	session.AddHandler(events.NewMessageCreateEvent(ctn).Handler)
	session.AddHandler(events.NewGuildCreateEvent(ctn).Handler)
	session.AddHandler(events.NewGuildDeleteEvent(ctn).Handler)
	session.AddHandler(events.NewGuildMemberAddEvent(ctn).Handler)

	return
}
