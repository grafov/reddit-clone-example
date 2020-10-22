package comment

import (
	"context"
	"database/sql"

	"reddit-clone-example/internal/user"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

// List returns full list of the comments for the selected
// story. Without paging for simplicity. Just limit by 1000 for demo
// purposes. The function not presented in HTTP API, it used
// internally.
func List(ctx context.Context, tx *sqlx.Tx, storyID uuid.UUID) ([]Comment, error) {
	var (
		list []Comment
		l    = log.Fork().With("fn", "list")
		err  error
	)
	// Retrieve the com from the storage.
	{
		const q = `
SELECT c.id, c.body, c.created_by, c.created_at, a.login
FROM comment c INNER JOIN account a ON c.created_by = a.id
WHERE c.story_id = $1
ORDER BY c.created_at
LIMIT 1000`
		var rows *sqlx.Rows
		l.With("sql", q)
		if rows, err = tx.QueryxContext(ctx, q, storyID); err != nil && err != sql.ErrNoRows {
			l.Log("err", err, "desc", "db select failed")
			return []Comment{}, errInternal
		}
		if err == sql.ErrNoRows {
			return []Comment{}, nil
		}
		defer rows.Close()
		for rows.Next() {
			var com Comment
			if err = rows.StructScan(&com); err != nil {
				l.Log("err", err, "desc", "db select failed")
				return []Comment{}, errInternal
			}

			com.Author = user.User{ID: com.CreatedBy, Login: com.AuthorName}
			list = append(list, com)
		}
	}

	return list, nil
}
