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
	api.POST("register", handleRegister)
	api.POST("login", handleLogin)
	api.GET("posts", handlePostList)
	api.GET("post/:id", handlePost)

	return r
}

func handleRoot(c *gin.Context) {
	c.File("/web/index.html")
}
