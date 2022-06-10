package events

import (
	"github.com/z4vr/subayai/internal/services/database"
	"github.com/z4vr/subayai/internal/services/discord"
	"github.com/z4vr/subayai/internal/services/leveling"
)

func New(d *discord.Discord, lp *leveling.Provider, db database.Database, cfg discord.Config) *EventHandler {
	return &EventHandler{
		d:   d,
		lp:  lp,
		db:  db,
		cfg: cfg,
	}
}

type EventHandler struct {
	d   *discord.Discord
	lp  *leveling.Provider
	db  database.Database
	cfg discord.Config
}
