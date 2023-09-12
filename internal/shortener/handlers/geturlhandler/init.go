package geturlhandler

import (
	"github.com/anoriar/shortener/internal/shortener/storage"
)

func InitializeGetHandler() *GetHandler {
	return NewGetHandler(storage.GetInstance())
}
