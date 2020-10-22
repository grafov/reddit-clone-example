package handle

import (
	"context"
	"net/http"

	"reddit-clone-example/internal/comment"
	"reddit-clone-example/internal/story"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func handleCreateComment(c *gin.Context) {
	// Parse and validate input args.
	var (
		com comment.Comment
		err error
	)
	{
		var id uuid.UUID
		if id, err = uuid.Parse(c.Param("story_id")); err != nil {
			c.JSON(http.StatusBadRequest, msg("invalid story id"))
			return
		}
		if err = c.ShouldBindJSON(&com); err != nil {
			c.JSON(http.StatusBadRequest, msg("invalid JSON request"))
			return
		}
		com.StoryID = id
	}

	if err = comment.Create(context.Background(), getUser(c), com); err != nil {
		c.JSON(http.StatusBadRequest, msg(err.Error()))
		return
	}

	// Frontend wants to reload the story with all the comments
	// just after comment creation.
	var s story.Story
	if s, err = story.Get(context.Background(), com.StoryID); err != nil {
		c.JSON(http.StatusBadRequest, msg(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, s)
}

func handleDeleteComment(c *gin.Context) {
	// Parse and validate input args.
	var (
		storyID, commentID uuid.UUID
		err                error
	)
	{
		if storyID, err = uuid.Parse(c.Param("story_id")); err != nil {
			c.JSON(http.StatusBadRequest, msg("invalid story id"))
			return
		}
		if commentID, err = uuid.Parse(c.Param("comment_id")); err != nil {
			c.JSON(http.StatusBadRequest, msg("invalid comment id"))
			return
		}

	}

	if err = comment.Delete(context.Background(), getUser(c), storyID, commentID); err != nil {
		c.JSON(http.StatusBadRequest, msg("can't delete"))
		return
	}

	// Frontend wants to reload the story with all the comments
	// just after comment deletion.
	var s story.Story
	if s, err = story.Get(context.Background(), storyID); err != nil {
		c.JSON(http.StatusBadRequest, msg(err.Error()))
		return
	}

	c.JSON(http.StatusOK, s)
}
