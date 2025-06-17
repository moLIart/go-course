package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/moLIart/gomoku-backend/internal/domain"
	"github.com/moLIart/gomoku-backend/internal/repositories"
	"github.com/moLIart/gomoku-backend/internal/services"
)

type registerRq struct {
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

type registerRs struct {
	Token string `json:"token"`
}

type loginRq struct {
	Nickname string `json:"nickname"`
	Password string `json:"password"`
}

type loginRs struct {
	Token string `json:"token"`
}

// HandleRegister регистрирует нового пользователя.
//
// @Summary      Register new player
// @Description  Creates a new player and returns a JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        registerRq  body      registerRq  true  "Registration data"
// @Success      200         {object}  registerRs
// @Failure      400         {object}  errorRs
// @Failure      409         {object}  errorRs
// @Failure      500         {object}  errorRs
// @Router       /api/v1/register [post]
func HandleRegister(uow *repositories.UnitOfWork, jwtSvc *services.JWTService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var rq registerRq
		if err := json.NewDecoder(r.Body).Decode(&rq); err != nil {
			writeErrorRs(w, http.StatusBadRequest, err)
			return
		}

		player, err := domain.NewPlayer(rq.Nickname, rq.Password)
		if err != nil {
			writeErrorRs(w, http.StatusBadRequest, err)
			return
		}

		if err := uow.Begin(r.Context()); err != nil {
			writeErrorRs(w, http.StatusInternalServerError, err)
			return
		}

		repository := uow.GetPlayerRepository()
		if err := repository.Insert(player, r.Context()); err != nil {
			if errors.Is(err, domain.ErrPlayerAlreadyExists) {
				uow.Complete(nil)
				writeErrorRs(w, http.StatusConflict, err)
				return
			}

			uow.Complete(err)
			writeErrorRs(w, http.StatusInternalServerError, err)
			return
		}

		if err := uow.Complete(nil); err != nil {
			writeErrorRs(w, http.StatusInternalServerError, err)
			return
		}

		tokenString, err := jwtSvc.Sign(player.Nickname)
		if err != nil {
			writeErrorRs(w, http.StatusInternalServerError, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&registerRs{Token: tokenString}); err != nil {
			writeErrorRs(w, http.StatusInternalServerError, err)
			return
		}
	})
}

// HandleLogin аутентифицирует пользователя.
//
// @Summary      Login player
// @Description  Authenticates player and returns a JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        loginRq  body      loginRq  true  "Login data"
// @Success      200      {object}  loginRs
// @Failure      400      {object}  errorRs
// @Failure      401      {object}  errorRs
// @Failure      500      {object}  errorRs
// @Router       /api/v1/login [post]
func HandleLogin(uow *repositories.UnitOfWork, jwtSvc *services.JWTService) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var rq loginRq
		if err := json.NewDecoder(r.Body).Decode(&rq); err != nil {
			writeErrorRs(w, http.StatusBadRequest, err)
			return
		}

		if err := uow.Begin(r.Context()); err != nil {
			writeErrorRs(w, http.StatusInternalServerError, err)
			return
		}

		repository := uow.GetPlayerRepository()

		player, err := repository.GetByNickname(rq.Nickname, r.Context())
		if err != nil {
			if errors.Is(err, domain.ErrPlayerNotFound) {
				uow.Complete(nil)
				writeErrorRs(w, http.StatusUnauthorized, err)
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

		if player.Password != rq.Password {
			writeErrorRs(w, http.StatusUnauthorized, errors.New("invalid credentials"))
			return
		}

		tokenString, err := jwtSvc.Sign(player.Nickname)
		if err != nil {
			writeErrorRs(w, http.StatusInternalServerError, err)
			return
		}

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&loginRs{Token: tokenString}); err != nil {
			writeErrorRs(w, http.StatusInternalServerError, err)
			return
		}
	})
}
