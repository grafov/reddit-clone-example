package handle

import (
	"context"
	"errors"
	"net/http"

	"redditclone/internal/user"

	"github.com/gin-gonic/gin"
)

// headers middleware checks for valid content type for API requests
func headers(c *gin.Context) {
	if c.GetHeader("Content-Type") != "application/json" {
		c.AbortWithStatusJSON(
			http.StatusBadRequest, errors.New("request payload not recognized"),
		)
		return
	}
	c.Next()
}

// auth middleware checks for authorization header
func auth(c *gin.Context) {
	var h = c.GetHeader("Authorization")

	if len(h) < 8 {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized, errors.New("access denied: badly formatted auth token"),
		)
		return
	}

	if !user.AuthCheck(context.Background(), h) {
		c.AbortWithStatusJSON(
			http.StatusUnauthorized, errors.New("access denied: wrong auth token"),
		)
		return
	}

	c.Next()
}
