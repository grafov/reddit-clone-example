package user

import (
	"context"
	"errors"
	"fmt"
	"redditclone/internal/config"

	"github.com/dgrijalva/jwt-go"
	"github.com/grafov/kiwi"
)

// AuthCheck checks for JWT token validity.
func AuthCheck(ctx context.Context, authkey string) (*Authbox, error) {
	secret := func(token *jwt.Token) (interface{}, error) {
		method, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok || method.Alg() != "HS256" {
			return nil, fmt.Errorf("bad sign method")
		}
		// Note that token secret should be passed as []byte
		// here. Refs to issues:
		// https://github.com/dgrijalva/jwt-go/issues/65
		// https://github.com/dgrijalva/jwt-go/issues/147
		// https://github.com/dgrijalva/jwt-go/issues/223
		return []byte(config.App.TokenSecret), nil
	}

	var (
		token *jwt.Token
		err   error
	)
	{
		if token, err = jwt.ParseWithClaims(authkey, &Authbox{}, secret); err != nil {
			kiwi.Log("err", err, "token", token)
			return nil, err
		}
		if !token.Valid {
			kiwi.Log("token", token)
			return nil, errors.New("bad token")
		}
	}

	var auth *Authbox
	{
		var ok bool
		if auth, ok = token.Claims.(*Authbox); !ok {
			return nil, errors.New("no payload in the token")
		}
		kiwi.Log("payload", auth)

	}

	return auth, nil
}
