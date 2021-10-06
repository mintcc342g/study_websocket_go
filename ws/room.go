package ws

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

const (
	redisChannel = "ws:study:%s" // TODO: move to another package
)

// Room ...
type Room struct {
	channel    string
	rc         *redis.Client
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

// MakeRoom ...
func MakeRoom(roomID string, rc *redis.Client) *Room {
	room := newRoom(roomID, rc)
	fmt.Println("ws", "Room", "New Room Created", "RoomID", roomID)
	go room.run()

	return room
}

// newRoom ...
func newRoom(roomID string, rc *redis.Client) *Room {
	return &Room{
		channel:    fmt.Sprintf(redisChannel, roomID),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		rc:         rc,
	}
}

func (r *Room) run() {
	done := make(chan struct{}) // NOTE: a signal to close the redis channel after all clients left the redis channel

	for {
		select {
		case client := <-r.register:
			r.clients[client] = true
			go client.Subscribe(done)

			go func() { // TODO: delete these temporal codes after creating read/writePump methods
				time.Sleep(60 * time.Second)
				r.unregister <- client
			}()

		case client := <-r.unregister:
			if _, ok := r.clients[client]; ok {
				close(client.pipe)
				delete(r.clients, client)
				done <- struct{}{}
				fmt.Println("ws", "Room", "run", "unregister client", client.nickname) // TODO: logger
			}

		case message := <-r.broadcast:
			for client := range r.clients {
				select {
				case client.pipe <- message:
				default:
					close(client.pipe)
					delete(r.clients, client)
				}
			}
		}
	}
}
