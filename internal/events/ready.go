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

func (l *ReadyEvent) Handler(session *discordgo.Session, event *discordgo.Ready) {
	err := session.UpdateListeningStatus("slash commands [WIP]")
	if err != nil {
		return
	}
	logrus.WithFields(logrus.Fields{
		"id":       event.User.ID,
		"username": event.User.String(),
	}).Info("Signed in as:")
	logrus.Infof("Invite link: %s", discordutils.GetInviteLink(event.User.ID))
}
