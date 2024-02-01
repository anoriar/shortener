package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/anoriar/shortener/internal/app"
	"github.com/anoriar/shortener/internal/shortener/config"
	"github.com/anoriar/shortener/internal/shortener/router"
	"github.com/anoriar/shortener/internal/shortener/server/tlscert"
)

// RunServer missing godoc.
func RunServer(app *app.App, r *router.Router) error {
	srv, err := createServer(app.Config, r)

	if err != nil {
		return err
	}
	gracefulShutdown(srv, app)

	return nil
}

func gracefulShutdown(srv *http.Server, app *app.App) {

	// Create a context that will be canceled on signal reception
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	go func() {
		<-ctx.Done()
		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown error: %v", err)
		}
	}()

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()

		if err := srv.ListenAndServe(); err != http.ErrServerClosed {
			log.Printf("Error starting the server: %v\n", err)
		}
	}()

	// background tasks
	wg.Add(1)
	go func() {
		defer wg.Done()
		app.DeleteUserURLsProcessor.Start(ctx)
	}()

	wg.Wait()

	fmt.Println("Graceful shutdown done")
}

func createServer(conf *config.Config, r *router.Router) (*http.Server, error) {
	var srv = &http.Server{Addr: conf.Host, Handler: r.Route()}
	if conf.EnableHTTPS {
		if !fileExists(tlscert.CertFilePath) || !fileExists(tlscert.PrivateKeyFilePath) {
			tlscert.GenerateTLSCert()
		}
		cert, err := tls.LoadX509KeyPair(tlscert.CertFilePath, tlscert.PrivateKeyFilePath)
		if err != nil {
			return nil, fmt.Errorf("error loading certificate and key: %v", err)
		}

		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
		srv.TLSConfig = tlsConfig
	}
	return srv, nil
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}
