package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"

	appPkg "github.com/anoriar/shortener/internal/app"
	"github.com/anoriar/shortener/internal/shortener/config/file"
	"github.com/anoriar/shortener/internal/shortener/server"

	"github.com/caarlos0/env/v6"
	"go.uber.org/zap"

	"github.com/anoriar/shortener/internal/shortener/config"
	"github.com/anoriar/shortener/internal/shortener/logger"
	"github.com/anoriar/shortener/internal/shortener/router"
)

var buildVersion string = "N/A"
var buildDate string = "N/A"
var buildCommit string = "N/A"

func main() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

	conf, err := createConfig()
	if err != nil {
		log.Fatalf("create config error %v", err.Error())
	}

	logger, err := logger.Initialize(conf.LogLevel)
	if err != nil {
		log.Fatalf("init logger error %v", err.Error())
	}

	defer logger.Sync()

	runProfiler(conf, logger)

	app, err := appPkg.NewApp(conf, logger)
	if err != nil {
		log.Fatalf("init app error %v", err.Error())
	}

	r := router.InitializeRouter(app)

	err = server.RunServer(app, r)
	if err != nil {
		log.Fatalf("init router error %v", err.Error())
	}

}

func createConfig() (*config.Config, error) {
	conf := config.NewConfig()
	parseFlags(conf)

	err := env.Parse(conf)
	if err != nil {
		return nil, fmt.Errorf("parse env error: %v", err)
	}

	err = file.LoadAndMergeConfig(conf)
	if err != nil {
		return nil, fmt.Errorf("parse config from file error: %v", err)
	}

	return conf, nil
}

func runProfiler(cnf *config.Config, logger *zap.Logger) {
	if cnf.ProfilerHost != "" {
		go func() {
			fmt.Println("Starting pprof server at " + cnf.Host)
			err := http.ListenAndServe(cnf.ProfilerHost, nil)
			if err != nil {
				log.Fatalf("profiler server error %v", err.Error())
			}
		}()
	}
}
