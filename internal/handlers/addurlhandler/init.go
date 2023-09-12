package addurlhandler

import (
	"github.com/anoriar/shortener/internal/config"
	"github.com/anoriar/shortener/internal/storage"
	"github.com/anoriar/shortener/internal/util"
)

func InitializeAddHandler(cnf *config.Config) *AddHandler {
	return NewAddHandler(storage.GetInstance(), util.NewKeyGen(), cnf.BaseURL)
}
