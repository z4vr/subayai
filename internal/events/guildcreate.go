package events

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di"
	"github.com/sirupsen/logrus"
	"github.com/z4vr/subayai/internal/services/config"
	"github.com/z4vr/subayai/internal/util/static"
	"github.com/z4vr/subayai/pkg/database"
	"github.com/z4vr/subayai/pkg/discordutils"
)

type GuildCreateEvent struct {
	db  database.Database
	cfg config.Provider
}

func NewGuildCreateEvent(ctn di.Container) *GuildCreateEvent {
	return &GuildCreateEvent{
		db:  ctn.Get(static.DiDatabase).(database.Database),
		cfg: ctn.Get(static.DiConfigProvider).(config.Provider),
	}
}

func (g *GuildCreateEvent) HandlerCreate(s *discordgo.Session, e *discordgo.GuildCreate) {

	// TODO: log earlier guild joins to prevent triggering this e

	limit := g.cfg.Config().Bot.GuildLimit
	if limit == -1 {
		return
	}

	if len(s.State.Guilds) >= limit {
		_, err := discordutils.SendMessageDM(s, e.OwnerID,
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

	return

}
