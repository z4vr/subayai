package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

type ReadyEvent struct {
	dc *Discord
}

func NewReadyEvent(dc *Discord) *ReadyEvent {
	return &ReadyEvent{
		dc: dc,
	}
}

func (l *ReadyEvent) Handler(s *discordgo.Session, e *discordgo.Ready) {
	err := s.UpdateListeningStatus("slash commands [WIP]")
	if err != nil {
		return
	}
	logrus.WithFields(logrus.Fields{
		"id":       e.User.ID,
		"username": e.User.String(),
	}).Info("Signed in as:")
	logrus.Infof("Invite link: %s", l.dc.GetInviteLink())
}
