package events

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di"
	"github.com/sirupsen/logrus"
	"github.com/z4vr/subayai/internal/services/config"
	"github.com/z4vr/subayai/internal/services/database"
	"github.com/z4vr/subayai/internal/util/static"
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

func (g *GuildCreateEvent) HandlerCreate(session *discordgo.Session, event *discordgo.GuildCreate) {

	// TODO: log earlier guild joins to prevent triggering this event

	limit := g.cfg.Config().Bot.GuildLimit
	if limit == -1 {
		return
	}

	if len(session.State.Guilds) >= limit {
		_, err := discordutils.SendMessageDM(session, event.OwnerID,
			fmt.Sprintf("Sorry, the instance owner disallowed me to join more than %d guilds.", limit))
		if err != nil {
			logrus.WithError(err).WithFields(logrus.Fields{
				"guild": event.ID,
			}).Error("Failed to send message")
			return
		}
		err = session.GuildLeave(event.ID)
		if err != nil {
			logrus.WithError(err).WithFields(logrus.Fields{
				"guild": event.ID,
			}).Error("Failed to leave guild")
		}

	}

	return

}
