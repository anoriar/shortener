package main

import (
	"github.com/anoriar/shortener/internal/config"
	"github.com/anoriar/shortener/internal/di"
	"github.com/anoriar/shortener/internal/router"
	"net/http"
)

func main() {
	run()
}

func run() {
	conf := config.NewConfig()
	parseFlags(conf)

	ctn, err := di.NewContainer(conf)
	if err != nil {
		panic(err)
	}

	r := ctn.Resolve("router").(*router.Router)

	err = http.ListenAndServe(conf.Host, r.Route())
	if err != nil {
		panic(err)
	}
}
