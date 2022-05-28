package events

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di"
	"github.com/sirupsen/logrus"
	"github.com/z4vr/subayai/internal/util/static"
	"github.com/z4vr/subayai/pkg/database"
	"github.com/z4vr/subayai/pkg/discordutils"
	leveling2 "github.com/z4vr/subayai/pkg/leveling"
	"math/rand"
	"time"
)

type VoiceStateUpdateEvent struct {
	db database.Database
	lp leveling2.LevelProvider
}

func NewVoiceStateUpdateEvent(ctn di.Container) *VoiceStateUpdateEvent {
	return &VoiceStateUpdateEvent{
		db: ctn.Get(static.DiDatabase).(database.Database),
		lp: ctn.Get(static.DiLevelProvider).(leveling2.LevelProvider),
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
		xpData               leveling2.Data
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
		// reward the leveling for the time spent in voice
		levelData := v.lp.Get(e.UserID, e.GuildID)
		if levelData == nil {
			levelData := &leveling2.Data{
				UserID:    e.UserID,
				GuildID:   e.GuildID,
				Level:     0,
				CurrentXP: 0,
				TotalXP:   0,
			}
			err := v.lp.Add(levelData)
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"uid": e.UserID,
					"gid": e.GuildID,
				}).WithError(err).Error("Failed to add user to level map")
			}
		}

		earnedXP := rand.Intn(60) + 25
		levelup := xpData.AddXP(earnedXP, false)

		err = v.lp.SaveLevelData(xpData)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"gid": e.GuildID,
				"uid": e.UserID,
			}).WithError(err).Error("Failed to update user leveling")
		}

		if levelup {
			// TODO: send a message to the user or to the bot channel
		}

	}

}
