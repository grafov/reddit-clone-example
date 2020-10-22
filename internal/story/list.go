package story

import (
	"context"
	"database/sql"
	"errors"

	"reddit-clone-example/internal/comment"
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
ORDER BY s.created_at DESC
LIMIT 1000`
		var rows *sqlx.Rows
		if rows, err = tx.QueryxContext(ctx, q); err != nil && err != sql.ErrNoRows {
			l.Log("err", err, "sql", q, "desc", "db select failed")
			return []Story{}, errInternal
		}
		if err == sql.ErrNoRows {
			return []Story{}, nil
		}
		defer rows.Close()
		for rows.Next() {
			var story Story
			if err = rows.StructScan(&story); err != nil {
				l.Log("err", err, "sql", q, "desc", "rows scan failed")
				return []Story{}, errInternal
			}

			story.Author = user.User{ID: story.CreatedBy, Login: story.AuthorName}
			list = append(list, story)
		}
	}

	// Add comments to the story list. They added in a separated
	// loop for avoid fetching nested query in the same
	// transaction.
	for i, s := range list {
		if list[i].Comments, err = comment.List(ctx, tx, s.ID); err != nil {
			l.Log("err", err, "desc", "load of comments failed")
			return []Story{}, errInternal
		}
	}

	if err = tx.Commit(); err != nil {
		log.Log("err", err, "desc", "can't commit")
		return []Story{}, errInternal
	}

	return list, nil
}

// ListByCategory returns list of the stories for the selected
// category. Without paging for simplicity. Just limit by 1000 for
// demo purposes.
func ListByCategory(ctx context.Context, category string) ([]Story, error) {
	var (
		tx, err = storage.DB.BeginTxx(ctx, nil)
		l       = log.Fork().With("fn", "list-by-category")
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
WHERE s.category = $1
ORDER BY s.created_at DESC
LIMIT 1000`
		var rows *sqlx.Rows
		if rows, err = tx.QueryxContext(ctx, q, category); err != nil && err != sql.ErrNoRows {
			l.Log("err", err, "sql", q, "desc", "db select failed")
			return []Story{}, errInternal
		}
		if err == sql.ErrNoRows {
			return []Story{}, nil
		}
		defer rows.Close()
		for rows.Next() {
			var story Story
			if err = rows.StructScan(&story); err != nil {
				l.Log("err", err, "sql", q, "desc", "rows scan failed")
				return []Story{}, errInternal
			}

			story.Author = user.User{ID: story.CreatedBy, Login: story.AuthorName}
			list = append(list, story)
		}
	}

	// Add comments to the story list. They added in a separated
	// loop for avoid fetching nested query in the same
	// transaction.
	for i, s := range list {
		if list[i].Comments, err = comment.List(ctx, tx, s.ID); err != nil {
			l.Log("err", err, "desc", "load of comments failed")
			return []Story{}, errInternal
		}
	}

	if err = tx.Commit(); err != nil {
		log.Log("err", err, "desc", "can't commit")
		return []Story{}, errInternal
	}

	return list, nil
}

// ListByUser returns full list of the stories for the selected
// user. Without paging for simplicity. Just limit by 1000 for demo
// purposes.
func ListByUser(ctx context.Context, name string) ([]Story, error) {
	var (
		tx, err = storage.DB.BeginTxx(ctx, nil)
		l       = log.Fork().With("fn", "list-by-user")
	)
	if err != nil {
		l.Log("err", err, "desc", "can't begin transaction")
		return []Story{}, errInternal
	}
	defer tx.Rollback()

	// Check the user first.
	var u user.User
	{
		if u, err = user.GetByName(ctx, tx, name); err != nil && err != sql.ErrNoRows {
			return []Story{}, errInternal
		}
		if err == sql.ErrNoRows {
			return []Story{}, errors.New("user not found")
		}
	}

	// Retrieve the story from the storage.
	var list []Story
	{
		const q = `
SELECT id, title, kind, body, category, score, created_by, created_at
FROM story
WHERE created_by = $1
ORDER BY created_at DESC
LIMIT 1000`
		var rows *sqlx.Rows
		if rows, err = tx.QueryxContext(ctx, q, u.ID); err != nil && err != sql.ErrNoRows {
			l.Log("err", err, "sql", q, "desc", "db select failed")
			return []Story{}, errInternal
		}
		if err == sql.ErrNoRows {
			return []Story{}, nil
		}
		defer rows.Close()
		for rows.Next() {
			var story Story
			if err = rows.StructScan(&story); err != nil {
				l.Log("err", err, "sql", q, "desc", "rows scan failed")
				return []Story{}, errInternal
			}
			story.Author = user.User{ID: u.ID, Login: u.Login}
			list = append(list, story)
		}
	}

	// Add comments to the story list. They added in a separated
	// loop for avoid fetching nested query in the same
	// transaction.
	for i, s := range list {
		if list[i].Comments, err = comment.List(ctx, tx, s.ID); err != nil {
			l.Log("err", err, "desc", "load of comments failed")
			return []Story{}, errInternal
		}
	}

	if err = tx.Commit(); err != nil {
		log.Log("err", err, "desc", "can't commit")
		return []Story{}, errInternal
	}

	return list, nil
}
