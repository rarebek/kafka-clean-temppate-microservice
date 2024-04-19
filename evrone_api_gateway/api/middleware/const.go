package middleware

type (
	ctxKeyRequestAuth int
)

const (
	RequestIDHeader                   = "X-Request-Id"
	RequestAuthCtx  ctxKeyRequestAuth = 0
)
