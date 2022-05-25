package embedify

import "github.com/bwmarrin/discordgo"

func (c *Colors) GetColorFromKind(kind string) int {
	switch kind {
	case "success":
		return c.SuccessColor
	case "general":
		return c.GeneralColor
	case "error":
		return c.ErrorColor
	case "warning":
		return c.WarningColor
	default:
		return 0x5865F2
	}
}

func (b *Builder) Embed(embedKind string, options *EmbedOptions) *discordgo.MessageEmbed {

	embed := &discordgo.MessageEmbed{
		Title:       options.Title,
		Description: options.Description,
		Color:       b.colors.GetColorFromKind(embedKind),
		Fields:      options.Fields,
		Author:      options.Author,
		Footer:      options.Footer,
		Image:       options.Image,
		Thumbnail:   options.Thumbnail,
	}
	return embed
}
