package actualGame

import (
	"bytes"
	"log"
	"net/http"
	"time"
	"fmt"
	"encoding/json"

	"github.com/gorilla/websocket"
	// "github.com/Toront0/poker/internal/types"
	"github.com/Toront0/poker/internal/types/game"
	"github.com/Toront0/poker/internal/utils"
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

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	Hub *Hub

	// The websocket connection.
	Conn *websocket.Conn

	// Buffered channel of outbound messages.
	Send chan []byte

	PlayerID int
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.

type WSRequest[T any] struct {
	Action string `json:"action"`
	Data T `json:"data"`
}

type WSEmojiRequest struct { 
	EmojiID int `json:"emojiId"`
	SenderId int `json:"senderId"`
}


//aka. check || call - actions that have only one field - playerId
type WSSimpleRequest struct {
	PlayerID int `json:"playerId"` 
}

type WSBetRequest struct {
	PlayerID int `json:"playerId"` 
	Bet int `json:"bet"`
}

type WSShowCardReq struct {
	PlayerID int `json:"playerId"` 
	Hand []string `json:"hand"`
}

type WSSetNextActionReq struct {
	PlayerID int `json:"playerId"` 
	NextAction string `json:"nextAction"`
}



func (c *Client) readPump(gameData *game.PokerTable) {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
		fmt.Println("WS CLOSED")
	}()
	c.Conn.SetReadLimit(maxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error { c.Conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		
		req := &WSRequest[WSEmojiRequest]{}

		err = json.Unmarshal(message, req)
		
		if err != nil {
			fmt.Printf("could not unmarshal websocket 3421 request %s", err)
		}

		if req.Action == "fold" {
			req := &WSRequest[WSSimpleRequest]{}

			err = json.Unmarshal(message, req)
			
			if err != nil {
				fmt.Printf("could not unmarshal websocket request 123 %s", err)
			}	

			fmt.Printf("ws folding")

			gameData.Mutex.Lock()

			idx := utils.SliceIndex(len(gameData.Players), func (i int) bool { return gameData.Players[i].ID == req.Data.PlayerID})

			gameData.Players[idx].Fold()

			gameData.Mutex.Unlock()

		}

		if req.Action == "call" {
			req := &WSRequest[WSBetRequest]{}

			err = json.Unmarshal(message, req)
			
			if err != nil {
				fmt.Printf("could not unmarshal websocket request 123 %s", err)
			}	

			gameData.Mutex.Lock()

			idx := utils.SliceIndex(len(gameData.Players), func (i int) bool { return gameData.Players[i].ID == req.Data.PlayerID})


			gameData.Players[idx].Call(req.Data.Bet)


			gameData.AddChipsToPot(req.Data.Bet)

			gameData.ResetRaisersActionAfterCall()

			gameData.Mutex.Unlock()
		}

		if req.Action == "check" {

			req := &WSRequest[WSSimpleRequest]{}

			err = json.Unmarshal(message, req)
			
			if err != nil {
				fmt.Printf("could not unmarshal websocket request 123 %s", err)
			}	

			gameData.Mutex.Lock()

			idx := utils.SliceIndex(len(gameData.Players), func (i int) bool { return gameData.Players[i].ID == req.Data.PlayerID})

			fmt.Println("CHECKING")

			gameData.Players[idx].Check()

			fmt.Println("AFTER CHECKING")

			gameData.Mutex.Unlock()
		}

		if req.Action == "raise" {

			req := &WSRequest[WSBetRequest]{}

			err = json.Unmarshal(message, req)
			
			if err != nil {
				fmt.Printf("could not unmarshal websocket request 123 %s", err)
			}	

			gameData.Mutex.Lock()

			idx := utils.SliceIndex(len(gameData.Players), func (i int) bool { return gameData.Players[i].ID == req.Data.PlayerID})

			gameData.Players[idx].Raise(req.Data.Bet)
		
			gameData.AddChipsToPot(req.Data.Bet)
			gameData.ResetPlayersActionAfterRaise(req.Data.PlayerID)
			
			gameData.Mutex.Unlock()
		}

		if req.Action == "emoji" {
			c.Hub.Broadcast <- message
		}

		if req.Action == "next-action" {
			req := &WSRequest[WSSetNextActionReq]{}

			err = json.Unmarshal(message, req)
			
			if err != nil {
				fmt.Printf("could not unmarshal websocket request 123 %s", err)
			}	

			gameData.Mutex.Lock()

			idx := utils.SliceIndex(len(gameData.Players), func (i int) bool { return gameData.Players[i].ID == req.Data.PlayerID})

			fmt.Println("req.Data.NextAction", req.Data.NextAction)

			gameData.Players[idx].SetNextAction(req.Data.NextAction)

			gameData.Mutex.Unlock()
		}

		if req.Action == "show-card" {
			// res := &WSRequest[WSShowCardReq]{}



			c.Hub.Broadcast <- message
		}


		fmt.Println("message from ws", string(message))
		// c.Hub.Broadcast <- message
	}
}


type WSResponse struct {
	Action string `json:"action"`
	Data game.PokerTable `json:"data"`
}

type WSRevealCards struct {
	Action string `json:"action"`
	Data []game.PokerPlayer `json:"data"`
}

type WSWinnerResponse struct {
	Action string `json:"action"`
	Winner []int `json:"winner"`
	PotentialStreetCards []string `json:"potentialStreetCards"`
	KickedPlayers []int `json:"kickedPlayers"`
}

type WSGameEndResponse struct {
	Action string `json:"action"`
	Winner game.PokerPlayer `json:"winner"`
}


func (c *Client) writePump(gameData *game.PokerTable) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}


	

			if string(message) == "global-changes" {
				gameData.Mutex.RLock()

				res := &WSResponse{
					Action: "global-changes",
					Data: gameData.SendData(c.PlayerID),
				}

				gameData.Mutex.RUnlock()

				fmt.Println("workds")

				bytes, _ := json.Marshal(res) 

				w.Write(bytes)
			} else if string(message) == "reveal-cards" {
				gameData.Mutex.RLock()

				res := &WSRevealCards{
					Action: "reveal-cards",
					Data: gameData.RevealPlayerCards(),
				}

				gameData.Mutex.RUnlock()
				

				bytes, _ := json.Marshal(res) 

				w.Write(bytes)
			} else if string(message) == "winner" {
				gameData.Mutex.RLock()

				res := &WSWinnerResponse{
					Action: "winner",
					Winner: gameData.GetWinner(),
					PotentialStreetCards: gameData.PotentialStreetCards,
				}

				gameData.Mutex.RUnlock()

				bytes, _ := json.Marshal(res) 

				w.Write(bytes)
			} else if string(message) == "game-end" {
				gameData.Mutex.RLock()

				res := &WSGameEndResponse{
					Action: "game-end",
					Winner: gameData.Players[0],
				}

				gameData.Mutex.RUnlock()

				bytes, _ := json.Marshal(res) 

				w.Write(bytes)
			} else {


				w.Write(message)
			}

			

			// Add queued chat messages to the current websocket message.
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
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