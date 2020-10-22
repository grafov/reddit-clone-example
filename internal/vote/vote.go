package vote

import (
	"context"

	"reddit-clone-example/internal/storage"
	"reddit-clone-example/internal/user"

	"github.com/google/uuid"
)

// Upvote increments user' score for the story.
func Upvote(ctx context.Context, author user.User, storyID uuid.UUID) ([]Vote, error) {
	var (
		tx, err = storage.DB.BeginTxx(ctx, nil)
		l       = log.Fork().With("fn", "upvote", "voter", author.ID, "story", storyID)
	)
	if err != nil {
		l.Log("err", err, "desc", "can't begin transaction")
		return []Vote{}, errInternal
	}
	defer tx.Rollback()

	// Check that article exists.
	{
		// XXX
	}

	// Increment vote for the user.
	{
		const q = `UPDATE story SET views = views + 1 WHERE id = $1`
		if _, err = tx.ExecContext(ctx, q, storyID); err != nil {
			l.Log("err", err, "sql", q, "desc", "views update failed")
			return []Vote{}, errInternal
		}
	}

	if err = tx.Commit(); err != nil {
		log.Log("err", err, "desc", "can't commit")
		return []Vote{}, errInternal
	}

	return []Vote{}, nil

}

// Neutral resets user' score for the story.
func Neutral(ctx context.Context, author user.User, storyID uuid.UUID) ([]Vote, error) {
	return []Vote{}, nil
}

// Downvote decrements user' score for the story.
func Downvote(ctx context.Context, author user.User, storyID uuid.UUID) ([]Vote, error ){
	return []Vote{}, nil
}
