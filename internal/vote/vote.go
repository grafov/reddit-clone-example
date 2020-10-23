package vote

import (
	"context"
	"database/sql"
	"errors"

	"reddit-clone-example/internal/storage"
	"reddit-clone-example/internal/user"

	"github.com/google/uuid"
)

// Upvote increments user' score for the story.
func Upvote(ctx context.Context, voter user.User, storyID uuid.UUID) ([]Vote, error) {
	return voteUpdate(ctx, voter, storyID, +1)
}

// Neutral resets user' score for the story.
func Neutral(ctx context.Context, voter user.User, storyID uuid.UUID) ([]Vote, error) {
	return voteUpdate(ctx, voter, storyID, 0)
}

// Downvote decrements user' score for the story.
func Downvote(ctx context.Context, voter user.User, storyID uuid.UUID) ([]Vote, error) {
	return voteUpdate(ctx, voter, storyID, -1)
}

func voteUpdate(ctx context.Context, voter user.User, storyID uuid.UUID, direction int64) ([]Vote, error) {
	var (
		tx, err = storage.DB.BeginTxx(ctx, nil)
		l       = log.Fork().With("fn", "vote-update", "voter", voter.ID, "story", storyID)
	)
	if err != nil {
		l.Log("err", err, "desc", "can't begin transaction")
		return []Vote{}, errInternal
	}
	defer tx.Rollback()

	// Check that article exists.
	{
		var score int64
		const q = `SELECT score FROM story WHERE id = $1`
		if err = tx.QueryRowxContext(ctx, q, storyID).Scan(&score); err != nil && err != sql.ErrNoRows {
			l.Log("err", err, "sql", q, "desc", "select story failed")
			return []Vote{}, err
		}
		if err == sql.ErrNoRows {
			return []Vote{}, errors.New("story not found")
		}
	}

	// Set vote for the user.
	var oldVote int64
	{
		const q = `SELECT vote FROM vote WHERE story_id = $1 AND account_id = $2`
		if err = tx.QueryRowxContext(ctx, q, storyID, voter.ID).Scan(&oldVote); err != nil && err != sql.ErrNoRows {
			l.Log("err", err, "sql", q, "desc", "vote load failed")
			return []Vote{}, errInternal
		}
		// Setup first vote for the user. Reset not need here.
		if err == sql.ErrNoRows && direction != 0 {
			const q = `INSERT INTO vote (story_id, account_id, vote) VALUES ($1, $2, $3)`
			if _, err = tx.ExecContext(ctx, q, storyID, voter.ID, direction); err != nil {
				l.Log("err", err, "sql", q, "desc", "vote insert failed")
				return []Vote{}, errInternal
			}
		}
		// This case handles score reset. Increment/decrement skiped here.
		if err == nil && direction == 0 {
			const q = `DELETE FROM vote WHERE story_id = $1 AND account_id = $2`
			if _, err = tx.ExecContext(ctx, q, storyID, voter.ID); err != nil {
				l.Log("err", err, "sql", q, "desc", "vote update failed")
				return []Vote{}, errInternal
			}
		}
		if err == nil && direction != 0 {
			const q = `UPDATE vote SET vote = $3 WHERE story_id = $1 AND account_id = $2`
			if _, err = tx.ExecContext(ctx, q, storyID, voter.ID, direction); err != nil {
				l.Log("err", err, "sql", q, "desc", "vote update failed")
				return []Vote{}, errInternal
			}
		}
	}

	// Update story' score counter. Update score only when vote
	// direction has changed.
	if oldVote != direction {
		const q = `UPDATE story SET score = score + $2 WHERE id = $1`
		if _, err = tx.ExecContext(ctx, q, storyID, -oldVote+direction); err != nil {
			l.Log("err", err, "sql", q, "desc", "vote load failed")
			return []Vote{}, errInternal
		}
	}

	// Load all article votes for the result.
	var votes []Vote
	if votes, _, err = list(ctx, tx, storyID); err != nil {
		log.Log("err", err, "desc", "can't load votes")
		return []Vote{}, errInternal
	}

	if err = tx.Commit(); err != nil {
		log.Log("err", err, "desc", "can't commit")
		return []Vote{}, errInternal
	}

	return votes, nil
}
