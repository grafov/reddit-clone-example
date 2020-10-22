package vote

import (
	"context"
	"database/sql"

	"reddit-clone-example/internal/storage"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// List returns current stats for the story.
func List(ctx context.Context, storyID uuid.UUID) ([]Vote, error) {
	var (
		tx, err = storage.DB.BeginTxx(ctx, nil)
		l       = log.Fork().With("fn", "list", "story", storyID)
	)
	if err != nil {
		l.Log("err", err, "desc", "can't begin transaction")
		return []Vote{}, errInternal
	}
	defer tx.Rollback()

	// Load all article votes for the result.
	var votes []Vote
	if votes, err = list(ctx, tx, storyID); err != nil {
		log.Log("err", err, "desc", "can't load votes")
		return []Vote{}, errInternal
	}
	if err = tx.Commit(); err != nil {
		log.Log("err", err, "desc", "can't commit")
		return []Vote{}, errInternal
	}

	return votes, nil
}

func list(ctx context.Context, tx *sqlx.Tx, storyID uuid.UUID) ([]Vote, error) {
	var (
		votes []Vote
		rows  *sqlx.Rows
		err   error
	)
	const q = `SELECT account_id, story_id, vote FROM vote WHERE story_id = $1`
	if rows, err = tx.QueryxContext(ctx, q, storyID); err != nil && err != sql.ErrNoRows {
		return []Vote{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var v Vote
		if err = rows.StructScan(&v); err != nil {
			return []Vote{}, err
		}
		votes = append(votes, v)
	}

	return votes, err
}
