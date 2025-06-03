package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/julienschmidt/httprouter"
)

type handlerClosureAlias = func(w http.ResponseWriter, r *http.Request, _ httprouter.Params)

type errorRs struct {
	Error string `json:"error"`
}

func writeErrorRs(w http.ResponseWriter, code int, err error) {
	log.Errorf("handler error: %v", err)

	w.WriteHeader(code)

	errorRs := errorRs{Error: strings.TrimSpace(err.Error())}
	if err = json.NewEncoder(w).Encode(errorRs); err != nil {
		log.Errorf("error writing response: %v", err)
	}
}
