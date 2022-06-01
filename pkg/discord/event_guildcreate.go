package discord

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

func (d *Discord) HandlerCreate(s *discordgo.Session, e *discordgo.GuildCreate) {

	// check if the joinedAt is older than the time
	if e.JoinedAt.Unix() <= time.Now().Unix() {
		fmt.Println(1)
		return
	}

	limit := d.config.GuildLimit
	if limit == -1 {
		return
	}

	if len(s.State.Guilds) >= limit {
		_, err := d.SendMessageDM(e.OwnerID,
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
