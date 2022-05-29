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

	ColorEmbedRed     = 0xd32f2f
	ColorEmbedDefault = 0x249ff2
	ColorEmbedGreen   = 0x8bc34a
	ColorEmbedYellow  = 0xffeb3b

	AppIcon = "https://raw.githubusercontent.com/z4vr/subayai/main/assets/app-icon.jpg"

	// EMOJIS

)
