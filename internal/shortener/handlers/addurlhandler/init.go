package addurlhandler

import (
	"github.com/anoriar/shortener/internal/shortener/config"
	"github.com/anoriar/shortener/internal/shortener/storage"
	"github.com/anoriar/shortener/internal/shortener/util"
)

func InitializeAddHandler(cnf *config.Config, storage storage.URLStorageInterface) *AddHandler {
	return NewAddHandler(storage, util.NewKeyGen(), cnf.BaseURL)
}
