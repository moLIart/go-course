package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/moLIart/gomoku-backend/internal/domain"
	log "github.com/sirupsen/logrus"
	"gopkg.in/guregu/null.v3"
)

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

type gameStateDto struct {
	ID            int          `json:"id"`
	Type          string       `json:"type"`
	CurrentPlayer int          `json:"current_player"`
	Winner        null.Int     `json:"winner,omitempty"`
	Size          int          `json:"size"`
	Board         [][]null.Int `json:"board"`
}

func mapToGameState(game *domain.Game) *gameStateDto {
	dto := &gameStateDto{
		ID:            int(game.ID),
		Type:          string(game.Type),
		Size:          game.Board.Size,
		CurrentPlayer: int(game.CurrentPlayer.ID),
	}

	if game.WinnerPlayer != nil {
		dto.Winner = null.IntFrom(int64(game.WinnerPlayer.ID))
	}

	dto.Size = game.Board.Size
	dto.Board = make([][]null.Int, dto.Size)
	for i := 0; i < dto.Size; i++ {
		dto.Board[i] = make([]null.Int, dto.Size)
		for j := 0; j < dto.Size; j++ {
			if game.Board.Data[i][j] != 0 {
				dto.Board[i][j] = null.IntFrom(int64(game.Board.Data[i][j]))
			} else {
				dto.Board[i][j] = null.Int{}
			}
		}
	}

	return dto
}
