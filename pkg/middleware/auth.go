package middleware

import (
	"context"
	//"ecom-api/pkg/auth"
	"net/http"
	"strings"

	"github.com/babishagetaneh1992/ecom-api/pkg/auth"
)

type contextKey string

const (
	userCtxKey = contextKey("userID")
	roleCtxKey = contextKey("role")
)

// AuthMiddleware checks JWT and attaches userID into request context
func AuthMiddleware(next http.Handler) http.Handler {
	return  http.HandlerFunc(func (w http.ResponseWriter, r *http.Request)  {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return 
		}

		// format: "Bearer token"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return 
		}

		claims, err := auth.VerifyToken(parts[1])
		if err != nil {
			http.Error(w, "Invalid Token:" +err.Error(), http.StatusUnauthorized)
			return 
		}

		// injecting userID in to context
		ctx := context.WithValue(r.Context(), userCtxKey, claims.UserID)
        ctx = context.WithValue(ctx, roleCtxKey, claims.Role) // could be ""

		next.ServeHTTP(w, r.WithContext(ctx))
	})

	
}

func AdminOnly(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        _, role := FromContext(r.Context())
        if role != "admin" {
            http.Error(w, "Forbidden: Admins only", http.StatusForbidden)
            return
        }
        next.ServeHTTP(w, r)
    })
}



// retrive userID from request context
func FromContext(ctx context.Context) (string, string) {
    uid, _ := ctx.Value(userCtxKey).(string)
    role, ok := ctx.Value(roleCtxKey).(string)
	if !ok || role == "" {
		role = "user"
	}
    return uid, role
}