package postgres

import (
	"database/sql"
	"fmt"
	"github.com/z4vr/subayai/pkg/database/dberr"
	"strings"

	_ "github.com/lib/pq"
)

type PGMiddleware struct {
	Db *sql.DB
}

func (p *PGMiddleware) setup() (err error) {
	// ping the database to make sure it's up
	if err := p.Db.Ping(); err != nil {
		return err
	}

	tx, err := p.Db.Begin()
	if err != nil {
		return err
	}

	// create guild table
	_, err = tx.Exec(`
	CREATE TABLE IF NOT EXISTS "guild" (
		"guild_id" varchar (25) NOT NULL,
		"bot_message_channel_id" varchar (25) NOT NULL DEFAULT '',
		"level_up_message" text DEFAULT 'Well done {user}, your Level of wasting time just advanced to {leveling}!',
		"afk_channel_id" varchar (25) NOT NULL DEFAULT '',
		"autorole_ids" varchar (25) NOT NULL DEFAULT '',
		PRIMARY KEY ("guild_id"));
	`)
	if err != nil {
		return
	}

	// create user table
	_, err = tx.Exec(`
	CREATE TABLE IF NOT EXISTS "user" (
		"user_id" varchar (25) NOT NULL,
		PRIMARY KEY ("user_id"));
	`)
	if err != nil {
		return
	}

	// create leveling table
	_, err = tx.Exec(`
	CREATE TABLE IF NOT EXISTS "leveling" (
		"entry_id" serial NOT NULL,
		"user_id" varchar (25) NOT NULL,
		"guild_id" varchar (25) NOT NULL,
		"leveling" integer NOT NULL DEFAULT 0,
		"current_xp" integer NOT NULL DEFAULT 0,
		"total_xp" integer NOT NULL DEFAULT 0,
		PRIMARY KEY ("entry_id"));
	`)
	if err != nil {
		return
	}

	// create timestamp table
	_, err = tx.Exec(`
	CREATE TABLE IF NOT EXISTS "timestamp" (
		"entry_id" serial NOT NULL,
		"user_id" varchar (25) NOT NULL,
		"guild_id" varchar (25) NOT NULL,
		"last_guild_message" timestamp NOT NULL,
		"last_voice_session" timestamp NOT NULL,
		PRIMARY KEY ("entry_id"));
	`)
	if err != nil {
		return
	}

	// create discordids table
	_, err = tx.Exec(`
	CREATE TABLE IF NOT EXISTS "discordids" (
		"entry_id" serial NOT NULL,
		"user_id" varchar (25) NOT NULL,
		"guild_id" varchar (25) NOT NULL,
		"last_guild_message" varchar (25) NOT NULL,
		"last_voice_session" varchar (25) NOT NULL,
		PRIMARY KEY ("entry_id"));
	`)
	if err != nil {
		return
	}

	return tx.Commit()
}

func (p *PGMiddleware) Connect(credentials ...interface{}) (err error) {
	// connect to the database
	creds := credentials[0].(Config)
	pgi := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		creds.Host, creds.Port, creds.Username, creds.Password, creds.Database)
	db, err := sql.Open("postgres", pgi)
	if err != nil {
		return err
	}

	// set the database
	p.Db = db

	// set up the database
	return p.setup()
}

func (p *PGMiddleware) getGuildSetting(guildID, setting string) (value string, err error) {

	err = p.Db.QueryRow(`
	SELECT `+setting+` FROM guild WHERE guild_id = $1;
	`, guildID).Scan(&value)

	err = wrapNotFound(err)

	return
}

func (p *PGMiddleware) setGuildSetting(guildID, setting, value string) (err error) {

	res, err := p.Db.Exec(`
	UPDATE guild SET `+setting+` = $1 WHERE guild_id = $2;
	`, value, guildID)
	if err != nil {
		return
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return
	}

	if rows == 0 {
		_, err = p.Db.Exec(`
		INSERT INTO guild (guild_id, `+setting+`) VALUES ($1, $2);
		`, guildID, value)
		if err != nil {
			return
		}
	}
	return
}

func wrapNotFound(err error) error {
	if err == sql.ErrNoRows {
		err = dberr.ErrNotFound
	}
	return err
}

// GetGuildBotMessageChannelID returns the bot message channel ID for the guild
func (p *PGMiddleware) GetGuildBotMessageChannelID(guildID string) (value string, err error) {
	return p.getGuildSetting(guildID, "bot_message_channel_id")
}

