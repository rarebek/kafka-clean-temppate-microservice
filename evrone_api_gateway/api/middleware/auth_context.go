package middleware

import (
	"context"
	"net/http"

	token_pkg "evrone_api_gateway/internal/pkg/token"
)

func AuthContext(jwtsecret string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var (
				token    string
				authData = make(map[string]string)
			)
			if token = r.URL.Query().Get("token"); len(token) == 0 {
				if token = r.Header.Get("Authorization"); len(token) > 10 {
					token = token[7:]
				}
			}
			claims, err := token_pkg.ParseJwtToken(token, jwtsecret)
			if err == nil && len(claims) != 0 {
				for key, value := range claims {
					if valStr, ok := value.(string); ok {
						authData[key] = valStr
					}
				}
				next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), RequestAuthCtx, authData)))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
