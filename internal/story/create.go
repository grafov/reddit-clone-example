package story

import (
	"context"

	"reddit-clone-example/internal/user"
	"reddit-clone-example/storage"
)

func Create(ctx context.Context, author user.User, story Story) error {
	var (
		tx, err = storage.DB.BeginTxx(ctx, nil)
		l       = log.Fork().With("fn", "login")
	)
	if err != nil {
		l.Log("err", err, "desc", "can't begin transaction")
		return err
	}
	defer tx.Rollback()

	return tx.Commit()
}
