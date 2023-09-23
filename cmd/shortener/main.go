package main

import (
	"github.com/anoriar/shortener/internal/shortener/config"
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
		return
	}

	logger, err := logger.Initialize(conf.LogLevel)

	defer logger.Sync()

	r := router.InitializeRouter(conf, logger)

	err = http.ListenAndServe(conf.Host, r.Route())
	if err != nil {
		logger.Fatal("Server error", zap.String("error", err.Error()))
		panic(err)
	}
}
