package domain

import (
	"testing"
	"time"
)

type mockBoard struct {
	Board
}

type mockPlayer struct {
	Player
}

func TestNewGame_SuccessPvP(t *testing.T) {
	board := &mockBoard{}
	player := &mockPlayer{}
	game, err := NewGame(PvP, &board.Board, &player.Player)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if game == nil {
		t.Fatal("expected game, got nil")
	}
	if game.Type != PvP {
		t.Errorf("expected type %v, got %v", PvP, game.Type)
	}
	if game.Board != &board.Board {
		t.Errorf("expected board to be set")
	}
	if game.CurrentPlayer != &player.Player {
		t.Errorf("expected current player to be set")
	}
	if game.Players[0] != &player.Player {
		t.Errorf("expected first player to be set")
	}
	if game.Players[1] != nil {
		t.Errorf("expected second player to be nil")
	}
	if time.Since(game.LastActivity) > time.Second {
		t.Errorf("expected recent last activity")
	}
}

func TestNewGame_SuccessPvA(t *testing.T) {
	board := &mockBoard{}
	player := &mockPlayer{}
	game, err := NewGame(PvA, &board.Board, &player.Player)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if game.Type != PvA {
		t.Errorf("expected type %v, got %v", PvA, game.Type)
	}
}

func TestNewGame_InvalidGameType(t *testing.T) {
	board := &mockBoard{}
	player := &mockPlayer{}
	game, err := NewGame("invalid", &board.Board, &player.Player)
	if err != ErrInvalidGameType {
		t.Errorf("expected ErrInvalidGameType, got %v", err)
	}
	if game != nil {
		t.Errorf("expected game to be nil")
	}
}

func TestNewGame_NilBoard(t *testing.T) {
	player := &mockPlayer{}
	game, err := NewGame(PvP, nil, &player.Player)
	if err != ErrInvalidBoard {
		t.Errorf("expected ErrInvalidBoard, got %v", err)
	}
	if game != nil {
		t.Errorf("expected game to be nil")
	}
}
func TestGame_Join_Success(t *testing.T) {
	board := &mockBoard{}
	player1 := &mockPlayer{Player: Player{Entity: Entity{ID: 1}}}
	player2 := &mockPlayer{Player: Player{Entity: Entity{ID: 2}}}
	game, err := NewGame(PvP, &board.Board, &player1.Player)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = game.Join(&player2.Player)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if game.Players[1] != &player2.Player {
		t.Errorf("expected second player to be set")
	}
}

func TestGame_Join_FullGame(t *testing.T) {
	board := &mockBoard{}
	player1 := &mockPlayer{Player: Player{Entity: Entity{ID: 1}}}
	player2 := &mockPlayer{Player: Player{Entity: Entity{ID: 2}}}
	player3 := &mockPlayer{Player: Player{Entity: Entity{ID: 3}}}
	game, err := NewGame(PvP, &board.Board, &player1.Player)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_ = game.Join(&player2.Player)

	err = game.Join(&player3.Player)
	if err != ErrFullGame {
		t.Errorf("expected ErrFullGame, got %v", err)
	}
}

func TestGame_Join_SamePlayer(t *testing.T) {
	board := &mockBoard{}
	player1 := &mockPlayer{Player: Player{Entity: Entity{ID: 1}}}
	game, err := NewGame(PvP, &board.Board, &player1.Player)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = game.Join(&player1.Player)
	if err != ErrCantJoinToSameGame {
		t.Errorf("expected ErrCantJoinToSameGame, got %v", err)
	}
}

func TestGame_IsReady_FalseWhenSecondPlayerNil(t *testing.T) {
	board := &mockBoard{}
	player1 := &mockPlayer{Player: Player{Entity: Entity{ID: 1}}}
	game, err := NewGame(PvP, &board.Board, &player1.Player)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if game.IsReady() {
		t.Errorf("expected IsReady to be false when second player is nil")
	}
}

func TestGame_IsReady_TrueWhenSecondPlayerSet(t *testing.T) {
	board := &mockBoard{}
	player1 := &mockPlayer{Player: Player{Entity: Entity{ID: 1}}}
	player2 := &mockPlayer{Player: Player{Entity: Entity{ID: 2}}}
	game, err := NewGame(PvP, &board.Board, &player1.Player)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	err = game.Join(&player2.Player)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !game.IsReady() {
		t.Errorf("expected IsReady to be true when second player is set")
	}
}
