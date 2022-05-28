package waterlink

import (
	"github.com/bwmarrin/discordgo"
	"github.com/lukasl-dev/waterlink/v2"
)

type Waterlink struct {
	s      *discordgo.Session
	client *waterlink.Client
	conn   *waterlink.Connection

	address         string
	creds           waterlink.Credentials
	opts            waterlink.ConnectionOptions
	reconnectionTry int
}
