package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/z4vr/subayai/pkg/database"
	"strconv"
	"time"
)

var (
	OAuthScopes = "bot%20applications.commands"

	InvitePermission = discordgo.PermissionViewChannel |
		discordgo.PermissionSendMessages |
		discordgo.PermissionEmbedLinks |
		discordgo.PermissionReadMessageHistory |
		discordgo.PermissionUseExternalEmojis |
		//discordgo.PermissionAddReactions |
		discordgo.PermissionManageRoles |
		discordgo.PermissionManageChannels

	Intents = discordgo.IntentsAll
)

type Discord struct {
	session *discordgo.Session
	config  Config
	db      database.Database
}

func New(c Config, db database.Database) (*Discord, error) {
	var t Discord
	var err error

	t.config = c
	t.db = db
	t.session, err = discordgo.New("Bot " + c.Token)

	t.session.Identify.Intents = discordgo.MakeIntent(Intents)

	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (d *Discord) Open() error {

	d.session.AddHandler(NewReadyEvent(d).Handler)
	d.session.AddHandler(NewGuildCreateEvent(d, &d.db, &d.config).HandlerCreate)
	d.session.AddHandler(NewGuildMemberAddEvent(d.db).HandlerAutoRole)

	return d.session.Open()
}

func (d *Discord) Close() {
	d.session.Close()
}

func (d *Discord) Session() *discordgo.Session {
	return d.session
}

func (d *Discord) AddHandler(f func()) {
	d.session.AddHandler(f)
}

func (d *Discord) GetInviteLink() string {
	return fmt.Sprintf("https://discord.com/api/oauth2/authorize?client_id=%s&scope=%s&permissions=%d",
		d.session.State.User.ID, OAuthScopes, InvitePermission)
}

func (d *Discord) GetDiscordSnowflakeCreationTime(snowflake string) (time.Time, error) {
	sfI, err := strconv.ParseInt(snowflake, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	timestamp := (sfI >> 22) + 1420070400000
	return time.Unix(timestamp/1000, timestamp), nil
}

func (d *Discord) SendMessageDM(userID, message string) (msg *discordgo.Message, err error) {
	ch, err := d.session.UserChannelCreate(userID)
	if err != nil {
		return
	}
	msg, err = d.session.ChannelMessageSend(ch.ID, message)
	return
}

func (d *Discord) SendEmbedMessageDM(userID string, embed *discordgo.MessageEmbed) (msg *discordgo.Message, err error) {
	ch, err := d.session.UserChannelCreate(userID)
	if err != nil {
		return
	}
	msg, err = d.session.ChannelMessageSendEmbed(ch.ID, embed)
	return
}

func (d *Discord) GetMember(userID, guildID string) (discordgo.Member, error) {
	member, err := d.session.State.Member(guildID, userID)
	if err == nil {
		return *member, nil
	}

	member, err = d.session.GuildMember(guildID, userID)
	return *member, err
}

func (d *Discord) GetGuild(id string) (discordgo.Guild, error) {
	guild, err := d.session.State.Guild(id)
	if err == nil {
		return *guild, nil
	}

	guild, err = d.session.Guild(id)
	return *guild, err
}

func (d *Discord) UsersInGuildVoice(guildID string) ([]string, error) {
	g, err := d.session.State.Guild(guildID)
	if err != nil {
		return nil, err
	}

	userIDs := make([]string, 0, len(g.VoiceStates))
	for _, vs := range g.VoiceStates {
		if vs.UserID != d.session.State.User.ID {
			userIDs = append(userIDs, vs.UserID)
		}
	}

	return userIDs, nil
}

func (d *Discord) FindUserVS(userID string) (discordgo.VoiceState, bool) {
	for _, g := range d.session.State.Guilds {
		for _, vs := range g.VoiceStates {
			for vs.UserID == userID {
				return *vs, true
			}
		}
	}
	return discordgo.VoiceState{}, false
}
