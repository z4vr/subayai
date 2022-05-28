package inits

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sarulabs/di"
	"github.com/z4vr/subayai/internal/services/database"
	"github.com/z4vr/subayai/internal/services/leveling"
	"github.com/z4vr/subayai/internal/util/static"
)

func NewLevelProvider(ctn di.Container) leveling.LevelProvider {
	s := ctn.Get(static.DiDiscordSession).(*discordgo.Session)
	db := ctn.Get(static.DiDatabase).(database.Database)

	provider := leveling.LevelProvider{
		Session: s,
		Db:      db,
	}

	// iterate over the current guilds
	// and load the data from the db into the map
	guilds := s.State.Guilds
	if len(guilds) == 0 {
		// assign an empty map to
		// the provider
		provider.Guilds = make(leveling.GuildData)

		return provider

	}

	guildData := make(leveling.GuildData, len(guilds))

	// iterate over the guilds
	for _, guild := range guilds {
		// create a new map for the guild
		guildData[guild.ID] = make(leveling.UserData, len(guild.Members))

		// iterate over the members
		for _, member := range guild.Members {
			levelData := provider.FetchFromDB(member.User.ID, guild.ID)
			if levelData == nil {
				continue
			}
			// add the data to the map
			guildData[guild.ID][member.User.ID] = levelData
		}
	}

	provider.Guilds = guildData

	return provider

}
