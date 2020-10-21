package handle

import (
	"context"
	"net/http"

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
	var h = c.GetHeader("Authorization")
	if h == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, msg("authentication header missed or not valid"))
		return
	}
	if !user.AuthCheck(context.Background(), h) {
		c.AbortWithStatusJSON(http.StatusUnauthorized, msg("wrong auth token"))
		return
	}
	c.Next()
}
