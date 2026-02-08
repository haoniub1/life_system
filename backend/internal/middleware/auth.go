package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/zeromicro/go-zero/rest/httpx"
	"life-system-backend/internal/types"
)

const UserIDKey = "userID"

type CustomClaims struct {
	UserID int64 `json:"userId"`
	jwt.RegisteredClaims
}

func AuthMiddleware(secret string) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			var token string
			var tokenSource string

			// Try Authorization header first (more reliable with localStorage)
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" {
				parts := strings.SplitN(authHeader, " ", 2)
				if len(parts) == 2 && parts[0] == "Bearer" {
					token = parts[1]
					tokenSource = "Authorization header"
				}
			}

			// Fallback to cookie
			if token == "" {
				if cookie, err := r.Cookie("token"); err == nil && cookie.Value != "" {
					token = cookie.Value
					tokenSource = "cookie"
				}
			}

			fmt.Printf("🔐 Auth check for %s - Token source: %s, Token present: %v, Secret length: %d\n",
				r.URL.Path, tokenSource, token != "", len(secret))

			if token == "" {
				fmt.Printf("❌ No token found\n")
				httpx.OkJson(w, types.CommonResp{
					Code:    401,
					Message: "unauthorized",
				})
				return
			}

			fmt.Printf("🔍 Token (first 20 chars): %s...\n", token[:min(20, len(token))])

			claims := &CustomClaims{}
			parsedToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(secret), nil
			})

			if err != nil {
				fmt.Printf("❌ Token parse error: %v\n", err)
				httpx.OkJson(w, types.CommonResp{
					Code:    401,
					Message: "invalid token",
				})
				return
			}

			if !parsedToken.Valid {
				fmt.Printf("❌ Token is not valid\n")
				httpx.OkJson(w, types.CommonResp{
					Code:    401,
					Message: "invalid token",
				})
				return
			}

			fmt.Printf("✅ Token valid, UserID: %d\n", claims.UserID)

			// Inject userID into context
			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func GetUserID(ctx context.Context) (int64, error) {
	userID, ok := ctx.Value(UserIDKey).(int64)
	if !ok {
		return 0, fmt.Errorf("user id not found in context")
	}
	return userID, nil
}
