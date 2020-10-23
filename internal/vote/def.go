package vote

import (
	"errors"

	"github.com/google/uuid"
	"github.com/grafov/kiwi"
)

var (
	log         = kiwi.Fork().With("pkg", "vote")
	errInternal = errors.New("service error, try later")
)

type Vote struct {
	UserID uuid.UUID `json:"user" db:"account_id"`
	Value  int64     `json:"vote" db:"vote"`

	// Internal fields.
	StoryID uuid.UUID `json:"-" db:"story_id"`
}
