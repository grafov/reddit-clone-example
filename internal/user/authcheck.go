package user

import "context"

func AuthCheck(ctx context.Context, token string) bool {
	return true
}
