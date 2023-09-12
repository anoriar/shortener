package addurlhandler

import (
	"github.com/anoriar/shortener/internal/shortener/config"
	"github.com/anoriar/shortener/internal/shortener/storage"
	"github.com/anoriar/shortener/internal/shortener/util"
)

func InitializeAddHandler(cnf *config.Config) *AddHandler {
	return NewAddHandler(storage.GetInstance(), util.NewKeyGen(), cnf.BaseURL)
}
