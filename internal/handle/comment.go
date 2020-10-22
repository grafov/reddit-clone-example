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
		s   story.Story
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

	// Really it returns whole comment on creation but API wants
	// display parent story with all related comments.
	if _, err = comment.Create(context.Background(), getUser(c), com); err != nil {
		c.JSON(http.StatusBadRequest, msg(err.Error()))
		return
	}

	if s, err = story.Get(context.Background(), com.StoryID); err != nil {
		c.JSON(http.StatusBadRequest, msg(err.Error()))
		return
	}

	c.JSON(http.StatusCreated, s)
}

func handleDeleteComment(c *gin.Context) {
	var (
		storyID, commentID uuid.UUID
		err                error
	)
	if storyID, err = uuid.Parse(c.Param("story_id")); err != nil {
		c.JSON(http.StatusBadRequest, msg("invalid story id"))
		return
	}
	if commentID, err = uuid.Parse(c.Param("comment_id")); err != nil {
		c.JSON(http.StatusBadRequest, msg("invalid comment id"))
		return
	}
	if err = comment.Delete(context.Background(), getUser(c), storyID, commentID); err != nil {
		c.JSON(http.StatusBadRequest, msg("can't delete"))
		return
	}

	c.JSON(http.StatusOK, msg("success"))
}
