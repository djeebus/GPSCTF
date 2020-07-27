package app

import (
	"encoding/json"
	"github.com/djeebus/gpsctf/Backend/db"
	"log"
	"sync"
	"time"
)

type PlayerClient interface {
	Close()
	Send(buffer []byte)
}

type GameProcessor struct {
	// public
	Game        *db.Game

	// private
	manager *Worker
	lock    sync.Mutex

	inProgress  bool
	start       time.Time
	finish      time.Time
	playerCount int
	players     map[int64]*playerInfo

	flagLatitude  float64
	flagLongitude float64

	cancel chan bool

	register   chan PlayerClient
	unregister chan PlayerClient
	clients    map[PlayerClient]bool
}

func (gp *GameProcessor) Register(client PlayerClient) {
	gp.register <- client
}

func (gp *GameProcessor) Unregister(client PlayerClient) {
	gp.unregister <- client
}

func (gp *GameProcessor) AddPlayer(player *db.Player) {
	gp.lock.Lock()
	defer gp.lock.Unlock()

	gp.playerCount += 1
	gp.players[player.PlayerID] = &playerInfo{
		player:    player,
		latitude:  0,
		longitude: 0,
	}
}

func (gp *GameProcessor) UpdatePlayer(playerId int64, latitude float64, longitude float64) error {
	player, ok := gp.players[playerId]
	if !ok {
		return &db.PlayerNotFoundError{PlayerID: playerId}
	}

	player.latitude = latitude
	player.longitude = longitude
	return nil
}

func (gp *GameProcessor) StartGame() error {
	game := gp.Game

	if gp.inProgress {
		return &GameInProgressError{GameID: game.GameID}
	}

	if gp.playerCount == 0 {
		return &GameHasNoPlayersError{GameID: game.GameID}
	}

	flagLng, flagLat := getRandomLocation(game.Latitude, game.Longitude, game.Radius)

	gp.start = time.Now()
	gp.finish = time.Now().Add(60 * time.Second)
	gp.flagLongitude = flagLng
	gp.flagLatitude = flagLat
	gp.inProgress = true

	return nil
}

func (gp *GameProcessor) ProcessGame() {
	ticker := time.NewTicker(1 * time.Second)

	for {
		select {

		// deal with clients
		case client := <-gp.register:
			gp.clients[client] = true
		case client := <-gp.unregister:
			if _, ok := gp.clients[client]; ok {
				delete(gp.clients, client)
				client.Close()
			}

		// app was canceled
		case <-gp.cancel:
			gp.endGame()
			return

		// run the app
		case <-ticker.C:
			if gp.inProgress {
				gameOver := gp.gameTick()
				if gameOver {
					gp.endGame()
				}
			}
		}
	}
}

func (gp *GameProcessor) endGame() {
	delete(gp.manager.games, gp.Game.GameID)

	for client := range gp.clients {
		delete(gp.clients, client)
		client.Close()
	}
}

const minimumFeetToWin = 50

func (gp *GameProcessor) gameTick() bool {
	var winner *playerInfo
	for playerId := range gp.players {
		player := gp.players[playerId]
		distance := computeDistance(
			player.latitude, player.longitude,
			gp.flagLatitude, gp.flagLongitude)
		if distance <= minimumFeetToWin {
			winner = player
			break
		}
	}

	if winner == nil {
		return false // keep playing
	}

	message := map[string]interface{}{
		"winner": "someone won!",
		"playerId": winner.player.PlayerID,
	}
	gp.broadcast(message)
	return true
}

func (gp *GameProcessor) broadcast(message interface{}) {
	buf, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to encode json message: %v", err)
		return
	}

	for client := range gp.clients {
		client.Send(buf)
	}
}
