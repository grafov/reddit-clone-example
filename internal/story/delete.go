package story

import (
	"context"
	"errors"

	"reddit-clone-example/internal/user"
	"reddit-clone-example/internal/comment"
	"reddit-clone-example/internal/vote"
	"reddit-clone-example/internal/storage"

	"github.com/google/uuid"
)

// Delete deletes a story from a storage.
func Delete(ctx context.Context, owner user.User, id uuid.UUID) error {
	var (
		tx, err = storage.DB.BeginTxx(ctx, nil)
		l       = log.Fork().With("fn", "delete", "user", owner.ID, "id", id)
	)
	if err != nil {
		l.Log("err", err, "desc", "can't begin transaction")
		return errInternal
	}
	defer tx.Rollback()

	const q = `DELETE FROM story WHERE id = $1 AND created_by = $2`
	if _, err = tx.ExecContext(ctx, q, id, owner.ID); err != nil {
		l.Log("err", err, "sql", q, "desc", "delete failed")
		return errors.New("can't delete story")
	}

	if err=comment.DeleteAll(ctx, id);err!= nil{
		log.Log("err", err, "desc", "can't delete story' comments")
	}

	if err=	vote.DeleteAll(ctx, id);err!=nil{
		log.Log("err", err, "desc", "can't delete story' votes")
	}

	if err = tx.Commit(); err != nil {
		log.Log("err", err, "desc", "can't commit")
		return errInternal
	}

	return nil
}
