package app

import (
	"github.com/djeebus/gpsctf/Backend/db"
	"sync"
)

type playerInfo struct {
	player    *db.Player
	latitude  float64
	longitude float64
}

type Worker struct {
	lock  sync.Mutex
	games map[int64]*GameProcessor
}

func NewWorker() *Worker {
	return &Worker{
		games: make(map[int64]*GameProcessor),
	}
}

func (w *Worker) Shutdown() {
	for gameId := range w.games {
		game := w.games[gameId]
		game.cancel <- true
	}
}

func (w *Worker) GetGameProcessor(gameId int64) (*GameProcessor, error) {
	processor, ok := w.games[gameId]
	if !ok {
		return nil, &GameDoesNotExistError{GameID: gameId}
	}

	return processor, nil
}

func (w *Worker) lockSelf() {
	w.lock.Lock()
}

func (w *Worker) unlockSelf() {
	w.lock.Unlock()
}

func (w *Worker) AddGame(game *db.Game) (*GameProcessor, error) {
	processor, ok := w.games[game.GameID]
	if ok {
		return processor, &GameInProgressError{GameID: game.GameID}
	}

	w.lockSelf()
	defer w.unlockSelf()

	gp := &GameProcessor{
		manager:    w,
		Game:       game,
		players:    make(map[int64]*playerInfo),
		cancel:     make(chan bool),
		register:   make(chan PlayerClient),
		unregister: make(chan PlayerClient),
		clients:    make(map[PlayerClient]bool),
	}
	w.games[game.GameID] = gp
	go gp.ProcessGame()

	return processor, nil
}

func (w *Worker) StartGame(gameId int64) (*GameProcessor, error) {
	processor, err := w.GetGameProcessor(gameId)
	if err != nil {
		return nil, err
	}

	err = processor.StartGame()
	if err != nil {
		return nil, err
	}

	return processor, nil
}
