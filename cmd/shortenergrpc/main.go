package main

import (
	"fmt"
	"log"

	appPkg "github.com/anoriar/shortener/internal/app"
	//"github.com/anoriar/shortener/internal/shortener/config/file"
	"github.com/anoriar/shortener/internal/shortener/server"

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

	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("create config error %v", err.Error())
	}

	logger, err := logger.Initialize(conf.LogLevel)
	if err != nil {
		log.Fatalf("init logger error %v", err.Error())
	}

	defer logger.Sync()

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
