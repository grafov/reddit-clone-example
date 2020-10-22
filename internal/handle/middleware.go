package handle

import (
	"context"
	"net/http"
	"strings"

	"reddit-clone-example/internal/user"

	"github.com/gin-gonic/gin"
)

const session = "session" // key for user info in gin context

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
	// Parse the token from the header. Take into account that the token prepended by Bearer
	// keyword.
	var (
		token string
		err   error
	)
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

	// Pass auth data into gin context.
	var u *user.User
	{
		if u, err = user.AuthCheck(context.Background(), token); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, msg(err.Error()))
			return
		}
		c.Set(session, *u)
	}

	c.Next()
}

// getUser is not a middleware but helper for extracting user info
// from gin context.
func getUser(c *gin.Context) user.User {
	var u user.User
	if v, ok := c.Get(session); ok {
		u = v.(user.User)
	}
	return u
}
