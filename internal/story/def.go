package story

import (
	"errors"
	"time"

	"reddit-clone-example/internal/comment"
	"reddit-clone-example/internal/user"
	"reddit-clone-example/internal/vote"

	"github.com/google/uuid"
	"github.com/grafov/kiwi"
)

var (
	log         = kiwi.Fork().With("pkg", "story")
	errInternal = errors.New("service error, try later")
)

const (
	textType = "text"
	linkType = "link"
)

// Story represents structure for keeping stories.
type Story struct {
	ID        uuid.UUID         `json:"id,omitempty" db:"id"`
	Author    user.User         `json:"author,omitempty"`
	Title     string            `json:"title" db:"title"`
	Type      string            `json:"type" db:"kind"`
	URL       string            `json:"url,omitempty"`
	Text      string            `json:"text,omitempty"`
	Category  string            `json:"category" db:"category"`
	Stat      int64             `json:"upvotePercentage,omitempty"`
	Score     int64             `json:"score"`
	Views     int64             `json:"views"`
	Votes     []vote.Vote       `json:"votes"`
	CreatedAt time.Time         `json:"created" db:"created_at"`
	Comments  []comment.Comment `json:"comments"`

	// Internal fields.
	CreatedBy  uuid.UUID `json:"-" db:"created_by"`
	AuthorName string    `json:"-" db:"login"`
	Body       string    `json:"-" db:"body"` // url or text stored to body
}

func (s *Story) MatchType() error {
	switch s.Type {
	case textType:
		s.Text = s.Body
	case linkType:
		s.URL = s.Body
	default:
		return errors.New("invalid story type")
	}
	return nil
}
