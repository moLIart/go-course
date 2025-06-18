// @title Gomoku API
// @version 1.0
// @description API for Gomoku game
// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	log "github.com/sirupsen/logrus"

	"github.com/moLIart/gomoku-backend/internal/handlers"
	"github.com/moLIart/gomoku-backend/internal/infra"
	"github.com/moLIart/gomoku-backend/internal/middleware"
	"github.com/moLIart/gomoku-backend/internal/repositories"
	"github.com/moLIart/gomoku-backend/internal/services"
)

const (
	appName = "gomoku-api"
)

var (
	fs           = flag.NewFlagSet(appName, flag.ExitOnError)
	jwtSecret    = fs.String("jwt-secret", "", "JWT secret key")
	serverAddr   = fs.String("server-addr", ":8080", "http endpoint")
	dbDataSource = fs.String("db-data-source", "", "postgres data source (e.g., postgres://user:password@localhost:5432/dbname?sslmode=disable)")
)

func main() {
	fs.Parse(os.Args[1:])

	// Setup services
	jwtSvc := services.NewJWTService(*jwtSecret)
	database := infra.NewDatabase(*dbDataSource)
	uow := repositories.NewUnitOfWork(database)

	// Setup routing
	stdMiddlewares := alice.New(middleware.ContentType("application/json"))
	authMiddlewares := stdMiddlewares.Append(middleware.JWTAuth(jwtSvc))

	router := httprouter.New()
	router.Handler("GET", "/swagger/*any", handlers.SwaggerUIHandler())

	router.Handler("POST", "/api/v1/register",
		stdMiddlewares.Then(handlers.HandleRegister(uow, jwtSvc)))
	router.Handler("POST", "/api/v1/login",
		stdMiddlewares.Then(handlers.HandleLogin(uow, jwtSvc)))

	router.Handler("POST", "/api/v1/games/",
		authMiddlewares.Then(handlers.HandleStartGame(uow)))
	router.Handler("GET", "/api/v1/games/:gameId",
		authMiddlewares.Then(handlers.HandleGetGameState(uow)))
	router.Handler("PUT", "/api/v1/games/:gameId/move",
		authMiddlewares.Then(handlers.HandleGameMove(uow)))
	router.Handler("PUT", "/api/v1/games/:gameId/join",
		authMiddlewares.Then(handlers.HandleGameJoin(uow)))

	// Configure HTTP server
	httpSrv := &http.Server{
		Addr:    *serverAddr,
		Handler: router,
	}

	// Set up HTTP Server graceful shutdown
	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		<-interruptChan

		log.Trace("Shutting down db connections...")
		database.Stop()

		log.Info("Shutting down HTTP Server...")
		if err := httpSrv.Shutdown(context.Background()); err != nil {
			log.Fatalf("HTTP Server forced to shutdown: %s", err)
		}
	}()

	// Start the database connection
	if err := database.Start(context.Background()); err != nil {
		log.Fatalf("Could not start database: %s", err)
	}

	// Starting the HTTP server
	log.Infof("Starting HTTP server on %s", httpSrv.Addr)
	if err := httpSrv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Could not start HTTP server: %s", err)
	}

	// Wait for the shutdown to complete
	wg.Wait()
	log.Info("Server gracefully stopped")
}
