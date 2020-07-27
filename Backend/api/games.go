package api

import (
	"encoding/json"
	"github.com/djeebus/gpsctf/Backend/db"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

var createGameSchema = `{
  "type": "object",
  "additionalProperties": false,
  "properties": {
    "name": {"type": "string"},
	"latitude": {"type": "number"},
	"longitude": {"type": "number"},
	"radius": {"type": "integer"}
  },
  "required": [
    "name",
	"latitude",
	"longitude",
	"radius"
  ]
}`

type CreateGameRequest struct {
	Name string `json:"name"`
	Latitude float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Radius uint8 `json:"radius"`
}

func handleCreateGame(worker *Worker, w http.ResponseWriter, request *http.Request) {
	err, buffer := validateSchema(request, createGameSchema)
	if err != nil {
		renderError(w, err)
		return
	}

	var model CreateGameRequest
	decoder := json.NewDecoder(buffer)
	err = decoder.Decode(&model)
	if err != nil {
		renderError(w, err)
		return
	}

	game, err := db.CreateGame(model.Name, model.Latitude, model.Longitude, model.Radius)
	if err != nil {
		renderError(w, err)
		return
	}

	_, err = worker.AddGame(game)
	if err != nil {
		renderError(w, err)
		return
	}

	RenderJson(w, game)
}

func handleStartGame(worker *Worker, w http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	gameId, err := strconv.ParseInt(vars["gameId"], 10, 64)
	if err != nil {
		renderError(w, err)
		return
	}

	gp, err := worker.StartGame(gameId)
	if err != nil {
		renderError(w, err)
		return
	}

	RenderJson(w, gp.game)
}
