package main

import (
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"

	"github.com/anoriar/shortener/internal/shortener/util/tls"

	"github.com/caarlos0/env/v6"
	"go.uber.org/zap"

	"github.com/anoriar/shortener/internal/shortener/config"
	"github.com/anoriar/shortener/internal/shortener/logger"
	"github.com/anoriar/shortener/internal/shortener/repository/url"
	"github.com/anoriar/shortener/internal/shortener/router"
)

var buildVersion string = "N/A"
var buildDate string = "N/A"
var buildCommit string = "N/A"

func main() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

	conf := config.NewConfig()
	parseFlags(conf)

	err := env.Parse(conf)
	if err != nil {
		log.Fatalf("parse config error %v", err.Error())
	}

	logger, err := logger.Initialize(conf.LogLevel)
	if err != nil {
		log.Fatalf("init logger error %v", err.Error())
	}

	defer logger.Sync()

	runProfiler(conf, logger)

	urlRepository, err := url.InitializeURLRepository(conf, logger)
	if err != nil {
		log.Fatalf("init repository error %v", err.Error())
	}
	defer urlRepository.Close()

	r, err := router.InitializeRouter(conf, urlRepository, logger)

	if err != nil {
		log.Fatalf("init router error %v", err.Error())
	}

	runServer(conf, r)

}

func runServer(conf *config.Config, r *router.Router) {
	var err error
	if conf.EnableHttps {
		tls.GenerateTLSCert()
		err = http.ListenAndServeTLS(conf.Host, tls.CertFilePath, tls.PrivateKeyFilePath, r.Route())
	} else {
		err = http.ListenAndServe(conf.Host, r.Route())
	}

	if err != nil {
		log.Fatalf("server error %v", err.Error())
	}
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
