package geturlhandler

import "github.com/anoriar/shortener/internal/storage"

func InitializeGetHandler() *GetHandler {
	return NewGetHandler(storage.GetInstance())
}
