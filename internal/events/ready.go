package events

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"github.com/z4vr/subayai/pkg/discordutils"
)

type ReadyEvent struct {
}

func NewReadyEvent() *ReadyEvent {
	return &ReadyEvent{}
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
	logrus.Infof("Invite link: %s", discordutils.GetInviteLink(e.User.ID))
}
