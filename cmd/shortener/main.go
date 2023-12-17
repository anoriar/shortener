package main

import (
	"fmt"
	"github.com/anoriar/shortener/internal/shortener/config"
	"github.com/anoriar/shortener/internal/shortener/logger"
	"github.com/anoriar/shortener/internal/shortener/repository/url"
	"github.com/anoriar/shortener/internal/shortener/router"
	"github.com/caarlos0/env/v6"
	"go.uber.org/zap"
	"net/http"
	_ "net/http/pprof"
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

	runProfiler(conf)

	logger, err := logger.Initialize(conf.LogLevel)
	if err != nil {
		panic(err)
	}

	defer logger.Sync()

	urlRepository, err := url.InitializeURLRepository(conf, logger)
	if err != nil {
		panic(err)
	}
	defer urlRepository.Close()

	r, err := router.InitializeRouter(conf, urlRepository, logger)

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

func runProfiler(cnf *config.Config) {
	if cnf.ProfilerHost != "" {
		go func() {
			fmt.Println("Starting pprof server at " + cnf.Host)
			err := http.ListenAndServe(cnf.ProfilerHost, nil)
			if err != nil {
				panic(err)
			}
		}()
	}
}
