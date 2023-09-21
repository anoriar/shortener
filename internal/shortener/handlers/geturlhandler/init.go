package geturlhandler

import (
	"github.com/anoriar/shortener/internal/shortener/storage"
)

func InitializeGetHandler(storage storage.URLStorageInterface) *GetHandler {
	return NewGetHandler(storage)
}
