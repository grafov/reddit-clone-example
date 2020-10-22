package comment

import (
	"context"
	"time"

	"reddit-clone-example/internal/user"
	"reddit-clone-example/storage"

	"github.com/google/uuid"
)

func Create(ctx context.Context, author user.User, storyID uuid.UUID, comment Comment) (Comment, error) {
	var (
		tx, err = storage.DB.BeginTxx(ctx, nil)
		l       = log.Fork().With("fn", "create")
	)
	if err != nil {
		l.Log("err", err, "desc", "can't begin transaction")
		return Comment{}, errInternal
	}
	defer tx.Rollback()

	// Prepare comment for keeping.
	{
		comment.ID = uuid.New()
		comment.Author = author
		comment.CreatedAt = time.Now()
	}

	// Save into a storage.
	{
		const q = `INSERT INTO comment (id, story_id, body, created_by, created_at) VALUES ($1, $2, $3, $4, $5)`
		if _, err = tx.ExecContext(ctx, q, comment.ID, comment.StoryID, comment.Body, comment.Author.ID, comment.CreatedAt); err != nil {
			l.Log("err", err, "desc", "new comment creation failed")
			return Comment{}, errInternal
		}
	}

	return comment, tx.Commit()
}
