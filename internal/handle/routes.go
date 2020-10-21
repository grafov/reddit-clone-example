package handle

import (
	"github.com/gin-gonic/gin"
)

func Route() *gin.Engine {
	var r = gin.Default()

	r.RedirectTrailingSlash = true
	r.GET("/", handleRoot)
	r.StaticFile("index.html", "/web/index.html")
	r.StaticFile("favicon.ico", "/web/favicon.ico")
	r.Static("static", "/web")

	var api = r.Group("api")
	api.POST("register", headers, handleRegister)
	api.POST("login", headers, handleLogin)
	api.GET("user/:name", handleProfile)
	api.GET("posts", handleStoryList)
	api.GET("post/:id", handleGetStory)
	api.POST("posts", auth, handleCreateStory)
	api.DELETE("post/:id", auth, handleDeleteStory)

	return r
}

func handleRoot(c *gin.Context) {
	c.File("/web/index.html")
}

func msg(text string) map[string]string {
	m := make(map[string]string)
	m["message"] = text
	return m
}
