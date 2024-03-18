package roomHandler

import (
	"bytes"
	"log"
	"time"
	"fmt"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)


// Client is a middleman between the websocket connection and the hub.
type Client struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	PlayerID int
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
		fmt.Println("WS CLOSED")
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		fmt.Println("message from ws", string(message))
		c.hub.broadcast <- message
	}
}


// type PlayerIn struct {
// 	Action string `json:"action"`
// 	PlayerID int `json:"playerId"`
// 	GameID int `json:"gameId"`
// }

type GameStartRes struct {
	Action string `json:"action"`
	PlayerID int `json:"playerId"`
	GameID int `json:"gameId"`
	GameName string `json:"gameName"`
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}


			w.Write(message)

			// if message[len(message) - 1] == 77 {
			// 	before, _ := bytes.CutSuffix(message, []byte(string(message[len(message) - 1])))

			// 	res := &GameStartRes{}
				
			// 	err := json.Unmarshal(before, res)

			// 	if err != nil {
			// 		fmt.Printf("could not marshal %s", err)
			// 		return
			// 	}


			// 	w.Write(before)
			// } else {
				
			// }

			


		
			

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
// func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
// 	conn, err := upgrader.Upgrade(w, r, nil)
	
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}

// 	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
// 	client.hub.register <- client

// 	fmt.Println("hub.clients", hub.clients)

// 	// Allow collection of memory referenced by the caller by doing all work in
// 	// new goroutines.
// 	go client.writePump()
// 	go client.readPump()
// }