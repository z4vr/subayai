package events

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di"
	"github.com/sirupsen/logrus"
	"github.com/z4vr/subayai/internal/services/database"
	xp2 "github.com/z4vr/subayai/internal/services/level"
	"github.com/z4vr/subayai/internal/util/static"
	"github.com/z4vr/subayai/pkg/discordutils"
	"math/rand"
	"time"
)

type VoiceStateUpdateEvent struct {
	db database.Database
}

func NewVoiceStateUpdateEvent(ctn di.Container) *VoiceStateUpdateEvent {
	return &VoiceStateUpdateEvent{
		db: ctn.Get(static.DiDatabase).(database.Database),
	}
}

func (v *VoiceStateUpdateEvent) HandlerXP(s *discordgo.Session, e *discordgo.VoiceStateUpdate) {

	// check if the member is a bot
	member, err := discordutils.GetMember(s, e.GuildID, e.UserID)
	if err != nil {
		logrus.
			WithFields(logrus.Fields{
				"gid": e.GuildID,
				"uid": e.UserID,
			}).
			Error("Failed to get member")
		return
	}

	if member.User.Bot {
		return
	}

	var (
		afkChannelID         string = "empty"
		lastSessionID        string = ""
		lastSessionTimestamp int64  = 0
		nowTimestamp         int64  = time.Now().Unix()
		xpData               *xp2.Data
	)

	afkChannelID, err = v.db.GetGuildAFKChannelID(e.GuildID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"gid": e.GuildID,
		}).WithError(err).Error("Failed to get afk channel id")
	}
	lastSessionID, err = v.db.GetLastVoiceSessionID(e.GuildID, e.UserID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"gid": e.GuildID,
			"uid": e.UserID,
		}).WithError(err).Error("Failed to get last session id")
	}
	lastSessionTimestamp, err = v.db.GetLastVoiceSessionTimestamp(e.GuildID, e.UserID)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"gid": e.GuildID,
			"uid": e.UserID,
		}).WithError(err).Error("Failed to get last session timestamp")
	}

	// scenario: user freshly joined voice channel, or rejoined from afk channel
	if e.BeforeUpdate == nil || e.BeforeUpdate.ChannelID == afkChannelID {
		// update our last records
		lastSessionID = e.SessionID
		lastSessionTimestamp = nowTimestamp

		// save them to db
		err = v.db.SetLastVoiceSessionID(e.GuildID, e.UserID, lastSessionID)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"gid": e.GuildID,
				"uid": e.UserID,
			}).WithError(err).Error("Failed to set last session id")
		}
		err = v.db.SetLastVoiceSessionTimestamp(e.GuildID, e.UserID, lastSessionTimestamp)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"gid": e.GuildID,
				"uid": e.UserID,
			}).WithError(err).Error("Failed to set last session timestamp")
		}

	}

	// scenario: user left voice channel
	if e.VoiceState.ChannelID == "" && e.VoiceState.SessionID == lastSessionID ||
		e.VoiceState.ChannelID == afkChannelID && e.VoiceState.SessionID == lastSessionID {
		// reward the level for the time spent in voice
		xpData, err = xp2.GetUserXP(e.UserID, e.GuildID, v.db)
		if err != nil && err == database.ErrValueNotFound {
			logrus.WithFields(logrus.Fields{
				"gid": e.GuildID,
				"uid": e.UserID,
			}).WithError(err).Error("Failed to get user level")
			xpData = &xp2.Data{
				UserID:    e.UserID,
				GuildID:   e.GuildID,
				CurrentXP: 0,
				TotalXP:   0,
				Level:     0,
			}
		}

		xpEarned := (rand.Intn(50) + 25) * (int(nowTimestamp - lastSessionTimestamp)) / 600
		_ = xpData.AddXP(xpEarned, false)

		err = xp2.UpdateUserXP(xpData, v.db)

	}

}
