package discordutils

import (
	"fmt"
	"github.com/z4vr/subayai/internal/util/static"
	"strconv"
	"time"
)

// GetInviteLink returns the invite link for the bot's account with the specified permissions.
func GetInviteLink(userID string) string {
	return fmt.Sprintf("https://discord.com/api/oauth2/authorize?client_id=%s&scope=%s&permissions=%d",
		userID, static.OAuthScopes, static.InvitePermission)
}

// GetDiscordSnowflakeCreationTime returns the time.Time object of the passed snowflake.
func GetDiscordSnowflakeCreationTime(snowflake string) (time.Time, error) {
	sfI, err := strconv.ParseInt(snowflake, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	timestamp := (sfI >> 22) + 1420070400000
	return time.Unix(timestamp/1000, timestamp), nil
}
