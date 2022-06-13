package cmds

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/z4vr/subayai/internal/services/database"
	"github.com/z4vr/subayai/internal/services/discord"
	"github.com/z4vr/subayai/pkg/stringutils"
	"github.com/zekrotja/ken"
	"strings"
)

type Profile struct {
}

var (
	_ ken.Command = (*Profile)(nil)
)

func (c *Profile) Name() string {
	return "profile"
}

func (c *Profile) Description() string {
	return "Shows the profile of a user."
}

func (c *Profile) Version() string {
	return "1.0.0"
}

func (c *Profile) Type() discordgo.ApplicationCommandType {
	return discordgo.ChatApplicationCommand
}

func (c *Profile) Options() []*discordgo.ApplicationCommandOption {
	return []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionUser,
			Name:        "user",
			Description: "The user you want to fetch.",
			Required:    true,
		},
	}
}

func (c *Profile) Run(ctx *ken.Ctx) (err error) {

	if err = ctx.Defer(); err != nil {
		return
	}

	dc := ctx.Get("discord").(*discord.Discord)
	db := ctx.Get("database").(database.Database)

	user := ctx.Options().GetByName("user").UserValue(ctx)
	member, err := dc.Session().GuildMember(ctx.Event.GuildID, user.ID)
	if err != nil {
		return
	}
	guild, err := dc.Session().Guild(ctx.Event.GuildID)
	if err != nil {
		return
	}

	membRoleIDs := make(map[string]struct{})
	for _, rID := range member.Roles {
		membRoleIDs[rID] = struct{}{}
	}

	maxPos := len(guild.Roles)
	roleColor := 0xffffff
	for _, guildRole := range guild.Roles {
		if _, ok := membRoleIDs[guildRole.ID]; ok && guildRole.Position < maxPos && guildRole.Color != 0 {
			maxPos = guildRole.Position
			roleColor = guildRole.Color
		}
	}

	if err != nil {
		return err
	}
	createdTime, err := dc.GetDiscordSnowflakeCreationTime(member.User.ID)
	if err != nil {
		return err
	}

	roles := make([]string, len(member.Roles))
	for i, rID := range member.Roles {
		roles[i] = "<@&" + rID + ">"
	}
	embed := &discordgo.MessageEmbed{
		Color: roleColor,
		Title: fmt.Sprintf("About %s#%s", member.User.Username, member.User.Discriminator),
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: member.User.AvatarURL(""),
		},
		Fields: []*discordgo.MessageEmbedField{
			{
				Inline: true,
				Name:   "Tag",
				Value:  member.User.Username + "#" + member.User.Discriminator,
			},
			{
				Inline: true,
				Name:   "Nickname",
				Value:  stringutils.EnsureNotEmpty(member.Nick, "*no nick*"),
			},
			{
				Name:  "ID",
				Value: "```\n" + member.User.ID + "\n```",
			},
			{
				Name: "Guild Joined",
				Value: stringutils.EnsureNotEmpty(member.JoinedAt.Format("02.01.2006, 15:04"),
					"*failed parsing timestamp*"),
			},
			{
				Name: "Account Created",
				Value: stringutils.EnsureNotEmpty(createdTime.Format("02.01.2006, 15:04"),
					"*failed parsing timestamp*"),
			},
			{
				Name:  "Roles",
				Value: stringutils.EnsureNotEmpty(strings.Join(roles, ", "), "*no roles assigned*"),
			},
		},
	}

	// if the user is a bot, don't grab the xp info
	// and skip progressbar creation
	if member.User.Bot {
		embed.Description = ":robot:  **This is a bot account**"
		return ctx.FollowUpEmbed(embed).Error
	}

	// get the xp info
	currentXP, err := db.GetUserCurrentXP(member.User.ID, ctx.Event.GuildID)
	if err != nil {
		return ctx.FollowUpError("Failed to get current XP.", "An error occurred").Error
	}
	totalXP, err := db.GetUserTotalXP(member.User.ID, ctx.Event.GuildID)
	if err != nil {
		return ctx.FollowUpError("Failed to get total XP.", "An error occurred").Error
	}
	currentLevel, err := db.GetUserLevel(member.User.ID, ctx.Event.GuildID)
	if err != nil {
		return ctx.FollowUpError("Failed to get current level.", "An error occurred").Error
	}

	xpFields := []*discordgo.MessageEmbedField{
		{
			Name:  "Current Level",
			Value: fmt.Sprintf("%d", currentLevel),
		},
		{
			Name:  "Current XP",
			Value: fmt.Sprintf("%d", currentXP),
		},
		{
			Name:  "Total XP",
			Value: fmt.Sprintf("%dXP", totalXP),
		},
	}
	embed.Fields = append(embed.Fields, xpFields...)

	return ctx.FollowUpEmbed(embed).Error
}
