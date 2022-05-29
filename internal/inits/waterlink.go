package inits

import (
	"github.com/sirupsen/logrus"
	"github.com/z4vr/subayai/pkg/waterlink"
)

func NewWaterlinkProvider() *waterlink.Waterlink {

	w, err := waterlink.New(session,
		waterlink.Config{Host: cfg.Config().Lavalink.Host, Password: cfg.Config().Lavalink.Password})

	if err != nil {
		logrus.Fatal("Failed to create waterlink provider: ", err)
		return nil
	}

	return w

}
