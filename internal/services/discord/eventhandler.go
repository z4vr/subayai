package discord

import (
	"github.com/z4vr/subayai/internal/services/config"
	"github.com/z4vr/subayai/internal/services/database"
	"github.com/z4vr/subayai/internal/services/leveling"
)

func NewEventHandler(d *Discord, lp *leveling.Provider, db database.Database, cfg config.Config) *EventHandler {
	return &EventHandler{
		d:   d,
		lp:  lp,
		db:  db,
		cfg: cfg,
	}
}

type EventHandler struct {
	d   *Discord
	lp  *leveling.Provider
	db  database.Database
	cfg config.Config
}
