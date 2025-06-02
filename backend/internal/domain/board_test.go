package domain_test

import (
	"testing"

	"github.com/moLIart/gomoku-backend/internal/domain"
)

func TestNewBoard_ValidSize(t *testing.T) {
	board, err := domain.NewBoard(5)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if board.GetSize() != 5 {
		t.Errorf("expected board size 5, got %d", board.GetSize())
	}
}

func TestNewBoard_InvalidSize(t *testing.T) {
	board, err := domain.NewBoard(2)
	if err == nil {
		t.Fatal("expected error for board size < 3, got nil")
	}
	if board != nil {
		t.Errorf("expected nil board for invalid size, got %+v", board)
	}
}

func TestBoard_PutAndIsOccupied(t *testing.T) {
	board, err := domain.NewBoard(3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	player, _ := domain.NewPlayer("TestPlayer", "securepassword")
	player.ID = 1

	err = board.Put(1, 1, player)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !board.IsOccupied(1, 1) {
		t.Errorf("expected position (1,1) to be occupied")
	}
	if !board.IsOccupied(1, 1, player) {
		t.Errorf("expected position (1,1) to be occupied by player")
	}
}

func TestBoard_Put_OutOfBounds(t *testing.T) {
	board, err := domain.NewBoard(3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	player, _ := domain.NewPlayer("TestPlayer", "securepassword")
	player.ID = 1

	err = board.Put(3, 3, player)
	if err == nil {
		t.Fatal("expected error for out of bounds, got nil")
	}
}

func TestBoard_Put_AlreadyOccupied(t *testing.T) {
	board, err := domain.NewBoard(3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	player, _ := domain.NewPlayer("TestPlayer", "securepassword")
	player.ID = 1

	_ = board.Put(0, 0, player)
	err = board.Put(0, 0, player)
	if err == nil {
		t.Fatal("expected error for already occupied position, got nil")
	}
}

func TestBoard_IsOutOfBounds(t *testing.T) {
	board, err := domain.NewBoard(3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !board.IsOutOfBounds(-1, 0) || !board.IsOutOfBounds(0, -1) ||
		!board.IsOutOfBounds(3, 0) || !board.IsOutOfBounds(0, 3) {
		t.Errorf("expected out of bounds for invalid indices")
	}
	if board.IsOutOfBounds(1, 1) {
		t.Errorf("expected (1,1) to be in bounds")
	}
}

func TestBoard_CheckWin_Horizontal(t *testing.T) {
	board, err := domain.NewBoard(5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	player, _ := domain.NewPlayer("Winner", "securepassword")
	player.ID = 1

	for col := 0; col < 3; col++ {
		_ = board.Put(0, col, player)
	}
	if !board.CheckWin(0, 2, player, 3) {
		t.Errorf("expected horizontal win")
	}
}

func TestBoard_CheckWin_Vertical(t *testing.T) {
	board, err := domain.NewBoard(5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	player, _ := domain.NewPlayer("Winner", "securepassword")
	player.ID = 1

	for row := 0; row < 3; row++ {
		_ = board.Put(row, 0, player)
	}
	if !board.CheckWin(2, 0, player, 3) {
		t.Errorf("expected vertical win")
	}
}

func TestBoard_CheckWin_Diagonal(t *testing.T) {
	board, err := domain.NewBoard(5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	player, _ := domain.NewPlayer("Winner", "securepassword")
	player.ID = 1

	for i := 0; i < 3; i++ {
		_ = board.Put(i, i, player)
	}
	if !board.CheckWin(2, 2, player, 3) {
		t.Errorf("expected diagonal win")
	}
}

func TestBoard_CheckWin_AntiDiagonal(t *testing.T) {
	board, err := domain.NewBoard(5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	player, _ := domain.NewPlayer("Winner", "securepassword")
	player.ID = 1

	for i := 0; i < 3; i++ {
		_ = board.Put(i, 2-i, player)
	}
	if !board.CheckWin(2, 0, player, 3) {
		t.Errorf("expected anti-diagonal win")
	}
}

func TestBoard_CheckWin_NoWin(t *testing.T) {
	board, err := domain.NewBoard(5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	player, _ := domain.NewPlayer("NoWin", "securepassword")
	player.ID = 1

	_ = board.Put(0, 0, player)
	_ = board.Put(0, 1, player)
	if board.CheckWin(0, 1, player, 3) {
		t.Errorf("did not expect win")
	}
}
