package ws

import (
	"fmt"

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
func MakeRoom(roomID string, rc *redis.Client) (room *Room) {
	room = newRoom(roomID, rc)
	fmt.Println("ws", "Room", "New Room Created", "Room", room) // TODO: logger

	go room.run()

	return
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

		case client := <-r.unregister:
			r.unregisterClient(client, done)

		case msg := <-r.broadcast:
			r.publish(msg)
		}
	}
}

// publish ...
func (r *Room) publish(msg []byte) (err error) {
	if err = r.rc.Publish(r.channel, msg).Err(); err != nil {
		fmt.Println("ws", "Room", "publish", "err", err) // TODO: logger
	}
	return
}

// unregisterClient ...
func (r *Room) unregisterClient(client *Client, done chan struct{}) {
	if _, ok := r.clients[client]; ok {
		close(client.pipe)
		delete(r.clients, client)
		done <- struct{}{}
		fmt.Println("ws", "Room", "run", "unregister client", client.nickname) // TODO: logger
	}
}
