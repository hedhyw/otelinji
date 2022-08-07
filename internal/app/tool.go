package app

import (
	"context"
	"database/sql"
	"fmt"
)

type core struct {
	db *sql.DB
}

func (c core) Health(ctx context.Context) (err error) {
	if err = c.db.PingContext(ctx); err != nil {
		return fmt.Errorf("db ping: %w", err)
	}

	return nil
}
