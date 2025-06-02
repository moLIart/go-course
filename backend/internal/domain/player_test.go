package domain_test

import (
	"testing"

	"github.com/moLIart/gomoku-backend/internal/domain"
)

func TestNewPlayer_ValidNickname(t *testing.T) {
	player, err := domain.NewPlayer("Alice", "securepassword")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if player.Nickname != "Alice" {
		t.Errorf("expected nickname 'Alice', got %s", player.Nickname)
	}
	if player.Score != 0 {
		t.Errorf("expected initial score 0, got %d", player.Score)
	}
}

func TestNewPlayer_ShortNickname(t *testing.T) {
	_, err := domain.NewPlayer("Al", "securepassword")
	if err == nil {
		t.Fatal("expected error for short nickname, got nil")
	}
}

func TestNewPlayer_LongNickname(t *testing.T) {
	_, err := domain.NewPlayer("ThisNicknameIsWayTooLongForTheGame", "securepassword")
	if err == nil {
		t.Fatal("expected error for long nickname, got nil")
	}
}

func TestPlayer_AddScore(t *testing.T) {
	player, _ := domain.NewPlayer("Bob", "securepassword")
	player.AddScore()
	if player.Score != 1 {
		t.Errorf("expected score 1 after AddScore, got %d", player.Score)
	}
}

func TestPlayer_DecScore(t *testing.T) {
	player, _ := domain.NewPlayer("Charlie", "securepassword")
	player.AddScore()
	player.DecScore()
	if player.Score != 0 {
		t.Errorf("expected score 0 after AddScore and DecScore, got %d", player.Score)
	}
}
func TestPlayer_Equal_SameID(t *testing.T) {
	player1, _ := domain.NewPlayer("Dave", "securepassword")
	player2, _ := domain.NewPlayer("Eve", "securepassword")
	player1.ID = 42
	player2.ID = 42

	if !player1.Equal(player2) {
		t.Errorf("expected players with same ID to be equal")
	}
}

func TestPlayer_Equal_DifferentID(t *testing.T) {
	player1, _ := domain.NewPlayer("Frank", "securepassword")
	player2, _ := domain.NewPlayer("Grace", "securepassword")
	player1.ID = 1
	player2.ID = 2

	if player1.Equal(player2) {
		t.Errorf("expected players with different IDs to not be equal")
	}
}

func TestPlayer_Equal_NilReceiver(t *testing.T) {
	var player1 *domain.Player
	player2, _ := domain.NewPlayer("Heidi", "securepassword")
	player2.ID = 1

	if player1.Equal(player2) {
		t.Errorf("expected nil receiver to not be equal to any player")
	}
}

func TestPlayer_Equal_NilArgument(t *testing.T) {
	player1, _ := domain.NewPlayer("Ivan", "securepassword")
	player1.ID = 1
	var player2 *domain.Player

	if player1.Equal(player2) {
		t.Errorf("expected player not to be equal to nil argument")
	}
}

func TestPlayer_Equal_BothNil(t *testing.T) {
	var player1 *domain.Player
	var player2 *domain.Player

	if player1.Equal(player2) {
		t.Errorf("expected nil players to not be equal")
	}
}
