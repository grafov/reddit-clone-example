package story

import (
	"context"

	"reddit-clone-example/internal/user"

	"github.com/google/uuid"
)

func Delete(ctx context.Context, owner user.User, id uuid.UUID) error {
	return nil
}
