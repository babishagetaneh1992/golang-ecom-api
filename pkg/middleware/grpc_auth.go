package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/babishagetaneh1992/ecom-api/services/user-ms/adaptors/grpc/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// GRPCAuthMiddleware is an HTTP middleware that calls user-ms VerifyToken
func GRPCAuthMiddleware(client pb.UserServiceClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
				return
			}

			token := parts[1]
			resp, err := client.VerifyToken(r.Context(), &pb.VerifyTokenRequest{Token: token})
			if err != nil || !resp.Valid {
				http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
				return
			}

			// Add userID and role to context using same keys as AuthMiddleware in auth.go
			ctx := context.WithValue(r.Context(), userCtxKey, resp.UserId)
			ctx = context.WithValue(ctx, roleCtxKey, resp.Role)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// UnaryGRPCAuthInterceptor is a gRPC interceptor that calls user-ms VerifyToken
func UnaryGRPCAuthInterceptor(client pb.UserServiceClient) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "Missing metadata")
		}

		authHeader := md["authorization"]
		if len(authHeader) == 0 {
			return nil, status.Error(codes.Unauthenticated, "Missing authorization header")
		}

		parts := strings.Split(authHeader[0], " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return nil, status.Error(codes.Unauthenticated, "Invalid authorization header format")
		}

		resp, err := client.VerifyToken(ctx, &pb.VerifyTokenRequest{Token: parts[1]})
		if err != nil || !resp.Valid {
			return nil, status.Errorf(codes.Unauthenticated, "Invalid token")
		}

		// Add userID to context
		ctx = context.WithValue(ctx, userCtxKey, resp.UserId)
		ctx = context.WithValue(ctx, roleCtxKey, resp.Role)

		return handler(ctx, req)
	}
}
