package user

import (
	"context"

	"redditclone/storage"
)

func Register(ctx context.Context, name, pass string) (token, message string) {
	var (
		tx, err = storage.DB.BeginTxx(ctx, nil)
		l       = log.Fork().With("fn", "register")
	)
	if err != nil {
		l.Log("err", err, "desc", "can't begin transaction")
		return "", "can't register"
	}
	defer tx.Rollback()

	// проверить наличие юзера
	// вернуть токен

	if err = tx.Commit(); err != nil {
		log.Log("err", err, "desc", "can't commit")
		return "", "can't register"
	}

	return "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7InVzZXJuYW1lIjoiU2FtcGxlIiwiaWQiOiI1ZjhjOTUxMzYyNjg4NjAwMDgyZDQyZGYifSwiaWF0IjoxNjAzMTI2MDM1LCJleHAiOjE2MDM3MzA4MzV9.diily4DgCZI-eNiqycVJHnfkPB2HuMZb6WRuBvOyLk4", ""
}
