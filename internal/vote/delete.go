package vote

import (
	"context"

	"github.com/google/uuid"
)

// DeleteAll deletes all votes for the story. Not exposed to HTTP API.
func DeleteAll(ctx context.Context, storyID uuid.UUID) error {
	return nil
}
