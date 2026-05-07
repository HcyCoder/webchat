package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/team/webchat-server/common/token"
)

type ctxKey string

const UserIDKey ctxKey = "user_id"

func GetUserID(ctx context.Context) string {
	if v, ok := ctx.Value(UserIDKey).(string); ok {
		return v
	}
	return ""
}

type AuthMiddleware struct {
	tm *token.Manager
}

func NewAuthMiddleware(tm *token.Manager) *AuthMiddleware {
	return &AuthMiddleware{tm: tm}
}

func (m *AuthMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}
		userID, err := m.tm.Validate(r.Context(), auth[7:])
		if err != nil || userID == "" {
			http.Error(w, `{"error":"token expired"}`, http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), UserIDKey, userID)
		next(w, r.WithContext(ctx))
	}
}
