package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

func (d *Discord) Ready(s *discordgo.Session, e *discordgo.Ready) {
	err := s.UpdateListeningStatus("slash commands [WIP]")
	if err != nil {
		return
	}
	logrus.WithFields(logrus.Fields{
		"id":       e.User.ID,
		"username": e.User.String(),
	}).Info("Signed in as:")
	logrus.Infof("Invite link: %s", d.GetInviteLink())
}
