package grpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	context2 "github.com/anoriar/shortener/internal/shortener/context"
	"github.com/anoriar/shortener/internal/shortener/services/auth"
)

const metadataTokenName = "token"

// AuthInterceptor missing godoc.
type AuthInterceptor struct {
	authenticator *auth.Authenticator
}

// NewAuthInterceptor missing godoc.
func NewAuthInterceptor(authenticator *auth.Authenticator) *AuthInterceptor {
	return &AuthInterceptor{authenticator: authenticator}
}

// Auth missing godoc.
func (ai *AuthInterceptor) Auth(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	shouldCreateNewToken := false

	var token string
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		values := md.Get(metadataTokenName)
		if len(values) > 0 {
			token = values[0]
		}
	}
	if len(token) == 0 {
		shouldCreateNewToken = true
	} else {
		isVerified, tokenPayload, err := ai.authenticator.GetToken(token)
		if err != nil {
			return nil, status.Error(codes.Internal, "internal error")
		}
		if !isVerified {
			shouldCreateNewToken = true
		} else {
			ctx = context.WithValue(ctx, context2.UserIDContextKey, tokenPayload.UserID)
		}

	}

	if shouldCreateNewToken {
		newToken, tokenPayload, err := ai.authenticator.CreateNewToken()
		if err != nil {
			return nil, status.Error(codes.Internal, "internal error")
		}
		ctx = context.WithValue(ctx, context2.UserIDContextKey, tokenPayload.UserID)
		token = newToken

	}

	resp, err := handler(ctx, req)
	if err != nil {
		return resp, err
	}

	md := metadata.Pairs(metadataTokenName, token)
	err = grpc.SendHeader(ctx, md)
	if err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}
	return resp, nil
}