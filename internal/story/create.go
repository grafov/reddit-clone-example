package story

import (
	"context"
	"errors"
	"time"

	"reddit-clone-example/internal/comment"
	"reddit-clone-example/internal/user"
	"reddit-clone-example/storage"

	"github.com/google/uuid"
)

func Create(ctx context.Context, author user.User, story Story) (Story, error) {
	var (
		tx, err = storage.DB.BeginTxx(ctx, nil)
		l       = log.Fork().With("fn", "create")
	)
	if err != nil {
		l.Log("err", err, "desc", "can't begin transaction")
		return Story{}, errInternal
	}
	defer tx.Rollback()

	// Prepare story for keeping.
	{
		story.ID = uuid.New()
		switch story.Type {
		case textType:
			story.Body = story.Text
		case linkType:
			story.Body = story.URL
		default:
			return Story{}, errors.New("unknown story type")
		}
		story.Author = author
		story.CreatedAt = time.Now()
		story.Comments = []comment.Comment{}
		story.Upvote(author.ID)
	}

	// Save into a storage.
	{
		const q = `INSERT INTO story (id, kind, title, body, category, created_by, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`
		if _, err = tx.ExecContext(ctx, q, story.ID, story.Type, story.Title, story.Body, story.Category, story.Author.ID, story.CreatedAt); err != nil {
			l.Log("err", err, "sql", q, "desc", "new story creation failed")
			return Story{}, errInternal
		}
	}

	return story, tx.Commit()
}
