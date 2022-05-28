package waterlink

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/gompus/snowflake"
	"github.com/lukasl-dev/waterlink/v2"
	"github.com/lukasl-dev/waterlink/v2/track"
	"github.com/lukasl-dev/waterlink/v2/track/query"
	"github.com/sirupsen/logrus"
	"math/rand"
	"time"
)

func New(s *discordgo.Session, c WaterlinkConfig) (*Waterlink, error) {
	var w Waterlink
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

func (w *Waterlink) handleVoiceServerUpdate(s *discordgo.Session, e *discordgo.VoiceServerUpdate) {
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

func (w *Waterlink) Connect() error {
	if w.conn != nil && !w.conn.Closed() {
		return errors.New("connection already established")
	}
	var err error
	w.conn, err = waterlink.Open(fmt.Sprintf("ws://%s", w.address), w.creds, w.opts)
	return err
}

func (w *Waterlink) Play(guildID, ident string) (track.Track, error) {
	tracks, err := w.client.LoadTracks(query.Of(ident))
	if err != nil {
		return track.Track{}, err
	}

	logrus.WithField("type", tracks.LoadType).
		WithField("n", len(tracks.Tracks)).
		Debug("Tracks loaded")

	if len(tracks.Tracks) == 0 {
		return track.Track{}, errors.New("no tracks have been loaded")
	}

	sf, err := snowflake.Parse(guildID)
	if err != nil {
		return track.Track{}, err
	}

	track := tracks.Tracks[0]
	return track, w.conn.Guild(sf).PlayTrack(track)
}

func (w *Waterlink) Pause(guildID string) error {
	sf, err := snowflake.Parse(guildID)
	if err != nil {
		return err
	}

	return w.conn.Guild(sf).SetPaused(true)
}

func (w *Waterlink) Resume(guildID string) error {
	sf, err := snowflake.Parse(guildID)
	if err != nil {
		return err
	}

	return w.conn.Guild(sf).SetPaused(false)
}

func (w *Waterlink) Stop(guildID string) error {
	sf, err := snowflake.Parse(guildID)
	if err != nil {
		return err
	}

	return w.conn.Guild(sf).Stop()
}

func (w *Waterlink) Disconnect() error {
	if w.conn == nil {
		return errors.New("connection not established")
	}
	return w.conn.Close()
}

func (w *Waterlink) tryReconnecting() {
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
