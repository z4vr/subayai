package waterlink

import (
	"github.com/bwmarrin/discordgo"
	"github.com/lukasl-dev/waterlink/v2"
)

type WaterlinkProvider struct {
	s      *discordgo.Session
	client *waterlink.Client
	conn   *waterlink.Connection

	address         string
	creds           waterlink.Credentials
	opts            waterlink.ConnectionOptions
	reconnectionTry int
}

type WaterlinkConfig struct {
	Host     string
	Password string
}
