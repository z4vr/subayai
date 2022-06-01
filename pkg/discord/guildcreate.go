package discord

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"github.com/z4vr/subayai/pkg/database"
)

type GuildCreateEvent struct {
	dc  *Discord
	db  *database.Database
	cfg *Config
}

func NewGuildCreateEvent(dc *Discord, db *database.Database, cfg *Config) *GuildCreateEvent {
	return &GuildCreateEvent{
		dc:  dc,
		db:  db,
		cfg: cfg,
	}
}

func (g *GuildCreateEvent) HandlerCreate(s *discordgo.Session, e *discordgo.GuildCreate) {

	// check if the joinedAt is older than the time
	if e.JoinedAt.Unix() <= time.Now().Unix() {
		fmt.Println(1)
		return
	}

	limit := g.cfg.GuildLimit
	if limit == -1 {
		return
	}

	if len(s.State.Guilds) >= limit {
		_, err := g.dc.SendMessageDM(e.OwnerID,
			fmt.Sprintf("Sorry, the instance owner disallowed me to join more than %d guilds.", limit))
		if err != nil {
			logrus.WithError(err).WithFields(logrus.Fields{
				"guild": e.ID,
			}).Error("Failed to send message")
			return
		}
		err = s.GuildLeave(e.ID)
		if err != nil {
			logrus.WithError(err).WithFields(logrus.Fields{
				"guild": e.ID,
			}).Error("Failed to leave guild")
		}

	}

}
