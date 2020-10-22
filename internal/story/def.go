package story

import (
	"errors"
	"time"

	"reddit-clone-example/internal/user"

	"github.com/google/uuid"
	"github.com/grafov/kiwi"
)

var (
	log         = kiwi.Fork().With("pkg", "story")
	errInternal = errors.New("internal error, try later")
)

const (
	textType = "text"
	linkType = "link"
)

type Story struct {
	ID        uuid.UUID `json:"id,omitempty" db:"id"`
	Author    user.User `json:"author,omitempty"`
	Title     string    `json:"title" db:"title"`
	Type      string    `json:"type" db:"kind"`
	URL       string    `json:"url,omitempty"`
	Text      string    `json:"text,omitempty"`
	Category  string    `json:"category" db:"category"`
	Stat      int64     `json:"upvotePercentage,omitempty"`
	Score     int64     `json:"score"`
	Views     int64     `json:"views"`
	Votes     []Vote    `json:"votes"`
	CreatedAt time.Time `json:"created" db:"created_at"`
	Comments  []string  `json:"comments"` // TODO

	CreatedBy  uuid.UUID `json:"-" db:"created_by"`
	AuthorName string    `json:"-" db:"login"`
	Body       string    `json:"-" db:"body"` // url or text stored to body
}

type Vote struct {
	User  uuid.UUID `json:"user"`
	Count int64     `json:"vote"`
}

func (s *Story) Upvote(user uuid.UUID) {
	for i, v := range s.Votes {
		if v.User == user {
			s.Votes[i].Count++
			s.Stat = 100 // TODO calculate this
			return
		}
	}
	s.Votes = append(s.Votes, Vote{user, 1})
	s.Stat = 100 // TODO calculate this
}
