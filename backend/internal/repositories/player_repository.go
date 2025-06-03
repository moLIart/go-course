package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/moLIart/gomoku-backend/internal/domain"
	"github.com/moLIart/gomoku-backend/pkg/errorx"
)

type PlayerRepository struct {
	tx *sqlx.Tx
}

var (
	sqlInsertPlayer = `
		INSERT INTO players (nickname, password, score) 
		VALUES ($1, $2, $3) 
		RETURNING player_id`

	sqlGetPlayerByNickname = `
		SELECT player_id, nickname, password, score 
		FROM players 
		WHERE nickname = $1
		LIMIT 1`
)

func NewPlayerRepository(tx *sqlx.Tx) *PlayerRepository {
	return &PlayerRepository{
		tx: tx,
	}
}

// Insert inserts a new player into the database.
// It returns an error if the player already has an existing ID,
// if acquiring a database connection fails, or if the SQL insert operation fails.
// On successful insertion, the player's ID is set to the newly generated value.
func (r *PlayerRepository) Insert(player *domain.Player, ctx context.Context) error {
	if player.ID != 0 {
		return errors.New("cannot insert player with existing ID")
	}

	scanner := r.tx.QueryRowxContext(ctx, sqlInsertPlayer, player.Nickname, player.Password, player.Score)
	if err := scanner.Scan(&player.ID); err != nil {
		// Check if the error is a PostgreSQL unique violation error (duplicate key)
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			return domain.ErrPlayerAlreadyExists
		}

		return errorx.Wrap(err, "insert player sql")
	}

	return nil
}

// GetByNickname retrieves a player from the database by their nickname.
// It returns a pointer to the Player domain object if found, or nil if no player exists with the given nickname.
// If an error occurs during the database operation, it wraps and returns the error.
// The context parameter is used to control the lifetime of the database query.
func (r *PlayerRepository) GetByNickname(nickname string, ctx context.Context) (*domain.Player, error) {
	player := &domain.Player{}

	scanner := r.tx.QueryRowxContext(ctx, sqlGetPlayerByNickname, nickname)
	if err := scanner.Scan(&player.ID, &player.Nickname, &player.Password, &player.Score); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrPlayerNotFound // Player not found
		}

		return nil, errorx.Wrap(err, "get player by nickname sql")
	}

	return player, nil
}
