package story

import (
	"context"
	"database/sql"

	"reddit-clone-example/internal/user"
	"reddit-clone-example/storage"

	"github.com/jmoiron/sqlx"
)

// List returns full list of the stories. Without paging for
// simplicity. Just limit by 1000 for demo purposes.
func List(ctx context.Context) ([]Story, error) {
	var (
		tx, err = storage.DB.BeginTxx(ctx, nil)
		l       = log.Fork().With("fn", "list")
	)
	if err != nil {
		l.Log("err", err, "desc", "can't begin transaction")
		return []Story{}, errInternal
	}
	defer tx.Rollback()

	// Retrieve the story from the storage.
	var list []Story
	{
		const q = `
SELECT s.id, s.title, s.kind, s.body, s.category, s.score, s.created_by, s.created_at, a.login
FROM story s INNER JOIN account a ON s.created_by = a.id
ORDER BY created_at DESC
LIMIT 1000`
		var rows *sqlx.Rows
		if rows, err = tx.QueryxContext(ctx, q); err != nil && err != sql.ErrNoRows {
			l.Log("err", err, "desc", "db select failed")
			return []Story{}, errInternal
		}
		if err == sql.ErrNoRows {
			return []Story{}, nil
		}
		defer rows.Close()
		for rows.Next() {
			var story Story
			if err = rows.StructScan(&story); err != nil {
				l.Log("err", err, "desc", "db select failed")
				return []Story{}, errInternal
			}

			story.Author = user.User{ID: story.CreatedBy, Login: story.AuthorName}
			// TODO append comments here
			list = append(list, story)
		}
	}

	if err = tx.Commit(); err != nil {
		log.Log("err", err, "desc", "can't commit")
		return []Story{}, errInternal
	}

	return list, nil
}
