package db

import "fmt"

type GameNotFoundError struct {
	GameID int64
}

func (err *GameNotFoundError) Error() string {
	return fmt.Sprintf("Game %d not found", err.GameID)
}

type PlayerNotFoundError struct {
	PlayerID int64
}

func (err *PlayerNotFoundError) Error() string {
	return fmt.Sprintf("Player %d not found", err.PlayerID)
}