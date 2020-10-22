package handle

import (
	"github.com/gin-gonic/gin"
)

func Route() *gin.Engine {
	var r = gin.Default()

	r.RedirectTrailingSlash = true
	r.StaticFile("index.html", "/web/index.html")
	r.StaticFile("favicon.ico", "/web/favicon.ico")
	r.Static("static", "/web")
	r.GET("/", handleRoot)
	r.GET("/a/:cat/:story_id", handleRoot)

	var api = r.Group("api")
	api.POST("register", headersCheck, handleRegister)
	api.POST("login", headersCheck, handleLogin)
	api.GET("user/:username", handleUserStories)
	api.GET("posts", handleStoryList)
	api.GET("posts/:catname", handleCategoryStories)
	api.GET("post/:story_id", handleGetStory)
	api.POST("posts", authCheck, handleCreateStory)
	api.DELETE("post/:story_id", authCheck, handleDeleteStory)
	api.POST("post/:story_id", authCheck, handleCreateComment)
	api.DELETE("post/:story_id/:comment_id", authCheck, handleDeleteComment)
	api.GET("post/:story_id/upvote", authCheck, )
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
