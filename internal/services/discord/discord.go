package discord

import (
	"errors"
	"fmt"
	"github.com/sarulabs/di/v2"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/z4vr/subayai/internal/services/config"
	"github.com/z4vr/subayai/internal/services/database"
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

	ErrUserHasRole    = errors.New("user has role")
	ErrUserHasNotRole = errors.New("user has not role")
)

type Discord struct {
	session *discordgo.Session
	cfg     config.Config
	db      database.Database
}

func New(ctn di.Container) (*Discord, error) {
	var t Discord
	var err error

	db := ctn.Get("database").(database.Database)
	cfg := ctn.Get("config").(config.Config)

	t.cfg = cfg
	t.db = db
	t.session, err = discordgo.New("Bot " + cfg.Discord.Token)

	if err != nil {
		return nil, err
	}

	t.session.Identify.Intents = discordgo.MakeIntent(Intents)

	t.session.State.TrackMembers = true
	t.session.State.TrackVoice = true

	// adding event handlers
	t.session.AddHandler(t.Ready)
	t.session.AddHandler(t.MessageLeveling)
	t.session.AddHandler(t.VoiceLeveling)
	t.session.AddHandler(t.AutoRole)
	t.session.AddHandler(t.GuildLimit)

	return &t, nil
}

func (d *Discord) Open() error {
	return d.Session().Open()
}

func (d *Discord) Close() error {
	return d.Session().Close()
}

func (d *Discord) Session() *discordgo.Session {
	return d.session
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

func (d *Discord) GetMember(userID, guildID string) (*discordgo.Member, error) {
	member, err := d.session.State.Member(guildID, userID)
	if err == nil {
		return member, nil
	}

	member, err = d.session.GuildMember(guildID, userID)
	return member, err
}

func (d *Discord) GetGuild(id string) (*discordgo.Guild, error) {
	guild, err := d.session.State.Guild(id)
	if err == nil {
		return guild, nil
	}

	guild, err = d.session.Guild(id)
	return guild, err
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

func (d *Discord) FindGuildTextChannel(guildID string) *discordgo.Channel {
	guild, err := d.GetGuild(guildID)
	if err != nil {
		return &discordgo.Channel{}
	}
	for _, c := range guild.Channels {
		if c.Type == discordgo.ChannelTypeGuildText {
			// check if bot has permission to write in channel
			return c
		}
	}
	return &discordgo.Channel{}
}

func (d *Discord) AddRoleToUser(userID, guildID, roleID string) error {
	member, err := d.GetMember(userID, guildID)
	if err != nil {
		return err
	}

	for _, r := range member.Roles {
		if r == roleID {
			return ErrUserHasRole
		}
	}

	err = d.session.GuildMemberRoleAdd(guildID, userID, roleID)
	if err != nil {
		return err
	}

	return nil
}

func (d *Discord) RemoveRoleFromUser(userID, guildID, roleID string) error {
	member, err := d.GetMember(userID, guildID)
	if err != nil {
		return err
	}

	for _, r := range member.Roles {
		if r == roleID {
			err = d.session.GuildMemberRoleRemove(guildID, userID, roleID)
			if err != nil {
				return err
			}
			return nil
		}
		return ErrUserHasNotRole
	}
	return ErrUserHasNotRole
}
