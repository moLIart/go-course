package handlers

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"

	"github.com/moLIart/gomoku-backend/pkg/middleware"
)

func RegisterHttpRoutes() http.Handler {
	router := httprouter.New()

	routerMiddlewares := alice.New(
		middleware.ContentType("application/json"),
	).Then(router)

	return routerMiddlewares
}
