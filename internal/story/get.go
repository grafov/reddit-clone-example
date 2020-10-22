package story

import (
	"context"
	"database/sql"
	"errors"

	"reddit-clone-example/internal/user"
	"reddit-clone-example/storage"

	"github.com/google/uuid"
)

func Get(ctx context.Context, id uuid.UUID) (Story, error) {
	var (
		tx, err = storage.DB.BeginTxx(ctx, nil)
		l       = log.Fork().With("fn", "get", "id", id)
	)
	if err != nil {
		l.Log("err", err, "desc", "can't begin transaction")
		return Story{}, errInternal
	}
	defer tx.Rollback()

	// Retrieve the story from the storage.
	var story Story
	{
		const q = `SELECT id, title, kind, body, category, score, created_by, created_at FROM story WHERE id = $1`
		if err = tx.QueryRowxContext(ctx, q, id).StructScan(&story); err != nil && err != sql.ErrNoRows {
			l.Log("err", err, "desc", "db select failed")
			return Story{}, errInternal
		}
		if err == sql.ErrNoRows {
			return Story{}, errors.New("story not found")
		}
	}

	// Fields postprocessing.
	{
		story.Views++
		const q = `UPDATE story SET views = views + 1 WHERE id = $1`
		if _, err = tx.ExecContext(ctx, q, id); err != nil {
			l.Log("err", err, "desc", "views update failed")
			return Story{}, errInternal
		}
		var u user.User
		if u, err = user.Get(ctx, tx, story.CreatedBy); err != nil {
			l.Log("err", err, "desc", "load of author info failed")
			return Story{}, errInternal
		}
		story.Author = u
	}

	if err = tx.Commit(); err != nil {
		log.Log("err", err, "desc", "can't commit")
		return Story{}, errInternal
	}

	return story, nil
}
