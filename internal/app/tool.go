package app

import (
	"context"
	"database/sql"
	"fmt"
)

type User struct{}

type repo struct {
	db *sql.DB
}

func (c repo) Health(ctx context.Context) (err error) {
	if err = c.db.PingContext(ctx); err != nil {
		return fmt.Errorf("db ping: %w", err)
	}

	return nil
}