// SetGuildBotMessageChannelID sets the bot message channel ID for the guild
func (p *PGMiddleware) SetGuildBotMessageChannelID(guildID, value string) (err error) {
	return p.setGuildSetting(guildID, "bot_message_channel_id", value)
}

// GetGuildLevelUpMessage returns the leveling up message for the guild
func (p *PGMiddleware) GetGuildLevelUpMessage(guildID string) (message string, err error) {
	return p.getGuildSetting(guildID, "level_up_message")
}

// SetGuildLevelUpMessage sets the leveling up message for the guild
func (p *PGMiddleware) SetGuildLevelUpMessage(guildID, message string) (err error) {
	return p.setGuildSetting(guildID, "level_up_message", message)
}

// GetGuildAFKChannelID returns the afk channel ID for the guild
func (p *PGMiddleware) GetGuildAFKChannelID(guildID string) (channelID string, err error) {
	return p.getGuildSetting(guildID, "afk_channel_id")
}

// SetGuildAFKChannelID sets the afk channel ID for the guild
func (p *PGMiddleware) SetGuildAFKChannelID(guildID, channelID string) (err error) {
	return p.setGuildSetting(guildID, "afk_channel_id", channelID)
}

// GetGuildAutoroleIDs returns the autorole IDs for the guild
func (p *PGMiddleware) GetGuildAutoroleIDs(guildID string) (roleIDs []string, err error) {

	roleString, err := p.getGuildSetting(guildID, "autorole_ids")
	if err != nil {
		return []string{}, err
	} else if roleString == "" {
		return []string{}, nil
	}

	return strings.Split(roleString, ";"), nil
}

// SetGuildAutoroleIDs sets the autorole IDs for the guild
func (p *PGMiddleware) SetGuildAutoroleIDs(guildID string, roleIDs []string) (err error) {
	roleString := strings.Join(roleIDs, ";")
	return p.setGuildSetting(guildID, "autorole_ids", roleString)
}

func (p *PGMiddleware) getUserSetting(userID, setting string) (value string, err error) {

	err = p.Db.QueryRow(`
	SELECT `+setting+` FROM user WHERE user_id = $1;
	`, userID).Scan(&value)

	err = wrapNotFound(err)

	return
}

func (p *PGMiddleware) setUserSetting(userID, setting, value string) (err error) {

	res, err := p.Db.Exec(`
	UPDATE user SET `+setting+` = $1 WHERE user_id = $2;
	`, value, userID)
	if err != nil {
		return
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return
	}

	if rows == 0 {
		_, err = p.Db.Exec(`
		INSERT INTO user (user_id, `+setting+`) VALUES ($1, $2);
		`, userID, value)
		if err != nil {
			return
		}
	}
	return
}

func (p *PGMiddleware) getLevelSetting(userID, guildID, setting string) (value int, err error) {

	err = p.Db.QueryRow(`
	SELECT `+setting+` FROM leveling WHERE user_id = $1 AND guild_id = $2;
	`, userID, guildID).Scan(&value)

	err = wrapNotFound(err)

	return
}

func (p *PGMiddleware) setLevelSetting(userID, guildID, setting string, value int) (err error) {

	res, err := p.Db.Exec(`
	UPDATE leveling SET `+setting+` = $1 WHERE user_id = $2 AND guild_id = $3;
	`, value, userID, guildID)
	if err != nil {
		return
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return
	}

	if rows == 0 {
		_, err = p.Db.Exec(`
		INSERT INTO leveling (user_id, guild_id, `+setting+`) VALUES ($1, $2, $3);
		`, userID, guildID, value)
		if err != nil {
			return
		}
	}
	return
}

// GetUserLevel returns the leveling for the user
func (p *PGMiddleware) GetUserLevel(userID, guildID string) (level int, err error) {
	return p.getLevelSetting(userID, guildID, "leveling")
}

// SetUserLevel sets the leveling for the user
func (p *PGMiddleware) SetUserLevel(userID, guildID string, level int) (err error) {
	return p.setLevelSetting(userID, guildID, "leveling", level)
}

// GetUserCurrentXP returns the current leveling for the user
func (p *PGMiddleware) GetUserCurrentXP(userID, guildID string) (xp int, err error) {
	return p.getLevelSetting(userID, guildID, "current_xp")
}

// SetUserCurrentXP sets the current leveling for the user
func (p *PGMiddleware) SetUserCurrentXP(userID, guildID string, xp int) (err error) {
	return p.setLevelSetting(userID, guildID, "current_xp", xp)
}

