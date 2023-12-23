package auth

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"

	context2 "github.com/anoriar/shortener/internal/shortener/context"
	"github.com/anoriar/shortener/internal/shortener/dto/auth"
	"github.com/anoriar/shortener/internal/shortener/entity"
	"github.com/anoriar/shortener/internal/shortener/repository/user"
	v1 "github.com/anoriar/shortener/internal/shortener/services/auth"
)

const cookieName = "token"
const cookieAge = 3600

// AuthMiddleware missing godoc.
type AuthMiddleware struct {
	signService    *v1.SignService
	userRepository user.UserRepositoryInterface
}

// NewAuthMiddleware missing godoc.
func NewAuthMiddleware(signService *v1.SignService, userRepository user.UserRepositoryInterface) *AuthMiddleware {
	return &AuthMiddleware{signService: signService, userRepository: userRepository}
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
			decodedToken, signature, err := am.signService.Decode(authCookie.Value)
			if err != nil {
				http.Error(w, "decode token error", http.StatusInternalServerError)
				return
			}
			if am.signService.Verify(decodedToken, signature) {
				tokenPayload := &auth.TokenPayload{}
				err = json.Unmarshal(decodedToken, tokenPayload)
				if err != nil {
					http.Error(w, "unmarshal token error", http.StatusInternalServerError)
					return
				}

				_, exists, err := am.userRepository.FindUserByID(tokenPayload.UserID)
				if err != nil {
					http.Error(w, "find user error", http.StatusInternalServerError)
					return
				}
				if exists {
					//#MENTOR: Лучше не передавать переменные через контекст, но тут пришлось
					// Есть ли более хорошее решение, как можно передать userID в хендлеры?
					ctx = context.WithValue(request.Context(), context2.UserIDContextKey, tokenPayload.UserID)
				} else {
					shouldCreateNewToken = true
				}

			} else {
				shouldCreateNewToken = true
			}
		}

		if shouldCreateNewToken {
			tokenPayload := auth.TokenPayload{UserID: uuid.NewString()}
			token, err := am.createNewToken(tokenPayload)
			if err != nil {
				http.Error(w, "create token error", http.StatusInternalServerError)
				return
			}
			err = am.userRepository.AddUser(entity.User{
				UUID: tokenPayload.UserID,
			})
			if err != nil {
				http.Error(w, "create user error", http.StatusInternalServerError)
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

func (am *AuthMiddleware) createNewToken(payload auth.TokenPayload) (string, error) {

	jsonTokenPayload, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	token := am.signService.Sign([]byte(jsonTokenPayload))
	return token, nil
}
