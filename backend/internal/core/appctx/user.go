package appctx

import (
	"context"

	"github.com/google/uuid"
)

// appctxkey is the type used to avoid colisions in the context package.
type appctxkey string

const appctxKeyUser appctxkey = "appctx:user"

// User is the base identity for a user. This can be used as a contract throughout the
// application for managing a user identity and associated permission lookups.
//
// The intent of this [User] type is to reduce excessive arguments for common fields
// across different packages and gives them a base to agree on an "identity" for a user.
//
// This creates a common pattern for methods to look like this
//
//	func GetUserFromDatabase(ctx context.Context, identity appctx.User, ...)
type User struct {
	ID uuid.UUID // User ID from the database
}

// WithUser stores the User within the context. User is assumed to be the identity
// of the requestor.
func WithUser(ctx context.Context, user User) context.Context {
	return context.WithValue(ctx, appctxKeyUser, user)
}

// GetUser returns the embedded user form the context if available.
func GetUser(ctx context.Context) (user User, ok bool) {
	user, ok = ctx.Value(appctxKeyUser).(User)
	return user, ok
}

// GetUserMust is the same as GetUser but panics if ok is false.
func GetUserMust(ctx context.Context) User {
	user, ok := GetUser(ctx)
	if !ok {
		panic("invalid to get user from context")
	}

	return user
}
