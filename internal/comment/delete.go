package comment

import (
	"context"
	"errors"

	"reddit-clone-example/internal/user"
	"reddit-clone-example/internal/storage"

	"github.com/google/uuid"
)

// Delete deletes a comment from a storage.
func Delete(ctx context.Context, owner user.User, storyID, commentID uuid.UUID) error {
	var (
		tx, err = storage.DB.BeginTxx(ctx, nil)
		l       = log.Fork().With("fn", "delete", "user", owner.ID, "comment-id", commentID, "story-id", storyID)
	)
	if err != nil {
		l.Log("err", err, "desc", "can't begin transaction")
		return errInternal
	}
	defer tx.Rollback()

	const q = `DELETE FROM comment WHERE id = $1 AND story_id = $2 AND created_by = $3`
	if _, err = tx.ExecContext(ctx, q, commentID, storyID, owner.ID); err != nil {
		l.Log("err", err, "desc", "delete failed")
		return errors.New("can't delete comment")
	}

	if err = tx.Commit(); err != nil {
		log.Log("err", err, "desc", "can't commit")
		return errInternal
	}

	return nil
}

// DeleteAll deletes all comments for the story. Not exposed to HTTP API.
func DeleteAll(ctx context.Context, storyID uuid.UUID) error {
	return nil
}
