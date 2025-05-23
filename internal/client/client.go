package client

import (
	"context"
	"fmt"

	"github.com/almottier/rivertui/config"
	"github.com/almottier/rivertui/internal/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/riverdriver/riverpgxv5"
)

// Client holds the database and River clients
type Client struct {
	Pool        *pgxpool.Pool
	RiverClient *river.Client[pgx.Tx]
}

// New creates a new client with database and River connections
func New(ctx context.Context, cfg *config.Config) (*Client, error) {
	pool, err := db.Connect(ctx, cfg.Database.URL, cfg.RefreshInterval)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	riverClient, err := river.NewClient[pgx.Tx](riverpgxv5.New(pool), &river.Config{})
	if err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to create River client: %w", err)
	}

	return &Client{
		Pool:        pool,
		RiverClient: riverClient,
	}, nil
}

// Close closes all client connections
func (c *Client) Close() {
	if c.Pool != nil {
		c.Pool.Close()
	}
}
