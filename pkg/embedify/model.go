package embedify

import "github.com/bwmarrin/discordgo"

type Builder struct {
	session *discordgo.Session
	colors  Colors
}

type Colors struct {
	ErrorColor   int
	SuccessColor int
	WarningColor int
	GeneralColor int
}

type EmbedOptions struct {
	Title       string
	Description string
	Thumbnail   *discordgo.MessageEmbedThumbnail
	Image       *discordgo.MessageEmbedImage
	Footer      *discordgo.MessageEmbedFooter
	Author      *discordgo.MessageEmbedAuthor
	Fields      []*discordgo.MessageEmbedField
}
