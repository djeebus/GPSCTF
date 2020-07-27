package api

import (
	"encoding/json"
	"github.com/djeebus/gpsctf/Backend/db"
	"github.com/djeebus/gpsctf/Backend/app"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"net/http"
	"strconv"
)

var createPlayerSchema = `{
  "type": "object",
  "additionalProperties": false,
  "properties": {
	"gameID": {"type": "integer"},
    "name": {"type": "string"}
  },
  "required": [
	"gameID",
    "name"
  ]
}`

type CreatePlayerRequest struct {
	GameID int64  `json:"gameId"`
	Name   string `json:"name"`
}

func handleCreatePlayer(worker *app.Worker, w http.ResponseWriter, request *http.Request) {
	err, buffer := validateSchema(request, createPlayerSchema)
	if err != nil {
		renderError(w, err)
		return
	}

	var model CreatePlayerRequest
	decoder := json.NewDecoder(buffer)
	err = decoder.Decode(&model)
	if err != nil {
		renderError(w, err)
		return
	}

	gp, err := worker.GetGameProcessor(model.GameID)
	if err != nil {
		renderError(w, err)
		return
	}

	player, err := db.CreatePlayer(model.GameID, model.Name)
	if err != nil {
		renderError(w, err)
		return
	}

	gp.AddPlayer(player)
	renderJson(w, player)
}

// some framework taken from https://github.com/gorilla/websocket/blob/master/examples/chat/client.go
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func handleProcessPlayerMessages(worker *app.Worker, w http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	playerId, err := strconv.ParseInt(vars["playerId"], 10, 64)
	if err != nil {
		renderError(w, err)
		return
	}

	player, err := db.GetPlayer(playerId)
	if err != nil {
		renderError(w, err)
		return
	}

	processor, err := worker.GetGameProcessor(player.GameID)
	if err != nil {
		renderError(w, err)
		return
	}

	conn, err := upgrader.Upgrade(w, request, nil)
	if err != nil {
		renderError(w, err)
		return
	}

	client := &Client{
		processor: processor,
		player:    player,
		conn:      conn,
		send:      make(chan []byte, 256),
	}
	client.processor.Register(client)

	go client.writeMessages()
	go client.readMessages()
}
