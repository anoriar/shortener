package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"

	"github.com/caarlos0/env/v6"
	"go.uber.org/zap"

	"github.com/anoriar/shortener/internal/shortener/config"
	"github.com/anoriar/shortener/internal/shortener/logger"
	"github.com/anoriar/shortener/internal/shortener/repository/url"
	"github.com/anoriar/shortener/internal/shortener/router"
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

	runProfiler(conf, logger)

	urlRepository, err := url.InitializeURLRepository(conf, logger)
	if err != nil {
		panic(err)
	}
	defer urlRepository.Close()

	r, err := router.InitializeRouter(conf, urlRepository, logger)

	if err != nil {
		logger.Fatal("init error", zap.String("error", err.Error()))
		os.Exit(1)
	}

	err = http.ListenAndServe(conf.Host, r.Route())
	if err != nil {
		logger.Fatal("Server exception", zap.String("exception", err.Error()))
		os.Exit(1)
	}
}

func runProfiler(cnf *config.Config, logger *zap.Logger) {
	if cnf.ProfilerHost != "" {
		go func() {
			fmt.Println("Starting pprof server at " + cnf.Host)
			err := http.ListenAndServe(cnf.ProfilerHost, nil)
			if err != nil {
				logger.Fatal("internal server error", zap.String("error", err.Error()))
				os.Exit(1)
			}
		}()
	}
}
