package auth

import (
	"encoding/json"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/anoriar/shortener/internal/shortener/domainerror"
	"github.com/anoriar/shortener/internal/shortener/dto/auth"
	"github.com/anoriar/shortener/internal/shortener/entity"
	"github.com/anoriar/shortener/internal/shortener/repository/user"
)

// Authenticator missing godoc.
type Authenticator struct {
	signService    *SignService
	userRepository user.UserRepositoryInterface
	logger         *zap.Logger
}

// NewAuthenticator missing godoc.
func NewAuthenticator(signService *SignService, userRepository user.UserRepositoryInterface, logger *zap.Logger) *Authenticator {
	return &Authenticator{signService: signService, userRepository: userRepository, logger: logger}
}

// GetToken missing godoc.
func (a *Authenticator) GetToken(srcToken string) (string, *auth.TokenPayload, error) {
	if srcToken == "" {

	}

	decodedToken, signature, err := a.signService.Decode(srcToken)
	if err != nil {
		a.logger.Error("decode token error", zap.String("error", err.Error()))
		return "", nil, domainerror.ErrInternal
	}
	if a.signService.Verify(decodedToken, signature) {
		tokenPayload := &auth.TokenPayload{}
		err = json.Unmarshal(decodedToken, tokenPayload)
		if err != nil {
			a.logger.Error("unmarshal token error", zap.String("error", err.Error()))
			return "", nil, domainerror.ErrInternal
		}

		_, exists, err := a.userRepository.FindUserByID(tokenPayload.UserID)
		if err != nil {
			a.logger.Error("find user error", zap.String("error", err.Error()))
			return "", nil, domainerror.ErrInternal
		}
		if exists {
			return "", tokenPayload, nil
		} else {
			return a.CreateNewToken()
		}

	} else {
		return a.CreateNewToken()
	}
}

// CreateNewToken missing godoc.
func (a *Authenticator) CreateNewToken() (string, *auth.TokenPayload, error) {
	tokenPayload := &auth.TokenPayload{UserID: uuid.NewString()}

	jsonTokenPayload, err := json.Marshal(tokenPayload)
	if err != nil {
		a.logger.Error("token payload marshal error", zap.String("error", err.Error()))
		return "", nil, domainerror.ErrInternal
	}
	newToken := a.signService.Sign([]byte(jsonTokenPayload))

	if err != nil {
		a.logger.Error("create newToken error", zap.String("error", err.Error()))
		return "", nil, domainerror.ErrInternal
	}
	err = a.userRepository.AddUser(entity.User{
		UUID: tokenPayload.UserID,
	})
	if err != nil {
		a.logger.Error("create user error", zap.String("error", err.Error()))
		return "", nil, domainerror.ErrInternal
	}

	return newToken, tokenPayload, nil
}
