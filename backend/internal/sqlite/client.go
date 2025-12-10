package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"przepisyapi/internal/sqlc"

	_ "embed"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed przepisy.sqlite
var noGoose []byte

type Client struct {
	Path string

	DB      *sql.DB
	Queries *sqlc.Queries
}

func New(ctx context.Context, dbPath string) (*Client, error) {
	c := &Client{
		Path: dbPath,
	}
	if !c.dbExists() {
		err := c.writeDb()
		if err != nil {
			return nil, fmt.Errorf("failed to write default db: %w", err)
		}
	}
	db, err := sql.Open("sqlite3", c.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite db: %w", err)
	}
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping sqlite db: %w", err)
	}
	c.DB = db
	c.Queries = sqlc.New(db)
	return c, nil
}

func (c *Client) dbExists() bool {
	_, err := os.Stat(c.Path)
	return os.IsExist(err) || err == nil
}

func (c *Client) writeDb() error {
	err := os.WriteFile(c.Path, noGoose, 0o644)
	if err != nil {
		return fmt.Errorf("failed to create default db: %w", err)
	}
	return nil
}
