package domain

import (
	"errors"
	"fmt"
)

type Player struct {
	Entity

	Nickname string
	Password string
	Score    int
}

var (
	ErrPlayerNotFound      = errors.New("player not found")
	ErrPlayerAlreadyExists = errors.New("player with same nickname already exists")
)

func NewPlayer(nickname string, password string) (*Player, error) {
	if len(nickname) < 3 {
		return nil, fmt.Errorf("nickname must be at least 3 characters long")
	}
	if len(nickname) > 20 {
		return nil, fmt.Errorf("nickname must be at most 20 characters long")
	}

	if len(password) < 6 || len(password) > 20 {
		return nil, fmt.Errorf("password must be between 6 and 20 characters long")
	}

	p := &Player{
		Nickname: nickname,
		Password: password,
		Score:    0,
	}

	return p, nil
}

func (p *Player) AddScore() {
	p.Score += 1
}

func (p *Player) DecScore() {
	p.Score -= 1
}

func (p *Player) Equal(other *Player) bool {
	if p == nil || other == nil {
		return false
	}
	return p.ID == other.ID
}
