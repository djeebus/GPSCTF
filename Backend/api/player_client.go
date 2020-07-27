package api

import (
	"encoding/json"
	"github.com/djeebus/gpsctf/Backend/db"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

type Client struct {
	processor *GameProcessor
	conn      *websocket.Conn
	send      chan []byte
	game      *db.Game
	player    *db.Player
}

const (
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
)

func (client *Client) writeMessages() {
	for {
		select {
		case message, ok := <-client.send:
			if !ok {
				err := client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					log.Printf("Failed to write close message: %s", err)
				}
				return
			}

			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Printf("Failed to get next writer: %s", err)
				return
			}
			_, err = w.Write(message)
			if err != nil {
				log.Printf("Failed to write message: %s", err)
				return
			}
		}
	}
}

func (client *Client) readMessages() {
	defer func() {
		client.processor.unregister <- client
		err := client.conn.Close()
		if err != nil {
			log.Printf("Failed to close socket connection: %v", err)
		}
	}()

	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Socket unexpectedly closed: %v", err)
			}
			break
		}

		var latlng LatLngMessage
		err = json.Unmarshal(message, &latlng)
		if err != nil {
			log.Printf("Failed to read lat long message: %v", err)
			return
		}

		err = client.processor.UpdatePlayer(client.player.PlayerID, latlng.Latitude, latlng.Longitude)
		if err != nil {
			log.Printf("Failed to update lat/lng: %v", err)
		}
	}
}

type LatLngMessage struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
