package comment

import (
	"errors"
	"time"

	"reddit-clone-example/internal/user"

	"github.com/google/uuid"
	"github.com/grafov/kiwi"
)

var (
	log         = kiwi.Fork().With("pkg", "comment")
	errInternal = errors.New("internal error, try later")
)

// Comment represents structure for keeping comments for stories.
type Comment struct {
	ID        uuid.UUID `json:"id" db:"id"`
	StoryID   uuid.UUID `json:"story_id" db:"story_id"`
	Body      string    `json:"comment" db:"body"`
	Author    user.User `json:"author,omitempty"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`

	// internal fields
	CreatedBy uuid.UUID `json:"-" db:"created_by"`
}
