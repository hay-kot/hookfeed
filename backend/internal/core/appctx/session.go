package appctx

import "context"

const appctxKeySession appctxkey = "appctx:session"

type Session struct {
	Jwt   string // Raw JWT Token
	Token string // Stateful auth token
}

func WithSession(ctx context.Context, jwt string, token string) context.Context {
	return context.WithValue(ctx, appctxKeySession, Session{Jwt: jwt, Token: token})
}

// GetSession returns the embedded user form the context if available.
func GetSession(ctx context.Context) (user Session, ok bool) {
	user, ok = ctx.Value(appctxKeySession).(Session)
	return user, ok
}

// GetSessionMust is the same as GetSession but panics if ok is false.
func GetSessionMust(ctx context.Context) Session {
	user, ok := GetSession(ctx)
	if !ok {
		panic("invalid to get user from context")
	}

	return user
}
