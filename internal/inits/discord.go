package inits

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di"
	"github.com/z4vr/subayai/internal/services/config"
	"github.com/z4vr/subayai/internal/util/static"
)

func NewDiscordSession(ctn di.Container) (session *discordgo.Session, err error) {
	cfg := ctn.Get(static.DiConfigProvider).(config.Provider)

	session, err = discordgo.New("Bot " + cfg.Config().Bot.Token)
	if err != nil {
		return
	}

	session.Identify.Intents = discordgo.MakeIntent(static.Intents)

	return
}
