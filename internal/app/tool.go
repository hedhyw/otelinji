package app

import (
	"context"
	"fmt"
)

func (c repo) GetUser(ctx context.Context, id string) (u model.User, err error) {
	err = c.db.
		QueryContext(ctx, `SELECT * FROM users WHERE id = ?`, id).
		Scan(&u)
	if err != nil {
		return model.User{}, fmt.Errorf("quering user: %w", err)
	}

	return u, nil
}
