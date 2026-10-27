package middlewares

import (
	"context"
	"net/http"
	"strconv"

	"github.com/deside01/hotel_notebook/backend/internal/utils"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(secretKey []byte) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("token")
			if err != nil {
				utils.Err(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			claims := &jwt.RegisteredClaims{}
			token, err := jwt.ParseWithClaims(cookie.Value, claims,
				func(token *jwt.Token) (any, error) {
					return secretKey, nil
				},
			)

			if err != nil || !token.Valid {
				utils.Err(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			userID, err := strconv.Atoi(claims.Subject)
			if err != nil {
				utils.Err(w, http.StatusUnauthorized, "unauthorized")
				return
			}

			ctx := context.WithValue(r.Context(), "user_id", userID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
