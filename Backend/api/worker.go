package api

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

func (worker *Worker) Shutdown() {
	for gameId := range worker.games {
		game := worker.games[gameId]
		game.cancel <- true
	}
}

func (worker *Worker) GetGameProcessor(gameId int64) (*GameProcessor, error) {
	processor, ok := worker.games[gameId]
	if !ok {
		return nil, &GameDoesNotExistError{GameID: gameId}
	}

	return processor, nil
}

func (worker *Worker) lockSelf() {
	worker.lock.Lock()
}

func (worker *Worker) unlockSelf() {
	worker.lock.Unlock()
}

func (worker *Worker) AddGame(game *db.Game) (*GameProcessor, error) {
	processor, ok := worker.games[game.GameID]
	if ok {
		return processor, &GameInProgressError{GameID: game.GameID}
	}

	worker.lockSelf()
	defer worker.unlockSelf()

	gp := &GameProcessor{
		worker:      worker,
		game:        game,
		players:     make(map[int64]*playerInfo),
		cancel:      make(chan bool),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		clients:     make(map[*Client]bool),
	}
	worker.games[game.GameID] = gp
	go gp.ProcessGame()

	return processor, nil
}

func (worker *Worker) StartGame(gameId int64) (*GameProcessor, error) {
	processor, err := worker.GetGameProcessor(gameId)
	if err != nil {
		return nil, err
	}

	err = processor.StartGame()
	if err != nil {
		return nil, err
	}

	return processor, nil
}
