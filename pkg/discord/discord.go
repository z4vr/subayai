package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/z4vr/subayai/pkg/database"
)

type Discord struct {
	session *discordgo.Session
	config  Config
	db      database.Database
}

func New(c Config, db database.Database) (*Discord, error) {
	var t Discord
	var err error

	t.config = c
	t.db = db
	t.session, err = discordgo.New("Bot " + c.Token)
	t.session.State.TrackVoice = true
	t.session.State.TrackMembers = true

	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (t *Discord) Open() error {

	t.session.AddHandler(NewReadyEvent().Handler)
	t.session.AddHandler(NewGuildCreateEvent(t.db, t.config).HandlerCreate)

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
