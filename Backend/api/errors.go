package api

import "fmt"

type GameNotInProgressError struct {
	GameID int64
}

func (err *GameNotInProgressError) Error() string {
	return fmt.Sprintf("Game #%d not in progress", err.GameID)
}

