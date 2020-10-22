package handle

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"reddit-clone-example/internal/story"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func handleStoryList(c *gin.Context) {
	list, err := story.List(context.Background())
	if err != nil {
		c.JSON(http.StatusBadRequest, msg(err.Error()))
		return
	}

	c.JSON(http.StatusOK, list)
}

func handleUserStories(c *gin.Context) {
	// Parse and validate args.
	var (
		name string
		err  error
	)
	if name = strings.TrimSpace(c.Param("name")); name == "" {
		c.JSON(http.StatusBadRequest, msg("empty user name"))
		return
	}

	// Retrieve the story.
	var s []story.Story
	if s, err = story.ListByUser(context.Background(), name); err != nil {
		c.JSON(http.StatusBadRequest, msg(err.Error()))
		return
	}

	c.JSON(http.StatusOK, s)
}

func handleCreateStory(c *gin.Context) {
	// Parse and validate input args.
	var (
		s   story.Story
		err error
	)
	{

		if err = c.ShouldBindJSON(&s); err != nil {
			c.JSON(http.StatusBadRequest, msg("invalid JSON request"))
			return
		}
		switch s.Type {
		case "text":
			s.URL = ""
		case "link":
			var u *url.URL
			u, err = url.Parse(s.URL)
			if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
				c.JSON(http.StatusBadRequest, msg("url is invalid"))
				return
			}
			s.Text = ""
		default:
			c.JSON(http.StatusBadRequest, msg("invalid story type"))
			return

		}
	}

	if s, err = story.Create(context.Background(), getUser(c), s); err != nil {
		c.JSON(http.StatusBadRequest, msg(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, s)
}

func handleGetStory(c *gin.Context) {
	// Parse and validate args.
	var (
		id  uuid.UUID
		err error
	)
	if id, err = uuid.Parse(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, msg("invalid story id"))
		return
	}

	// Retrieve the story.
	var s story.Story
	if s, err = story.Get(context.Background(), id); err != nil {
		c.JSON(http.StatusBadRequest, msg(err.Error()))
		return
	}

	c.JSON(http.StatusOK, s)
}

func handleDeleteStory(c *gin.Context) {
	var (
		id  uuid.UUID
		err error
	)
	if id, err = uuid.Parse(c.Param("id")); err != nil {
		c.JSON(http.StatusBadRequest, msg("invalid story id"))
		return
	}
	if err = story.Delete(context.Background(), getUser(c), id); err != nil {
		c.JSON(http.StatusBadRequest, msg("can't delete"))
		return
	}

	c.JSON(http.StatusOK, msg("success"))
}
