package postgres

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/lib/pq"
	"github.com/z4vr/subayai/internal/models"
	"github.com/z4vr/subayai/internal/services/database"
)

type PGMiddleware struct {
	Db *sql.DB
}

var _ database.Database = (*PGMiddleware)(nil)

func (p *PGMiddleware) setup() (err error) {
	if err = p.Db.Ping(); err != nil {
		return
	}

	tx, err := p.Db.Begin()
	if err != nil {
		return
	}

	// create guild table
	_, err = tx.Exec(`
	CREATE TABLE IF NOT EXISTS "guild" (
		"guild_id" varchar (25) NOT NULL,
		"bot_message_channel_id" varchar (25) NOT NULL DEFAULT '',
		"level_up_message" text DEFAULT 'Well done {user}, your Level of wasting time just advanced to {level}!',
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

	// create level table
	_, err = tx.Exec(`
	CREATE TABLE IF NOT EXISTS "level" (
		"entry_id" serial NOT NULL,
		"user_id" varchar (25) NOT NULL,
		"guild_id" varchar (25) NOT NULL,
		"level" integer NOT NULL,
		"currentxp" integer NOT NULL,
		"totalxp" integer NOT NULL,
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
		"last_message_ts" timestamp NOT NULL,
		"last_session_ts" timestamp NOT NULL,
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
		"last_message_id" varchar (25) NOT NULL,
		"last_session_id" varchar (25) NOT NULL,
		PRIMARY KEY ("entry_id"));
	`)
	if err != nil {
		return
	}

	return
}

func (p *PGMiddleware) Connect(credentials ...interface{}) (err error) {
	creds := credentials[0].(models.DatabaseCreds)
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s",
		creds.User, creds.Password, creds.Host, creds.Database)

	if p.Db, err = sql.Open("p", dsn); err != nil {
		return
	}

	err = p.setup()
	return
}

func (p *PGMiddleware) Close() {
	if p.Db != nil {
		p.Db.Close()
	}
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
		err = database.ErrValueNotFound
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

// GetGuildLevelUpMessage returns the level up message for the guild
func (p *PGMiddleware) GetGuildLevelUpMessage(guildID string) (value string, err error) {
	return p.getGuildSetting(guildID, "level_up_message")
}

// SetGuildLevelUpMessage sets the level up message for the guild
func (p *PGMiddleware) SetGuildLevelUpMessage(guildID, value string) (err error) {
	return p.setGuildSetting(guildID, "level_up_message", value)
}

// GetGuildAfkChannelID returns the afk channel ID for the guild
func (p *PGMiddleware) GetGuildAfkChannelID(guildID string) (value string, err error) {
	return p.getGuildSetting(guildID, "afk_channel_id")
}

// SetGuildAfkChannelID sets the afk channel ID for the guild
func (p *PGMiddleware) SetGuildAfkChannelID(guildID, value string) (err error) {
	return p.setGuildSetting(guildID, "afk_channel_id", value)
}

// GetGuildAutoroleIDs returns the autorole IDs for the guild
func (p *PGMiddleware) GetGuildAutoroleIDs(guildID string) (value []string, err error) {

	roleString, err := p.getGuildSetting(guildID, "autorole_ids")
	if err != nil {
		return []string{}, err
	} else if roleString == "" {
		return []string{}, nil
	}

	return strings.Split(roleString, ";"), nil
}

// SetGuildAutoroleIDs sets the autorole IDs for the guild
func (p *PGMiddleware) SetGuildAutoroleIDs(guildID string, value []string) (err error) {
	roleString := strings.Join(value, ";")
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
