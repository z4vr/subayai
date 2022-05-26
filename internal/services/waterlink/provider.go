package waterlink

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/gompus/snowflake"
	"github.com/lukasl-dev/waterlink/v2"
	"github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

func NewWaterlinkProvider(s *discordgo.Session, c WaterlinkConfig) (*WaterlinkProvider, error) {
	var w WaterlinkProvider
	var err error

	w.s = s
	w.address = c.Host
	w.creds = waterlink.Credentials{
		Authorization: c.Password,
		UserID:        snowflake.MustParse(w.s.State.User.ID),
		ResumeKey:     "subayaiSession",
	}

	w.client, err = waterlink.NewClient(fmt.Sprintf("http://%s", c.Host), w.creds)
	if err != nil {
		return nil, err
	}

	if err = w.Connect(); err != nil {
		return nil, err
	}

	w.s.AddHandler(w.handleVoiceServerUpdate)

	return &w, nil
}

func (w *WaterlinkProvider) handleVoiceServerUpdate(s *discordgo.Session, e *discordgo.VoiceServerUpdate) {
	logrus.
		WithField("guild", e.GuildID).
		WithField("sessionID", s.State.SessionID).
		Debugf("Update voice server: %+v", e)

	g := w.conn.Guild(snowflake.MustParse(e.GuildID))
	err := g.UpdateVoice(s.State.SessionID, e.Token, e.Endpoint)
	if err != nil {
		logrus.
			WithError(err).
			WithField("guild", e.GuildID).
			WithField("sessionID", s.State.SessionID).
			Error("Voice server update failed")
	}
}

func (w *WaterlinkProvider) Connect() error {
	if w.conn != nil && !w.conn.Closed() {
		return errors.New("connection already established")
	}
	var err error
	w.conn, err = waterlink.Open(fmt.Sprintf("ws://%s", w.address), w.creds, w.opts)
	return err
}

func (w *WaterlinkProvider) tryReconnecting() {
	timeout := time.Duration(w.reconnectionTry)*500*time.Millisecond + time.Duration(rand.Intn(900)+100)*time.Millisecond

	if timeout > 30*time.Second {
		timeout = 30 * time.Second
	}

	logrus.
		WithField("try", w.reconnectionTry).
		WithField("timeout", timeout).
		Warn("Lavalink: Trying to reconnect ...")
	time.Sleep(timeout)

	w.reconnectionTry++
	err := w.Connect()
	if err != nil {
		logrus.WithError(err)
		w.tryReconnecting()
		return
	}

	w.reconnectionTry = 0
	logrus.Info("Lavalink connection re-established")
}
