package comment

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"reddit-clone-example/internal/user"
	"reddit-clone-example/internal/storage"

	"github.com/google/uuid"
)

// Create creates a comment for a story. Author should be logged in.
func Create(ctx context.Context, author user.User, comment Comment) error {
	var (
		tx, err = storage.DB.BeginTxx(ctx, nil)
		l       = log.Fork().With("fn", "create", "author", author.ID, "story", comment.StoryID)
	)
	if err != nil {
		l.Log("err", err, "desc", "can't begin transaction")
		return errInternal
	}
	defer tx.Rollback()

	// Check the story really exists.
	{
		var id uuid.UUID
		const q = `SELECT id FROM story WHERE id = $1`
		if err = tx.QueryRowxContext(ctx, q, comment.StoryID).Scan(&id); err != nil && err != sql.ErrNoRows {
			l.Log("err", err, "sql", q, "desc", "db select failed")
			return errInternal
		}
		if err == sql.ErrNoRows {
			return errors.New("story not found")
		}
	}

	// Prepare comment for keeping.
	{
		comment.ID = uuid.New()
		comment.Author = author
		comment.CreatedAt = time.Now()
	}

	// Save into a storage.
	{
		const q = `INSERT INTO comment (id, story_id, body, created_by, created_at) VALUES ($1, $2, $3, $4, $5)`
		if _, err = tx.ExecContext(ctx, q, comment.ID, comment.StoryID, comment.Comment, comment.Author.ID, comment.CreatedAt); err != nil {
			l.Log("err", err, "sql", q, "desc", "new comment creation failed")
			return errInternal
		}
	}

	return tx.Commit()
}