// GetUserTotalXP returns the total leveling for the user
func (p *PGMiddleware) GetUserTotalXP(userID, guildID string) (xp int, err error) {
	return p.getLevelSetting(userID, guildID, "total_xp")
}

// SetUserTotalXP sets the total leveling for the user
func (p *PGMiddleware) SetUserTotalXP(userID, guildID string, xp int) (err error) {
	return p.setLevelSetting(userID, guildID, "total_xp", xp)
}

func (p *PGMiddleware) getTimestampSetting(userID, guildID, setting string) (value int64, err error) {

	err = p.Db.QueryRow(`
	SELECT `+setting+` FROM timestamp WHERE user_id = $1 AND guild_id = $2;
	`, userID, guildID).Scan(&value)

	err = wrapNotFound(err)

	return
}

func (p *PGMiddleware) setTimestampSetting(userID, guildID, setting string, value int64) (err error) {

	res, err := p.Db.Exec(`
	UPDATE timestamp SET `+setting+` = $1 WHERE user_id = $2 AND guild_id = $3;
	`, value, userID, guildID)
	if err != nil {
		return
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return
	}

	if rows == 0 {
		_, err = p.Db.Exec(`
		INSERT INTO timestamp (user_id, guild_id, `+setting+`) VALUES ($1, $2, $3);
		`, userID, guildID, value)
		if err != nil {
			return
		}
	}
	return
}

// GetLastMessageTimestamp returns the timestamp of the last message sent by the user
func (p *PGMiddleware) GetLastMessageTimestamp(userID, guildID string) (timestamp int64, err error) {
	return p.getTimestampSetting(userID, guildID, "last_guild_message")
}

// SetLastMessageTimestamp sets the timestamp of the last message sent by the user
func (p *PGMiddleware) SetLastMessageTimestamp(userID, guildID string, timestamp int64) (err error) {
	return p.setTimestampSetting(userID, guildID, "last_guild_message", timestamp)
}

// GetLastVoiceSessionTimestamp returns the timestamp of the last voice session by the user
func (p *PGMiddleware) GetLastVoiceSessionTimestamp(userID, guildID string) (timestamp int64, err error) {
	return p.getTimestampSetting(userID, guildID, "last_voice_session")
}

// SetLastVoiceSessionTimestamp sets the timestamp of the last voice session by the user
func (p *PGMiddleware) SetLastVoiceSessionTimestamp(userID, guildID string, timestamp int64) (err error) {
	return p.setTimestampSetting(userID, guildID, "last_voice_session", timestamp)
}

func (p *PGMiddleware) getDiscordIDsSetting(userID, guildID, setting string) (value string, err error) {

	err = p.Db.QueryRow(`
	SELECT `+setting+` FROM discord_ids WHERE user_id = $1 AND guild_id = $2;
	`, userID, guildID).Scan(&value)

	err = wrapNotFound(err)

	return
}

func (p *PGMiddleware) setDiscordIDsSetting(userID, guildID, setting string, value string) (err error) {

	res, err := p.Db.Exec(`
	UPDATE discord_ids SET `+setting+` = $1 WHERE user_id = $2 AND guild_id = $3;
	`, value, userID, guildID)
	if err != nil {
		return
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return
	}

	if rows == 0 {
		_, err = p.Db.Exec(`
		INSERT INTO discord_ids (user_id, guild_id, `+setting+`) VALUES ($1, $2, $3);
		`, userID, guildID, value)
		if err != nil {
			return
		}
	}
	return
}

// GetLastMessageID returns the last message id sent by the user
func (p *PGMiddleware) GetLastMessageID(userID, guildID string) (id string, err error) {
	return p.getDiscordIDsSetting(userID, guildID, "last_guild_message")
}

// SetLastMessageID sets the last message id sent by the user
func (p *PGMiddleware) SetLastMessageID(userID, guildID string, id string) (err error) {
	return p.setDiscordIDsSetting(userID, guildID, "last_guild_message", id)
}

// GetLastVoiceSessionID returns the last voice session id sent by the user
func (p *PGMiddleware) GetLastVoiceSessionID(userID, guildID string) (id string, err error) {
	return p.getDiscordIDsSetting(userID, guildID, "last_voice_session")
}

// SetLastVoiceSessionID sets the last voice session id sent by the user
func (p *PGMiddleware) SetLastVoiceSessionID(userID, guildID string, id string) (err error) {
	return p.setDiscordIDsSetting(userID, guildID, "last_voice_session", id)
}

func (p *PGMiddleware) Close() {
	if p.Db != nil {
		err := p.Db.Close()
		if err != nil {
			return
		}
	}
}
