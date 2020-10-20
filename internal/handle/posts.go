package handle

import (
	"github.com/gin-gonic/gin"
)

func handlePostList(c *gin.Context) {
	c.JSON(200, gin.H{"token": "zzz"})
}

func handlePost(c *gin.Context) {
	c.JSON(200, gin.H{"token": "zzz"})
}
