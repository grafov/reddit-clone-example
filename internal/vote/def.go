package vote

import (
	"errors"

	"github.com/google/uuid"
	"github.com/grafov/kiwi"
)

var (
	log         = kiwi.Fork().With("pkg", "vote")
	errInternal = errors.New("internal error, try later")
)

type Vote struct {
	User  uuid.UUID `json:"user"`
	Count int64     `json:"vote"`
}
