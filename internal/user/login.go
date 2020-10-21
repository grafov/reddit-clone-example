package user

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"redditclone/internal/config"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// Login sign in user into system by his name and pass. It returns
// token on success or message on failure.
func Login(ctx context.Context, name, pass string) (token, message string) {
	var (
		l   = log.Fork().With("fn", "login")
		err error
	)

	const q = `SELECT id, login FROM account WHERE login = $1 AND pass = $2`

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": name,
		"id":       "123",
	})
	if token, err = t.SignedString([]byte(config.App.TokenSecret)); err != nil {
		l.Log("err", err, "desc", "fail to sign token")
		message = "auth failed"
	}

	//	token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjp7InVzZXJuYW1lIjoiU2FtcGxlIiwiaWQiOiI1ZjhjOTUxMzYyNjg4NjAwMDgyZDQyZGYifSwiaWF0IjoxNjAzMTI2MDM1LCJleHAiOjE2MDM3MzA4MzV9.diily4DgCZI-eNiqycVJHnfkPB2HuMZb6WRuBvOyLk4"

	return
}

func auth(c *gin.Context) {
	var h = c.GetHeader("Authorization")
	const hmacSampleSecret = "replace-this-sample"
	const BEARER_SCHEMA = "Bearer"

	if len(h) < 8 {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized, errors.New("access denied: wrong auth token"),
		)
		return
	}

	token, err := validateToken(h[len(BEARER_SCHEMA):])
	if token.Valid {
		claims := token.Claims.(jwt.MapClaims)
		fmt.Println(claims)
	} else {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	c.Next()
}

func validateToken(encodedToken string) (*jwt.Token, error) {
	return jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
		if _, isvalid := token.Method.(*jwt.SigningMethodHMAC); !isvalid {
			return nil, fmt.Errorf("invalid token", token.Header["alg"])

		}
		return []byte("sessid"), nil
	})

}
