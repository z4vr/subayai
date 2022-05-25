package events

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di"
	"github.com/sirupsen/logrus"
	"github.com/z4vr/subayai/internal/services/database"
	"github.com/z4vr/subayai/internal/util/static"
	"github.com/z4vr/subayai/pkg/stringutils"
)

type GuildMemberAddEvent struct {
	db database.Database
}

func NewGuildMemberAddEvent(ctn di.Container) *GuildMemberAddEvent {
	return &GuildMemberAddEvent{
		db: ctn.Get(static.DiDatabase).(database.Database),
	}
}

func (g *GuildMemberAddEvent) HandlerAutoRole(session *discordgo.Session, event *discordgo.GuildMemberAdd) {
	autoroleIDs, err := g.db.GetGuildAutoroleIDs(event.GuildID)
	if err != nil && err == database.ErrValueNotFound {
		logrus.WithField("gid", event.GuildID).Warn("no autoroles found")
	}
	invalidAutoRoleIDs := make([]string, 0)
	for _, rid := range autoroleIDs {
		err = session.GuildMemberRoleAdd(event.GuildID, event.User.ID, rid)
		if apiErr, ok := err.(*discordgo.RESTError); ok && apiErr.Message.Code == discordgo.ErrCodeUnknownRole {
			invalidAutoRoleIDs = append(invalidAutoRoleIDs, rid)
		} else if err != nil {
			logrus.WithError(err).WithField("gid", event.GuildID).WithField("uid",
				event.User.ID).Error("Failed setting autorole for member")
		}
	}
	if len(invalidAutoRoleIDs) > 0 {
		newAutoRoleIDs := make([]string, 0, len(autoroleIDs)-len(invalidAutoRoleIDs))
		for _, rid := range autoroleIDs {
			if !stringutils.ContainsAny(rid, invalidAutoRoleIDs) {
				newAutoRoleIDs = append(newAutoRoleIDs, rid)
			}
		}
		err = g.db.SetGuildAutoroleIDs(event.GuildID, newAutoRoleIDs)
		if err != nil {
			logrus.WithError(err).WithField("gid", event.GuildID).WithField("uid",
				event.User.ID).Error("Failed updating auto role settings")
		}
	}
}