package sqlitestorage

import (
	_ "embed"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"

	"github.com/frioux/shortlinks/shortlinks"
)

//go:embed schema.sql
var schema string

func Connect(dsn string) (*Client, error) {
	if dsn == "" {
		dsn = "file:db.db"
	}
	db, err := sqlx.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(schema); err != nil {
		return nil, fmt.Errorf("couldn't create schema: %s", err)
	}

	return &Client{db: db}, nil
}

type Client struct {
	db *sqlx.DB
}

func (c Client) Shortlink(from string) (shortlinks.Shortlink, error) {
	ret := shortlinks.Shortlink{}
	err := c.db.Get(&ret, `SELECT "from", "to" FROM shortlinks WHERE "from" = ? AND "deleted" IS NULL`, from)

	if err != nil {
		return ret, fmt.Errorf("couldn't load shortlink (%s): %w", from, err)
	}
	return ret, nil
}

func (c Client) CreateShortlink(s shortlinks.Shortlink) error {
	_, err := c.db.Exec(`INSERT INTO shortlinks("from", "to") VALUES (?, ?)
			  ON CONFLICT("from") DO
			  UPDATE SET "to"=excluded."to", "deleted" = null`, s.From, s.To)

	if err != nil {
		return fmt.Errorf("couldn't insert shortlink (%s): %w", s.From, err)
	}
	return nil
}

func (c Client) DeleteShortlink(from string) error {
	_, err := c.db.Exec(`UPDATE shortlinks SET "deleted"=CURRENT_TIMESTAMP WHERE "from" = ?`, from)

	if err != nil {
		return fmt.Errorf("couldn't delete shortlink (%s): %w", from, err)
	}
	return nil
}

func (c Client) AllShortlinks() ([]shortlinks.Shortlink, error) {
	ret := []shortlinks.Shortlink{}
	err := c.db.Select(&ret, `SELECT "to", "from" FROM shortlinks WHERE "deleted" IS NULL ORDER BY "from"`)
	if err != nil {
		return nil, fmt.Errorf("couldn't load shortlinks: %w", err)
	}
	return ret, nil
}

func (c Client) DeletedShortlinks() ([]shortlinks.Shortlink, error) {
	ret := []shortlinks.Shortlink{}
	err := c.db.Select(&ret, `SELECT "to", "from" FROM shortlinks WHERE "deleted" IS NOT NULL ORDER BY "from"`)
	if err != nil {
		return nil, fmt.Errorf("couldn't load shortlinks: %w", err)
	}
	return ret, nil
}

func (c Client) History(from string) ([]shortlinks.History, error) {
	ret := []shortlinks.History{}
	err := c.db.Select(&ret, `SELECT "to", "from", "when" FROM history WHERE "from" = ?`, from)
	if err != nil {
		return nil, fmt.Errorf("couldn't load history (for %s): %w", from, err)
	}
	return ret, nil
}

func (c Client) InsertHistory(h shortlinks.History) error {
	_, err := c.db.Exec(`INSERT INTO history("from", "to", "when") VALUES (?, ?, CURRENT_TIMESTAMP)`, h.From, h.To)
	if err != nil {
		return fmt.Errorf("couldn't insert history (%s): %w", h.From, err)
	}
	return nil
}
