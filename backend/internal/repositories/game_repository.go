package repositories

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"slices"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/moLIart/gomoku-backend/internal/domain"

	_ "github.com/lib/pq"
)

type GameRepository struct {
	tx *sqlx.Tx
}

func NewGameRepository(tx *sqlx.Tx) *GameRepository {
	return &GameRepository{
		tx: tx,
	}
}

var (
	sqlGetGameById = `
		SELECT 
			game_id, type, board, current_player_id, winner_player_id, first_player_id, second_player_id, last_activity,
			fp.player_id as fp_id, fp.nickname as fp_nickname, fp.password as fp_password, fp.score as fp_score,
			sp.player_id as sp_id, sp.nickname as sp_nickname, sp.password as sp_password, sp.score as sp_score 
		FROM games
			LEFT JOIN players AS fp ON fp.player_id = games.first_player_id
			LEFT JOIN players AS sp ON sp.player_id = games.second_player_id
		WHERE game_id = $1
		LIMIT 1`

	sqlInsertGame = `
		INSERT INTO games (type, board, current_player_id, winner_player_id, first_player_id, second_player_id, last_activity)
		VALUES ($1, $2::jsonb, $3, $4, $5, $6, $7)
		RETURNING game_id`

	sqlUpdateGame = `
		UPDATE games
		SET type = $1, board = $2::jsonb, current_player_id = $3, winner_player_id = $4, first_player_id = $5, second_player_id = $6, last_activity = $7
		WHERE game_id = $8`
)

// This struct matches the SELECT columns in sqlGetGameById
type gameWithPlayersRow struct {
	GameID          int32         `db:"game_id"`
	Type            string        `db:"type"`
	Board           boardDto      `db:"board"`
	CurrentPlayerID int32         `db:"current_player_id"`
	WinnerPlayerID  sql.NullInt32 `db:"winner_player_id"`
	FirstPlayerID   int32         `db:"first_player_id"`
	SecondPlayerID  sql.NullInt32 `db:"second_player_id"`
	LastActivity    time.Time     `db:"last_activity"`

	FPID       int32  `db:"fp_id"`
	FPNickname string `db:"fp_nickname"`
	FPPassword string `db:"fp_password"`
	FPScore    int32  `db:"fp_score"`

	SPID       sql.NullInt32  `db:"sp_id"`
	SPNickname sql.NullString `db:"sp_nickname"`
	SPPassword sql.NullString `db:"sp_password"`
	SPScore    sql.NullInt32  `db:"sp_score"`
}

type boardDto struct {
	Size int       `json:"size"`
	Data [][]int32 `json:"data"`
}

func DtoFromBoard(b *domain.Board) boardDto {
	dto := boardDto{
		Size: b.Size,
		Data: make([][]int32, b.Size),
	}

	for i := range b.Data {
		dto.Data[i] = make([]int32, b.Size)
		copy(dto.Data[i], b.Data[i])
	}

	return dto
}

func (b *boardDto) Value() (driver.Value, error) {
	j, err := json.Marshal(b)
	return j, err
}

func (p *boardDto) Scan(src any) error {
	source, ok := src.([]byte)
	if !ok {
		return errors.New("type assertion .([]byte) failed")
	}

	return json.Unmarshal(source, p)
}

func (r *GameRepository) GetById(id int32, ctx context.Context) (*domain.Game, error) {
	var row gameWithPlayersRow
	err := r.tx.
		QueryRowxContext(ctx, sqlGetGameById, id).
		StructScan(&row)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrGameNotFound
		}
		return nil, err
	}

	game := &domain.Game{Board: &domain.Board{}, Players: [2]*domain.Player{}}

	game.ID = row.GameID
	game.Type = domain.GameType(row.Type)
	game.LastActivity = row.LastActivity

	game.Board.Size = row.Board.Size
	game.Board.Data = row.Board.Data

	game.Players[0] = &domain.Player{
		Entity: domain.Entity{
			ID: row.FPID,
		},
		Nickname: row.FPNickname,
		Password: row.FPPassword,
		Score:    int(row.FPScore),
	}

	if row.SPID.Valid {
		game.Players[1] = &domain.Player{
			Entity: domain.Entity{
				ID: row.SPID.Int32,
			},
			Nickname: row.SPNickname.String,
			Password: row.SPPassword.String,
			Score:    int(row.SPScore.Int32),
		}
	} else {
		game.Players[1] = nil
	}

	curIdx := slices.IndexFunc(game.Players[:], func(p *domain.Player) bool {
		return p != nil && p.ID == row.CurrentPlayerID
	})
	game.CurrentPlayer = game.Players[curIdx]

	if row.WinnerPlayerID.Valid {
		winnerIdx := slices.IndexFunc(game.Players[:], func(p *domain.Player) bool {
			return p != nil && p.ID == row.WinnerPlayerID.Int32
		})

		game.WinnerPlayer = game.Players[winnerIdx]
	} else {
		game.WinnerPlayer = nil
	}

	return game, nil
}

func (r *GameRepository) Save(game *domain.Game, ctx context.Context) error {
	boardDto := DtoFromBoard(game.Board)
	winnerPlayerID := sql.NullInt32{}
	if game.WinnerPlayer != nil {
		winnerPlayerID = sql.NullInt32{Int32: game.WinnerPlayer.ID, Valid: true}
	}

	secondPlayerID := sql.NullInt32{}
	if game.Players[1] != nil {
		secondPlayerID = sql.NullInt32{Int32: game.Players[1].ID, Valid: true}
	}

	boardJson, err := json.Marshal(boardDto)
	if err != nil {
		return err
	}

	if game.ID != 0 {
		// Update existing game
		_, err := r.tx.ExecContext(ctx, sqlUpdateGame,
			game.Type,
			boardJson,
			game.CurrentPlayer.ID,
			winnerPlayerID,
			game.Players[0].ID,
			secondPlayerID,
			game.LastActivity,
			game.ID)
		if err != nil {
			return err
		}
	} else {
		// Insert new game
		err := r.tx.QueryRowContext(ctx, sqlInsertGame,
			game.Type,
			boardJson,
			game.CurrentPlayer.ID,
			winnerPlayerID,
			game.Players[0].ID,
			secondPlayerID,
			game.LastActivity).Scan(&game.ID)
		if err != nil {
			return err
		}
	}

	return nil
}
