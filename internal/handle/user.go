package handle

import (
	"context"
	"net/http"

	"reddit-clone-example/internal/user"

	"github.com/gin-gonic/gin"
)

type authRequest struct {
	Name string `json:"username"`
	Pass string `json:"password"`
}

func handleLogin(c *gin.Context) {
	var r authRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "incorrect arguments"})
		return
	}

	token, message := user.Login(context.Background(), r.Name, r.Pass)
	if token != "" {
		c.JSON(http.StatusOK, gin.H{"token": token})
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{"message": message})
}

func handleRegister(c *gin.Context) {
	var r authRequest
	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "incorrect arguments"})
		return
	}

	token, message := user.Register(context.Background(), r.Name, r.Pass)
	if token != "" {
		c.JSON(http.StatusOK, gin.H{"token": token})
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{"message": message})
}
