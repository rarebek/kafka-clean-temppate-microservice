package middleware

import (
	"net/http"

	"github.com/casbin/casbin/v2"
	"go.uber.org/zap"
)

func Authorizer(e *casbin.CachedEnforcer, logger *zap.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			data, ok := r.Context().Value(RequestAuthCtx).(map[string]string)
			if !ok {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			ok, err := e.Enforce(data["sub"], r.URL.Path, r.Method)
			if err != nil {
				logger.Error("middleware authorizer",
					zap.Error(err),
					zap.String("sub", data["sub"]),
					zap.String("path", r.URL.Path),
					zap.String("method", r.Method),
				)
			}
			if ok {
				next.ServeHTTP(w, r)
				return
			}
			http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
		})
	}
}
