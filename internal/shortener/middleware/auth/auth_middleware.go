package auth

import (
	"context"
	"errors"
	"net/http"

	context2 "github.com/anoriar/shortener/internal/shortener/context"
	v1 "github.com/anoriar/shortener/internal/shortener/services/auth"
)

const cookieName = "token"
const cookieAge = 3600

// AuthMiddleware missing godoc.
type AuthMiddleware struct {
	authenticator *v1.Authenticator
}

// NewAuthMiddleware missing godoc.
func NewAuthMiddleware(authenticator *v1.Authenticator) *AuthMiddleware {
	return &AuthMiddleware{authenticator: authenticator}
}

// Auth missing godoc.
func (am *AuthMiddleware) Auth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		ctx := request.Context()

		shouldCreateNewToken := false
		authCookie, err := request.Cookie(cookieName)
		if err != nil || authCookie.Value == "" {
			if errors.Is(err, http.ErrNoCookie) {
				shouldCreateNewToken = true
			} else {
				http.Error(w, "can not authorized", http.StatusUnauthorized)
				return
			}
		} else {
			isVerified, tokenPayload, err := am.authenticator.GetToken(authCookie.Value)
			if err != nil {
				http.Error(w, "get token error", http.StatusInternalServerError)
				return
			}
			if !isVerified {
				shouldCreateNewToken = true
			} else {
				ctx = context.WithValue(request.Context(), context2.UserIDContextKey, tokenPayload.UserID)
			}
		}

		if shouldCreateNewToken {
			token, tokenPayload, err := am.authenticator.CreateNewToken()
			if err != nil {
				http.Error(w, "create token error", http.StatusInternalServerError)
				return
			}

			ctx = context.WithValue(request.Context(), context2.UserIDContextKey, tokenPayload.UserID)
			http.SetCookie(w, &http.Cookie{
				Name:     cookieName,
				Value:    token,
				Path:     "/",
				MaxAge:   cookieAge,
				HttpOnly: true,
			})
		}

		h.ServeHTTP(w, request.WithContext(ctx))
	})
}
