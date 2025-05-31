package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/moLIart/gomoku-backend/internal/handlers"
	"github.com/moLIart/gomoku-backend/pkg/errorx"
)

const (
	appName = "gomoku-api"
)

var (
	fs   = flag.NewFlagSet(appName, flag.ExitOnError)
	addr = fs.String("addr", ":8080", "http endpoint")

// dbConnStr = fs.String("db-connection-string", "", "db connection string")
)

func main() {
	errorx.MustNoError(fs.Parse(os.Args[1:]), "parse args")

	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)

	// Configure HTTP server
	httpSrv := &http.Server{
		Addr:    *addr,
		Handler: handlers.RegisterHttpRoutes(),
	}

	// Shutting down the server gracefully
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		<-interruptChan

		log.Println("Shutting down HTTP Server...")
		if err := httpSrv.Shutdown(context.Background()); err != nil {
			log.Fatalf("HTTP Server forced to shutdown: %s", err)
		}
	}()

	// Starting the HTTP server
	log.Println("Starting HTTP server on :8080")
	if err := httpSrv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Could not start HTTP server: %s", err)
	}

	// Wait for the shutdown to complete
	wg.Wait()
	log.Println("Server gracefully stopped")
}
