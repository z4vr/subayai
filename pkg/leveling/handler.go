package leveling

import (
	"math/rand"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

type LevelingHandlers struct {
	lp *Provider
}

func NewLevelingHandlers(lp *Provider) *LevelingHandlers {
	return &LevelingHandlers{
		lp: lp,
	}
}

func (l *LevelingHandlers) MessageCreate(s *discordgo.Session, e *discordgo.MessageCreate) {

	var (
		botMessageChannelID string
		levelUpMessage      string
		levelData           *LevelData
		err                 error
	)

	if e.Author.ID == s.State.User.ID {
		return
	}

	if e.Author.Bot {
		return
	}

	if e.GuildID == "" {
		return
	}

	levelData, err = l.lp.Get(e.Author.ID, e.GuildID)
	if err != nil {
		return
	}

	// In case of a new user, we need to create a new entry in the database
	// and the guildMal.lp.
	if levelData == nil {
		levelData = &LevelData{
			UserID:    e.Author.ID,
			GuildID:   e.GuildID,
			Level:     0,
			CurrentXP: 0,
			TotalXP:   0,
		}
		l.lp.Set(e.Author.ID, e.GuildID, levelData)
		l.lp.SaveToDB(levelData)
	}

	// Now calulate the Xl.lp and level ul.lp if needed.
	earnedXP := rand.Intn(60) + 25
	levelup := levelData.AddXP(earnedXP, false)

	// If the user leveled ul.lp, we need to send a message to the channel
	if levelup {
		botMessageChannelID, err = l.lp.db.GetGuildBotMessageChannelID(e.GuildID)
		if err != nil {
			return
		} else if botMessageChannelID == "" {
			botMessageChannelID = e.ChannelID
		}

		levelUpMessage = strings.Replace(levelUpMessage, "{user}", e.Author.Mention(), -1)
		levelUpMessage = strings.Replace(levelUpMessage, "{leveling}", strconv.Itoa(levelData.Level), -1)

		_, err = s.ChannelMessageSend(botMessageChannelID, levelUpMessage)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"guildID":   e.GuildID,
				"userID":    e.Author.ID,
				"channelID": botMessageChannelID,
			}).WithError(err).Error("Failed to send level ul.lp message")
			return
		}
	}

	// Save the ul.lpdated level data to the database.
	err = l.lp.SaveToDB(levelData)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"guildID": levelData.GuildID,
			"userID":  levelData.UserID,
		}).WithError(err).Error("Failed to save level data to DB")
		return
	}

	return

}
