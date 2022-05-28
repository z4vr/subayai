package inits

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di"
	"github.com/sirupsen/logrus"
	"github.com/z4vr/subayai/internal/services/config"
	"github.com/z4vr/subayai/internal/util/static"
	waterlink2 "github.com/z4vr/subayai/pkg/waterlink"
)

func NewWaterlinkProvider(ctn di.Container) *waterlink2.Waterlink {

	cfg := ctn.Get(static.DiConfigProvider).(config.Provider)
	session := ctn.Get(static.DiDiscordSession).(*discordgo.Session)

	w, err := waterlink2.New(session,
		waterlink2.WaterlinkConfig{Host: cfg.Config().Lavalink.Host, Password: cfg.Config().Lavalink.Password})

	if err != nil {
		logrus.Fatal("Failed to create waterlink provider: ", err)
		return nil
	}

	return w

}
