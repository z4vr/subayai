package inits

import (
	"fmt"
	"github.com/z4vr/subayai/internal/services/config"

	"github.com/z4vr/subayai/internal/util/static"

	"github.com/lukasl-dev/waterlink/v2"
	"github.com/sarulabs/di"
	"github.com/sirupsen/logrus"
)

type WaterlinkProvider struct {
	wlClient *waterlink.Client
}

func NewWaterlinkProvider(ctn di.Container) (p *WaterlinkProvider) {

	var err error

	p = &WaterlinkProvider{}
	cfg := ctn.Get(static.DiConfigProvider).(config.Provider)

	creds := waterlink.Credentials{
		Authorization: cfg.Config().Lavalink.Authorization,
	}

	p.wlClient, err = waterlink.NewClient(fmt.Sprintf("http://%s",
		cfg.Config().Lavalink.Host), creds)
	if err != nil {
		logrus.WithError(err).Fatal("failed to create waterlink client")
	}

	// TODO: finish up the lavalink provider

	return

}
