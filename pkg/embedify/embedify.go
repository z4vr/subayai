package embedify

import (
	"github.com/bwmarrin/discordgo"
	"github.com/z4vr/subayai/internal/static"
	"time"
)

// EmbedBuilder provides a builder pattern to
// create a discordgo message embed.
type EmbedBuilder struct {
	emb *discordgo.MessageEmbed
}

// New returns a fresh EmbedBuilder.
func New() *EmbedBuilder {
	return &EmbedBuilder{&discordgo.MessageEmbed{Color: static.ColorEmbedDefault}}
}

func (e *EmbedBuilder) SetTitle(title string) *EmbedBuilder {
	e.emb.Title = title
	return e
}

func (e *EmbedBuilder) SetDescription(description string) *EmbedBuilder {
	e.emb.Description = description
	return e
}

func (e *EmbedBuilder) SetURL(url string) *EmbedBuilder {
	e.emb.URL = url
	return e
}

func (e *EmbedBuilder) SetError() *EmbedBuilder {
	e.emb.Color = static.ColorEmbedRed
	return e
}

func (e *EmbedBuilder) SetSuccess() *EmbedBuilder {
	e.emb.Color = static.ColorEmbedGreen
	return e
}

func (e *EmbedBuilder) SetWarning() *EmbedBuilder {
	e.emb.Color = static.ColorEmbedYellow
	return e
}

func (e *EmbedBuilder) SetAuthor(author, iconURL string) *EmbedBuilder {
	e.emb.Author = &discordgo.MessageEmbedAuthor{
		Name:    author,
		IconURL: iconURL,
	}
	return e
}

func (e *EmbedBuilder) SetThumbnail(url string) *EmbedBuilder {
	e.emb.Thumbnail = &discordgo.MessageEmbedThumbnail{
		URL: url,
	}
	return e
}

func (e *EmbedBuilder) SetImage(url string) *EmbedBuilder {
	e.emb.Image = &discordgo.MessageEmbedImage{
		URL: url,
	}
	return e
}

func (e *EmbedBuilder) SetFooter(text, iconURL string) *EmbedBuilder {
	e.emb.Footer = &discordgo.MessageEmbedFooter{
		Text:    text,
		IconURL: iconURL,
	}
	return e
}

func (e *EmbedBuilder) SetTimestamp(timestamp time.Time) *EmbedBuilder {
	e.emb.Timestamp = timestamp.Format("02/01/2006 15:04:05")
	return e
}

func (e *EmbedBuilder) Build() *discordgo.MessageEmbed {
	return e.emb
}
