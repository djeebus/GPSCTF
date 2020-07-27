package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/djeebus/gpsctf/Backend/db"
	"github.com/gorilla/websocket"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	db.OpenDatabase()

	m.Run()

	err := db.CloseDatabase()
	if err != nil {
		log.Fatalf("Failed to close db: %v", err)
	}
}

func jsonToReader(model interface{}) (*bytes.Reader, error) {
	buf, err := json.Marshal(model)
	if err != nil {
		return nil, err
	}

	r := bytes.NewReader(buf)
	_, err = r.Seek(0, io.SeekStart)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func makeApiRequest(
	t *testing.T, server *httptest.Server, client *http.Client, method string, path string, requestBody interface{}, responseBody interface{}) {

	var err error
	fmt.Printf("API request: %s %s\n", method, path)

	var requestBytes io.Reader
	if requestBody != nil {
		requestBytes, err = jsonToReader(requestBody)
		if err != nil {
			t.Fatalf("API request: request body: %v", err)
		}
	} else {
		requestBytes = nil
	}

	req, err := http.NewRequest(method, server.URL + path, requestBytes)
	if err != nil {
		t.Fatalf("API request: request: %v", err)
	}

	response, err := client.Do(req)
	if err != nil {
		t.Fatalf("API request: submit: %v", err)
	}

	fmt.Printf("API request: %s %s [%d]", method, path, response.StatusCode)

	if response.StatusCode >= 400 {
		t.Fatalf("API request: status code %d", response.StatusCode)
	}

	if responseBody != nil {
		body, err := ioutil.ReadAll(response.Body)
		err = json.Unmarshal(body, responseBody)
		if err != nil {
			t.Fatalf("API request: unmarshal")
		}
	}
}

func TestProcessTasks(t *testing.T) {
	mux := NewRouter()
	server := httptest.NewServer(mux)
	defer server.Close()

	client := &http.Client{}

	createGame := map[string]interface{}{
		"name": "a test app",
		"latitude": 35,
		"longitude": 35,
		"radius": 0,
	}
	type CreateGameResponse struct {
		GameID int64 `json:"gameId"`
	}
	var game CreateGameResponse
	makeApiRequest(t, server, client, "POST", "/games/", &createGame, &game)

	// add players
	createPlayer1 := map[string]interface{}{
		"name": "player #1",
		"gameID": game.GameID,
	}
	type CreatePlayerResponse struct {
		PlayerID int64 `json:"playerId"`
	}
	var player1 CreatePlayerResponse
	makeApiRequest(t, server, client, "POST", "/players/", &createPlayer1, &player1)

	createPlayer2 := map[string]interface{}{
		"name": "player #2",
		"gameID": game.GameID,
	}
	var player2 CreatePlayerResponse
	makeApiRequest(t, server, client, "POST", "/players/", &createPlayer2, &player2)

	// connect clients to the server
	url := "ws" + strings.TrimPrefix(server.URL, "http")
	url1 := fmt.Sprintf("%s/players/%d/:connect", url, player1.PlayerID)
	ws1, response, err := websocket.DefaultDialer.Dial(url1, nil)
	if err != nil {
		body, _ := ioutil.ReadAll(response.Body)
		t.Fatalf("Failed to dial: %v\n%s", err, string(body))
	}
	defer ws1.Close()

	url2 := fmt.Sprintf("%s/players/%d/:connect", url, player2.PlayerID)
	ws2, response, err := websocket.DefaultDialer.Dial(url2, nil)
	if err != nil {
		body, _ := ioutil.ReadAll(response.Body)
		t.Fatalf("Failed to dial: %v\n%s", err, string(body))
	}
	defer ws2.Close()

	// start the app
	makeApiRequest(t, server, client, "POST",  fmt.Sprintf("/games/%d/:start", game.GameID), nil, nil)

	type Pair struct {
		x *websocket.Conn
		y map[string]float64
	}
	messages := []Pair {
		{ws1, map[string]float64{"latitude": 36, "longitude": 36}},
		{ws2, map[string]float64{"latitude": 36, "longitude": 36}},
		{ws1, map[string]float64{"latitude": 37, "longitude": 37}},
		{ws2, map[string]float64{"latitude": 35, "longitude": 35}},
	}

	for index, item := range messages {
		err = ws1.WriteJSON(item)
		if err != nil {
			t.Fatalf("Failed to send message #%d: %v", index, err)
		}
	}

	msg := make(map[string]interface{})
	err = ws1.ReadJSON(&msg)
	if err != nil {
		t.Fatalf("Failed to read message: %v", err)
	}

	winnerId, ok := msg["playerId"].(float64)
	if !ok {
		t.Fatalf("Failed to find the winner")
	}

	winnerIdInt := int64(winnerId)
	if winnerIdInt != player2.PlayerID {
		t.Fatalf("The wrong player won! %d != %d", winnerIdInt, player2.PlayerID)
	}
}
