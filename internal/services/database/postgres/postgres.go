package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/z4vr/subayai/internal/models"
	"github.com/z4vr/subayai/internal/services/database"
)

type PGMiddleware struct {
	Db *sql.DB
}

var (
	_                database.Database = (*PGMiddleware)(nil)
	ErrValueNotFound                   = errors.New("value not found")
)

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
		"afk_channel_id" text NOT NULL DEFAULT '',
		"autorole_ids" text NOT NULL DEFAULT '',
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
