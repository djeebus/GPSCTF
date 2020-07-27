package app

import "fmt"

type GameInProgressError struct {
	GameID int64
}

func (err *GameInProgressError) Error() string {
	return fmt.Sprintf("Game #%d already in progress", err.GameID)
}

type GameHasNoPlayersError struct {
	GameID int64
}

func (err *GameHasNoPlayersError) Error() string {
	return fmt.Sprintf("Game #%d has no players", err.GameID)
}

type GameDoesNotExistError struct {
	GameID int64
}

func (err *GameDoesNotExistError) Error() string {
	return fmt.Sprintf("Game #%d does not exist", err.GameID)
}
