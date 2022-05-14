package static

import (
	"github.com/bwmarrin/discordgo"
)

const (

	// BOT INTENTS AND PERMS

	Intents = discordgo.IntentsAll

	InvitePermission = discordgo.PermissionViewChannel |
		discordgo.PermissionSendMessages |
		discordgo.PermissionEmbedLinks |
		discordgo.PermissionReadMessageHistory |
		discordgo.PermissionUseExternalEmojis |
		//discordgo.PermissionAddReactions |
		discordgo.PermissionManageRoles |
		discordgo.PermissionManageChannels

	OAuthScopes = "bot%20applications.commands"

	// BOT STATICS

	ColorEmbedError   = 0xd32f2f
	ColorEmbedDefault = 0x249ff2
	ColorEmbedGray    = 0xb0bec5
	ColorEmbedGreen   = 0x8bc34a

	// EMOJIS

)
