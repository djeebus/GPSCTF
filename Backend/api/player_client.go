package api

import (
	"encoding/json"
	"github.com/djeebus/gpsctf/Backend/app"
	"github.com/djeebus/gpsctf/Backend/db"
	"github.com/gorilla/websocket"
	"log"
)

type Client struct {
	processor *app.GameProcessor
	conn      *websocket.Conn
	send      chan []byte
	game      *db.Game
	player    *db.Player
}

func (c *Client) Send(buffer []byte) {
	c.send <- buffer
}

func (c *Client) Close() {
	close(c.send)
}

func (c *Client) writeMessages() {
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				err := c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					log.Printf("Failed to write close message: %s", err)
				}
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
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

func (c *Client) readMessages() {
	defer func() {
		c.processor.Unregister(c)
		err := c.conn.Close()
		if err != nil {
			log.Printf("Failed to close socket connection: %v", err)
		}
	}()

	for {
		_, message, err := c.conn.ReadMessage()
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

		err = c.processor.UpdatePlayer(c.player.PlayerID, latlng.Latitude, latlng.Longitude)
		if err != nil {
			log.Printf("Failed to update lat/lng: %v", err)
		}
	}
}

type LatLngMessage struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}
