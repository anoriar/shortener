package auth

// TokenPayload payload токен авторизации пользователя
type TokenPayload struct {
	UserID string `json:"user_id"` // id пользователя на сайте
}
