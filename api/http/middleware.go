// JWT 校验中间件
package http

import (
	"net/http"
	"strings"

	"github.com/fire-disposal/health_DT_go/internal/auth"
	"github.com/golang-jwt/jwt/v4"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "未提供Token", http.StatusUnauthorized)
			return
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.ParseWithClaims(tokenStr, &auth.Claims{}, func(token *jwt.Token) (interface{}, error) {
			return auth.JwtSecret(), nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Token无效", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
