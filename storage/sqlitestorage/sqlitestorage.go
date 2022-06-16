package sqlitestorage

import (
	"embed"
	"fmt"
	"io/fs"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/frioux/dh"

	"github.com/frioux/shortlinks/shortlinks"
)

//go:embed dh
var dhFS embed.FS

func Connect(dsn string) (*Client, error) {
	if dsn == "" {
		dsn = "file:db.db"
	}
	db, err := sqlx.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	var found struct { C int }
	const sql = `SELECT COUNT(*) AS c FROM main.sqlite_master WHERE "name" = 'dh_migrations' AND "type" = 'table'`;
	if err := db.Get(&found, sql); err != nil {
		return nil, fmt.Errorf("db.Get: %w", err)
	}

	e := dh.NewMigrator()
	if found.C != 1 {
		if err := e.MigrateOne(db, dh.DHMigrations, "000-sqlite"); err != nil {
			return nil, fmt.Errorf("dh.Migrator.MigrateOne: %w", err)
		}
	}

	fss, _ := fs.Sub(dhFS, "dh")
	if err := e.MigrateAll(db, fss); err != nil {
		return nil, fmt.Errorf("dh.Migrator.MigrateAll: %w", err)
	}

	return &Client{db: db}, nil
}

type Client struct {
	db *sqlx.DB
}

func (c Client) Shortlink(from string) (shortlinks.Shortlink, error) {
	ret := shortlinks.Shortlink{}
	err := c.db.Get(&ret, `SELECT "from", "to", "description" FROM shortlinks WHERE "from" = ? AND "deleted" IS NULL`, from)

	if err != nil {
		return ret, fmt.Errorf("couldn't load shortlink (%s): %w", from, err)
	}
	return ret, nil
}

func (c Client) CreateShortlink(s shortlinks.Shortlink) error {
	_, err := c.db.Exec(`INSERT INTO shortlinks("from", "to", "description") VALUES (?, ?, ?)
			  ON CONFLICT("from") DO
			  UPDATE SET
			  "to"          = "excluded"."to",
			  "deleted"     = null,
			  "description" = "excluded"."description"`, s.From, s.To, s.Description)

	if err != nil {
		return fmt.Errorf("couldn't insert shortlink (%s): %w", s.From, err)
	}
	return nil
}

func (c Client) DeleteShortlink(from, who string) error {
	if err := c.InsertHistory(shortlinks.History{From: from, To: "«deleted»", Who: who}); err != nil {
		return fmt.Errorf("couldn't insert delete history for shortlink (%s): %w", from, err)
	}

	_, err := c.db.Exec(`UPDATE shortlinks SET "deleted"=CURRENT_TIMESTAMP WHERE "from" = ?`, from)

	if err != nil {
		return fmt.Errorf("couldn't delete shortlink (%s): %w", from, err)
	}
	return nil
}

func (c Client) AllShortlinks() ([]shortlinks.Shortlink, error) {
	ret := []shortlinks.Shortlink{}
	err := c.db.Select(&ret, `SELECT "to", "from", "description" FROM shortlinks WHERE "deleted" IS NULL ORDER BY "from"`)
	if err != nil {
		return nil, fmt.Errorf("couldn't load shortlinks: %w", err)
	}
	return ret, nil
}

func (c Client) DeletedShortlinks() ([]shortlinks.Shortlink, error) {
	ret := []shortlinks.Shortlink{}
	err := c.db.Select(&ret, `SELECT "to", "from", "description" FROM shortlinks WHERE "deleted" IS NOT NULL ORDER BY "from"`)
	if err != nil {
		return nil, fmt.Errorf("couldn't load shortlinks: %w", err)
	}
	return ret, nil
}

func (c Client) History(from string) ([]shortlinks.History, error) {
	ret := []shortlinks.History{}
	err := c.db.Select(&ret, `SELECT "to", "from", "when", "who", "description" FROM history WHERE "from" = ?`, from)
	if err != nil {
		return nil, fmt.Errorf("couldn't load history (for %s): %w", from, err)
	}
	return ret, nil
}

func (c Client) InsertHistory(h shortlinks.History) error {
	_, err := c.db.Exec(`INSERT INTO history("from", "to", "when", "who", "description") VALUES (?, ?, CURRENT_TIMESTAMP, ?, ?)`, h.From, h.To, h.Who, h.Description)
	if err != nil {
		return fmt.Errorf("couldn't insert history (%s): %w", h.From, err)
	}
	return nil
}
