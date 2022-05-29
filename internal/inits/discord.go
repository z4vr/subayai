package inits

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di"
	"github.com/z4vr/subayai/internal/events"
	"github.com/z4vr/subayai/internal/services/config"
	static2 "github.com/z4vr/subayai/internal/static"
)

func NewDiscordSession(ctn di.Container) (session *discordgo.Session, err error) {
	cfg := ctn.Get(static2.DiConfigProvider).(config.Provider)

	session, err = discordgo.New("Bot " + cfg.Config().Bot.Token)
	if err != nil {
		return
	}

	session.Identify.Intents = discordgo.MakeIntent(static2.Intents)

	// Register handlers
	// Ready handlers
	session.AddHandler(events.NewReadyEvent().Handler)
	// Message handlers
	session.AddHandler(events.NewMessageCreateEvent(ctn).HandlerXP)
	// Guild create handlers
	//session.AddHandler(events.NewGuildCreateEvent(ctn).HandlerCreate)
	// Guild delete handlers
	//session.AddHandler(events.NewGuildDeleteEvent(ctn).Handler)
	// Guild member add handlers
	//session.AddHandler(events.NewGuildMemberAddEvent(ctn).HandlerAutoRole)

	return
}
