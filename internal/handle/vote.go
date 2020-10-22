package handle

import (
	"context"
	"net/http"

	"reddit-clone-example/internal/story"
	"reddit-clone-example/internal/vote"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func handleUpvote(c *gin.Context) {
	var (
		id  uuid.UUID
		err error
	)
	if id, err = uuid.Parse(c.Param("story_id")); err != nil {
		c.JSON(http.StatusBadRequest, msg("invalid story id"))
		return
	}

	var votes []vote.Vote
	if votes, err = vote.Upvote(context.Background(), getUser(c), id); err != nil {
		c.JSON(http.StatusBadRequest, msg(err.Error()))
		return
	}

	// Frontend wants to reload the story with all the comments
	// just after comment creation.
	var s story.Story
	{
		if s, err = story.Get(context.Background(), id); err != nil {
			c.JSON(http.StatusBadRequest, msg(err.Error()))
			return
		}
		s.Votes = votes
	}

	c.JSON(http.StatusCreated, s)
}

func handleUnvote(c *gin.Context) {
	var (
		id  uuid.UUID
		err error
	)
	if id, err = uuid.Parse(c.Param("story_id")); err != nil {
		c.JSON(http.StatusBadRequest, msg("invalid story id"))
		return
	}

	var votes []vote.Vote
	if votes, err = vote.Neutral(context.Background(), getUser(c), id); err != nil {
		c.JSON(http.StatusBadRequest, msg(err.Error()))
		return
	}

	// Frontend wants to reload the story with all the comments
	// just after comment creation.
	var s story.Story
	{
		if s, err = story.Get(context.Background(), id); err != nil {
			c.JSON(http.StatusBadRequest, msg(err.Error()))
			return
		}
		s.Votes = votes
	}

	c.JSON(http.StatusCreated, s)
}

func handleDownvote(c *gin.Context) {
	var (
		id  uuid.UUID
		err error
	)
	if id, err = uuid.Parse(c.Param("story_id")); err != nil {
		c.JSON(http.StatusBadRequest, msg("invalid story id"))
		return
	}

	var votes []vote.Vote
	if votes, err = vote.Downvote(context.Background(), getUser(c), id); err != nil {
		c.JSON(http.StatusBadRequest, msg(err.Error()))
		return
	}

	// Frontend wants to reload the story with all the comments
	// just after comment creation.
	var s story.Story
	{
		if s, err = story.Get(context.Background(), id); err != nil {
			c.JSON(http.StatusBadRequest, msg(err.Error()))
			return
		}
		s.Votes = votes
	}

	c.JSON(http.StatusCreated, s)
}
