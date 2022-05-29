package waterlink

import (
	"github.com/lukasl-dev/waterlink/v2"
	"github.com/z4vr/subayai/pkg/discord"
)

type Waterlink struct {
	dc     *discord.Discord
	client *waterlink.Client
	conn   *waterlink.Connection

	address         string
	creds           waterlink.Credentials
	opts            waterlink.ConnectionOptions
	reconnectionTry int
}
