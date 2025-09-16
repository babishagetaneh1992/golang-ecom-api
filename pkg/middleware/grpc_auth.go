package middleware

import (
	"context"
	"strings"

	//"ecom-api/pkg/auth"
	"github.com/babishagetaneh1992/ecom-api/pkg/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func UnaryAuthInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "Missing metadata")
	}

	// Expect "authorization: Bearer <token>"
	authHeader := md["authorization"]
	if len(authHeader) == 0 {
		return nil, status.Error(codes.Unauthenticated, "Missing authorization header")
	}

	parts := strings.Split(authHeader[0], " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, status.Error(codes.Unauthenticated, "Invalid authorization header format")
	}

	claims, err := auth.VerifyToken(parts[1])
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "Invalid token: %v", err)
	}

	// Add userID to context
	ctx = context.WithValue(ctx, "userID", claims.UserID)

	// Continue
	return handler(ctx, req)
}
