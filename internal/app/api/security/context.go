package security

import "context"

const authenticatedUserKey = "authenticated:user"

func WithAuthenticatedUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, authenticatedUserKey, user)
}

func AuthenticatedUser(ctx context.Context) *User {
	user, ok := ctx.Value(authenticatedUserKey).(*User)

	if !ok {
		return nil
	}

	return user
}
