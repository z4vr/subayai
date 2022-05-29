package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/z4vr/subayai/pkg/discord/events"
)

type Discord struct {
	session *discordgo.Session
}

func New(c Config) (*Discord, error) {
	var t Discord
	var err error

	t.session, err = discordgo.New("Bot " + c.Token)
	t.session.State.TrackVoice = true
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (t *Discord) Open() error {

	t.session.AddHandler(events.NewReadyEvent().Handler)

	return t.session.Open()
}

func (t *Discord) Close() {
	t.session.Close()
}

func (t *Discord) Session() *discordgo.Session {
	return t.session
}

func (t *Discord) AddHandler(f func()) {
	t.session.AddHandler(f)
}
