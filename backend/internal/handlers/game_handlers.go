package handlers

import (
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
)

type startGameRq struct {
	Type string `json:"game_type"`
	Size int    `json:"board_size"`
}

func HandleStartGame() handlerClosureAlias {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	}
}

func HandleJoinGame() handlerClosureAlias {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	}
}

func HandleGameMove() handlerClosureAlias {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	}
}

// HandleGetGameState возвращает текущее состояние игры.
//
// @Summary      Get game state
// @Description  Returns the current state of the game (SSE endpoint)
// @Tags         game
// @Produce      json
// @Param        gameId   path      int     true  "Game ID"
// @Success      200      {object}  map[string]interface{}
// @Failure      400      {object}  errorRs
// @Failure      404      {object}  errorRs
// @Failure      500      {object}  errorRs
// @Router       /api/v1/games/{gameId} [get]
func HandleGetGameState() handlerClosureAlias {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Header().Add("Cache-Control", "no-cache, must-revalidate")
		w.Header().Add("Connection", "keep-alive")

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
			return
		}

		ticker := time.NewTicker(time.Second / 3)
		defer ticker.Stop()

		for {
			select {
			case <-r.Context().Done():
				// Client closed connection
				return
			case <-ticker.C:
				w.Write([]byte(`{"state": "new"}`)) // Replace with actual state
				flusher.Flush()
				continue
			}
		}

	}
}
