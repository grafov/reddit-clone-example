package vote

import (
	"context"
	"database/sql"

	"reddit-clone-example/internal/storage"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// List returns list of votes for the story and upvoting percent value.
func List(ctx context.Context, storyID uuid.UUID) ([]Vote, int, error) {
	var (
		tx, err = storage.DB.BeginTxx(ctx, nil)
		l       = log.Fork().With("fn", "list", "story", storyID)
	)
	if err != nil {
		l.Log("err", err, "desc", "can't begin transaction")
		return []Vote{}, 0, errInternal
	}
	defer tx.Rollback()

	// Load all article votes for the result.
	var (
		votes   []Vote
		upvotes int
	)
	if votes, upvotes, err = list(ctx, tx, storyID); err != nil {
		log.Log("err", err, "desc", "can't load votes")
		return []Vote{}, 0, errInternal
	}
	if err = tx.Commit(); err != nil {
		log.Log("err", err, "desc", "can't commit")
		return []Vote{}, 0, errInternal
	}

	return votes, upvotes, nil
}

// It returns list of votes and calculates upvote rating.
func list(ctx context.Context, tx *sqlx.Tx, storyID uuid.UUID) ([]Vote, int, error) {
	var (
		votes []Vote
		up    int
		rows  *sqlx.Rows
		err   error
	)
	const q = `SELECT account_id, story_id, vote FROM vote WHERE story_id = $1`
	if rows, err = tx.QueryxContext(ctx, q, storyID); err != nil && err != sql.ErrNoRows {
		return []Vote{}, 0, err
	}
	defer rows.Close()
	for rows.Next() {
		var v Vote
		if err = rows.StructScan(&v); err != nil {
			return []Vote{}, 0, err
		}
		if v.Value > 0 {
			up++
		}
		votes = append(votes, v)
	}

	var upvotes int
	if len(votes) > 0 {
		upvotes = up * 100 / len(votes)
	}

	return votes, upvotes, err
}
