package embedify

import (
	"github.com/bwmarrin/discordgo"
	"github.com/z4vr/subayai/pkg/discordutils"
)

func NewBuilderProvider(session *discordgo.Session, colors Colors) *Builder {
	return &Builder{
		session: session,
		colors:  colors,
	}
}

func (b *Builder) SendEmbed(channelID string, embed *discordgo.MessageEmbed) (*discordgo.Message, error) {
	return b.session.ChannelMessageSendEmbed(channelID, embed)
}

func (b *Builder) SendEmbedDM(channelID string, embed *discordgo.MessageEmbed) (*discordgo.Message, error) {
	return discordutils.SendEmbedMessageDM(b.session, channelID, embed)
}
