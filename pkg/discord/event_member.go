package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"github.com/z4vr/subayai/pkg/database/dberr"
	"github.com/z4vr/subayai/pkg/stringutils"
)

func (d *Discord) HandlerAutoRole(s *discordgo.Session, e *discordgo.GuildMemberAdd) {
	autoroleIDs, err := d.db.GetGuildAutoroleIDs(e.GuildID)
	if err != nil && err == dberr.ErrNotFound {
		logrus.WithField("guildID", e.GuildID).Warn("no autoroles found")
	}
	invalidAutoRoleIDs := make([]string, 0)
	for _, rid := range autoroleIDs {
		err = s.GuildMemberRoleAdd(e.GuildID, e.User.ID, rid)
		if apiErr, ok := err.(*discordgo.RESTError); ok && apiErr.Message.Code == discordgo.ErrCodeUnknownRole {
			invalidAutoRoleIDs = append(invalidAutoRoleIDs, rid)
		} else if err != nil {
			logrus.WithError(err).WithField("guildID", e.GuildID).WithField("userID",
				e.User.ID).Error("Failed setting autorole for member")
		}
	}
	if len(invalidAutoRoleIDs) > 0 {
		newAutoRoleIDs := make([]string, 0, len(autoroleIDs)-len(invalidAutoRoleIDs))
		for _, rid := range autoroleIDs {
			if !stringutils.ContainsAny(rid, invalidAutoRoleIDs) {
				newAutoRoleIDs = append(newAutoRoleIDs, rid)
			}
		}
		err = d.db.SetGuildAutoroleIDs(e.GuildID, newAutoRoleIDs)
		if err != nil {
			logrus.WithError(err).WithField("guildID", e.GuildID).WithField("userID",
				e.User.ID).Error("Failed updating auto role settings")
		}
	}
}
