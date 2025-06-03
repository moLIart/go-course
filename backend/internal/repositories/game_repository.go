package repositories

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"slices"
	"time"

	"github.com/moLIart/gomoku-backend/internal/domain"
	"github.com/moLIart/gomoku-backend/internal/infra"

	_ "github.com/lib/pq"
)

type GameRepository struct {
	db *infra.Database
}

func NewGameRepository(db *infra.Database) *GameRepository {
	return &GameRepository{
		db: db,
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
	conn, err := r.db.AcquireConn()
	if err != nil {
		return nil, err
	}

	var row gameWithPlayersRow
	err = conn.
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
