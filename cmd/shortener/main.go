package main

import (
	"github.com/anoriar/shortener/internal/shortener/config"
	database "github.com/anoriar/shortener/internal/shortener/db"
	"github.com/anoriar/shortener/internal/shortener/logger"
	"github.com/anoriar/shortener/internal/shortener/router"
	"github.com/caarlos0/env/v6"
	"go.uber.org/zap"
	"net/http"
)

func main() {
	run()
}

func run() {
	conf := config.NewConfig()
	parseFlags(conf)

	err := env.Parse(conf)
	if err != nil {
		panic(err)
	}

	logger, err := logger.Initialize(conf.LogLevel)
	if err != nil {
		panic(err)
	}

	defer logger.Sync()

	db, err := database.InitializeDatabase(conf.DatabaseDSN)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	r, err := router.InitializeRouter(conf, logger, db)

	if err != nil {
		logger.Fatal("init error", zap.String("error", err.Error()))
		panic(err)
	}

	err = http.ListenAndServe(conf.Host, r.Route())
	if err != nil {
		logger.Fatal("Server exception", zap.String("exception", err.Error()))
		panic(err)
	}
}
