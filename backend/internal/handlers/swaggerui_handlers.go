package handlers

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"

	_ "github.com/moLIart/gomoku-backend/docs"
)

func SwaggerUIHandler() http.Handler {
	return httpSwagger.WrapHandler
}
