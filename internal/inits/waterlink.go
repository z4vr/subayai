package inits

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di"
	"github.com/sirupsen/logrus"
	"github.com/z4vr/subayai/internal/services/config"
	"github.com/z4vr/subayai/internal/services/waterlink"
	"github.com/z4vr/subayai/internal/util/static"
)

func NewWaterlinkProvider(ctn di.Container) *waterlink.WaterlinkProvider {

	cfg := ctn.Get(static.DiConfigProvider).(config.Provider)
	session := ctn.Get(static.DiDiscordSession).(*discordgo.Session)

	w, err := waterlink.NewWaterlinkProvider(session,
		waterlink.WaterlinkConfig{Host: cfg.Config().Lavalink.Host, Password: cfg.Config().Lavalink.Password})

	if err != nil {
		logrus.Fatal("Failed to create waterlink provider: ", err)
		return nil
	}

	return w

}
