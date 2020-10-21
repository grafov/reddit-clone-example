package handle

import (
	"context"
	"net/http"
	"strings"

	"redditclone/internal/user"

	"github.com/gin-gonic/gin"
)

// headers middleware checks for valid content type for API requests
func headers(c *gin.Context) {
	if c.GetHeader("Content-Type") != "application/json" {
		c.AbortWithStatusJSON(http.StatusBadRequest, msg("request payload not recognized"))
		return
	}
	c.Next()
}

// auth middleware checks for authorization header
func auth(c *gin.Context) {
	var (
		token string
		err   error
	)
	// Parse the token from the header. Take into account that the token prepended by Bearer
	// keyword.
	{
		h := c.GetHeader("Authorization")
		if len(h) < 8 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, msg("authorization header missed or not valid"))
			return
		}
		s := strings.SplitN(h, "Bearer ", 2)
		if len(s) < 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, msg("badly formatted authorization header (Bearer missed)"))
			return
		}
		token = s[1]
	}

	var a *user.Authbox
	if a, err = user.AuthCheck(context.Background(), token); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, msg(err.Error()))
		return
	}
	c.Set("auth", a)

	c.Next()
}
