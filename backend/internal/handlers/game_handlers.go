package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/moLIart/gomoku-backend/internal/domain"
	"github.com/moLIart/gomoku-backend/internal/middleware"
	"github.com/moLIart/gomoku-backend/internal/repositories"
)

type startGameRq struct {
	Type string `json:"game_type"`
	Size int    `json:"board_size"`
}

type moveGameRq struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

func HandleStartGame(uow *repositories.UnitOfWork) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var rq startGameRq
		if err := json.NewDecoder(r.Body).Decode(&rq); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if err := uow.Begin(r.Context()); err != nil {
			writeErrorRs(w, http.StatusInternalServerError, err)
			return
		}

		playerName := r.Context().Value(middleware.AuthPlayerNameKey).(string)

		players := uow.GetPlayerRepository()
		player, err := players.GetByNickname(playerName, r.Context())
		if err != nil {
			err = uow.Complete(err)
			writeErrorRs(w, http.StatusInternalServerError, err)
			return
		}

		board, err := domain.NewBoard(rq.Size)
		if err != nil {
			err = uow.Complete(err)
			writeErrorRs(w, http.StatusBadRequest, err)
			return
		}

		game, err := domain.NewGame(domain.GameType(rq.Type), board, player)
		if err != nil {
			err = uow.Complete(err)
			writeErrorRs(w, http.StatusBadRequest, err)
			return
		}

		games := uow.GetGameRepository()
		if err := games.Save(game, r.Context()); err != nil {
			err = uow.Complete(err)
			writeErrorRs(w, http.StatusInternalServerError, err)
			return
		}

		if err := uow.Complete(nil); err != nil {
			writeErrorRs(w, http.StatusInternalServerError, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(mapToGameState(game)); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	})
}

func HandleGetGameState(uow *repositories.UnitOfWork) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)

		gameId, err := strconv.Atoi(params.ByName("gameId"))
		if err != nil {
			http.Error(w, "Invalid game ID", http.StatusNotFound)
			return
		}

		if err := uow.Begin(r.Context()); err != nil {
			writeErrorRs(w, http.StatusInternalServerError, err)
			return
		}

		games := uow.GetGameRepository()
		game, err := games.GetById(int32(gameId), r.Context())
		if err != nil {
			if err == domain.ErrGameNotFound {
				uow.Complete(nil)

				http.Error(w, "Game not found", http.StatusNotFound)
				return
			}

			err = uow.Complete(err)
			writeErrorRs(w, http.StatusInternalServerError, err)
			return
		}

		if err := uow.Complete(nil); err != nil {
			writeErrorRs(w, http.StatusInternalServerError, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(mapToGameState(game)); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	})
}

func HandleGameJoin(uow *repositories.UnitOfWork) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
		playerName := r.Context().Value(middleware.AuthPlayerNameKey).(string)

		gameId, err := strconv.Atoi(params.ByName("gameId"))
		if err != nil {
			http.Error(w, "Invalid game ID", http.StatusNotFound)
			return
		}

		if err := uow.Begin(r.Context()); err != nil {
			writeErrorRs(w, http.StatusInternalServerError, err)
			return
		}

		player, err := uow.GetPlayerRepository().
			GetByNickname(playerName, r.Context())
		if err != nil {
			err = uow.Complete(err)
			writeErrorRs(w, http.StatusInternalServerError, err)
			return
		}

		games := uow.GetGameRepository()
		game, err := games.GetById(int32(gameId), r.Context())
		if err != nil {
			if err == domain.ErrGameNotFound {
				uow.Complete(nil)

				http.Error(w, "Game not found", http.StatusNotFound)
				return
			}

			err = uow.Complete(err)
			writeErrorRs(w, http.StatusInternalServerError, err)
			return
		}

		if err := game.Join(player); err != nil {
			err = uow.Complete(err)
			writeErrorRs(w, http.StatusBadRequest, err)
			return
		}

		if err := games.Save(game, r.Context()); err != nil {
			err = uow.Complete(err)
			writeErrorRs(w, http.StatusInternalServerError, err)
			return
		}

		if err := uow.Complete(nil); err != nil {
			writeErrorRs(w, http.StatusInternalServerError, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(mapToGameState(game)); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	})
}

func HandleGameMove(uow *repositories.UnitOfWork) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := r.Context().Value(httprouter.ParamsKey).(httprouter.Params)
		playerName := r.Context().Value(middleware.AuthPlayerNameKey).(string)

		gameId, err := strconv.Atoi(params.ByName("gameId"))
		if err != nil {
			http.Error(w, "Invalid game ID", http.StatusNotFound)
			return
		}

		if err := uow.Begin(r.Context()); err != nil {
			writeErrorRs(w, http.StatusInternalServerError, err)
			return
		}

		players := uow.GetPlayerRepository()

		player, err := players.GetByNickname(playerName, r.Context())
		if err != nil {
			err = uow.Complete(err)
			writeErrorRs(w, http.StatusInternalServerError, err)
			return
		}

		games := uow.GetGameRepository()
		game, err := games.GetById(int32(gameId), r.Context())
		if err != nil {
			if err == domain.ErrGameNotFound {
				uow.Complete(nil)

				http.Error(w, "Game not found", http.StatusNotFound)
				return
			}

			err = uow.Complete(err)
			writeErrorRs(w, http.StatusInternalServerError, err)
			return
		}

		var move moveGameRq
		if err := json.NewDecoder(r.Body).Decode(&move); err != nil {
			writeErrorRs(w, http.StatusBadRequest, err)
			return
		}

		if err := game.Move(move.Row, move.Col, player); err != nil {
			err = uow.Complete(err)
			writeErrorRs(w, http.StatusBadRequest, err)
			return
		}

		if ok, winner := game.HasWinner(); ok {
			winner.AddScore()

			if err := players.Save(winner, r.Context()); err != nil {
				err = uow.Complete(err)
				writeErrorRs(w, http.StatusInternalServerError, err)
				return
			}
		}

		if err := games.Save(game, r.Context()); err != nil {
			err = uow.Complete(err)
			writeErrorRs(w, http.StatusInternalServerError, err)
			return
		}

		if err := uow.Complete(nil); err != nil {
			writeErrorRs(w, http.StatusInternalServerError, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(mapToGameState(game)); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	})
}
