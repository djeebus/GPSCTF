package api

import (
	"github.com/gorilla/mux"
	"net/http"
)

func NewRouter() *mux.Router {
	r := mux.NewRouter()

	worker := &Worker{
		games: make(map[int64]*GameProcessor),
	}

	r.Path("/status").HandlerFunc(handleStatus)

	gamesRouter := r.Path("/games/").Subrouter()
	gamesRouter.Methods("POST").HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			handleCreateGame(worker, w, r)
		})

	r.Path("/games/{gameId}/start").Methods("POST").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleStartGame(worker, w, r)
	})

	r.Path("/players/").Methods("POST").HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			handleCreatePlayer(worker, w, r)
		})
	r.Path("/players/{playerId}/:connect").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handleProcessPlayerMessages(worker, w, r)
	})
	return r
}
