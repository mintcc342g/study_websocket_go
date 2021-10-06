package ws

import (
	"fmt"

	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
)

// Client ...
type Client struct {
	nickname string
	room     *Room
	conn     *websocket.Conn
	pipe     chan []byte
}

// NewClient ...
func NewClient(conn *websocket.Conn, room *Room, nickname string) *Client {
	client := &Client{
		nickname: nickname,
		conn:     conn,
		room:     room,
		pipe:     make(chan []byte, 256),
	}

	return client
}

// Enter ...
func (c *Client) Enter() {
	// go c.writePump()
	// go c.readPump()

	c.room.register <- c
}

func (c *Client) writePump() {
	// TODO
}

func (c *Client) readPump() {
	// TODO
}

// Subscribe ...
func (c *Client) Subscribe(done <-chan struct{}) {
	pubsub := c.room.rc.Subscribe(c.room.channel)

	ch := pubsub.Channel()

	go func(ch <-chan *redis.Message) {
		for msg := range ch {
			if msg == nil {
				fmt.Println("ws", "Room", "subscribe", "no messages in the channel", c.room.channel) // TODO: logger
				return
			}

			c.pipe <- []byte(msg.Payload)
			fmt.Println("ws", "Room", "subscribe", "received a message", "client", c.nickname, "msg", msg.Payload) // TODO: delete after testing
		}
	}(ch)

	<-done // NOTE: a pubsub channel for a client will be closed when struct{}{} be passed to the done go channel
	err := pubsub.Close()
	fmt.Println("ws", "Room", "subscribe", "closed pubsub channel", "err", err) // TODO: logger
}
