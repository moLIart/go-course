package domain

import (
	"errors"
	"time"
)

type GameType string

const (
	PvP GameType = "pvp"
	PvA GameType = "pva"
)

var (
	ErrInvalidGameType    = errors.New("invalid game type")
	ErrFullGame           = errors.New("game is full")
	ErrCantJoinToSameGame = errors.New("can't join to the same game")
	ErrInvalidBoard       = errors.New("invalid board")
	ErrNotYourTurn        = errors.New("it's not your turn")
	ErrGameNotReady       = errors.New("game is not ready")
	ErrGameNotFound       = errors.New("game is not found")
)

type Game struct {
	Entity

	Type  GameType
	Board *Board

	CurrentPlayer *Player
	WinnerPlayer  *Player
	Players       [2]*Player

	LastActivity time.Time
}

func NewGame(gtype GameType, board *Board, firstPlayer *Player) (*Game, error) {
	if gtype != PvP && gtype != PvA {
		return nil, ErrInvalidGameType
	}

	if board == nil {
		return nil, ErrInvalidBoard
	}

	game := &Game{
		Type:          gtype,
		Board:         board,
		CurrentPlayer: firstPlayer,
		WinnerPlayer:  nil,
		Players:       [2]*Player{firstPlayer, nil},
		LastActivity:  time.Now(),
	}

	return game, nil
}

func (g *Game) Join(player *Player) error {
	if g.Players[1] != nil {
		return ErrFullGame
	}

	if g.Players[0].Equal(player) {
		return ErrCantJoinToSameGame
	}

	g.Players[1] = player
	g.LastActivity = time.Now()
	return nil
}

func (g *Game) IsReady() bool {
	return g.Players[1] != nil
}

func (g *Game) Move(row, col int, player *Player) error {
	if !g.IsReady() {
		return ErrGameNotReady
	}

	if g.CurrentPlayer == nil || !g.CurrentPlayer.Equal(player) {
		return ErrNotYourTurn
	}

	if err := g.Board.Put(row, col, player); err != nil {
		return err
	}

	if g.Board.CheckWin(row, col, player, 3) { // Assuming 3 in a row to win
		g.WinnerPlayer = player
	} else {
		if g.Players[0].Equal(player) {
			// Switch to the second player
			g.CurrentPlayer = g.Players[1]
		} else if g.Players[1].Equal(player) {
			// Switch to the first player
			g.CurrentPlayer = g.Players[0]
		}
	}

	g.LastActivity = time.Now()
	return nil
}

func (g *Game) HasWinner() (bool, *Player) {
	if g.WinnerPlayer != nil {
		return true, g.WinnerPlayer
	}
	return false, nil
}
