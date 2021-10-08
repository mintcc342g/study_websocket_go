package ws

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
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
	go c.readPump()
	go c.writePump()

	c.room.register <- c
}

// readPump pumps messages from the websocket connection to the room.
// See https://github.com/gorilla/websocket.
func (c *Client) readPump() {
	defer func() {
		c.room.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))                                                           // NOTE: wait pong
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil }) // NOTE: extend deadline if got pong

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				fmt.Println("ws", "Room", "readPump", "ReadMessage", "err", err) // TODO: logger
			}
			break
		}

		c.broadcast(message) // NOTE: msg's trip ... client server -> client.room.broadcast -> redis pub/sub -> many client.pipe -> clients servers
	}
}

// writePump pumps messages from the room to the websocket connection.
// See https://github.com/gorilla/websocket.
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	// NOTE: This ticker is needed to send 'ping' to a client server from our server.
	// The client server would send a 'pong' if they get the 'ping' signal from ours.
	// Then, our server could get the 'pong' sent by the client at the c.readPump().

	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.pipe: // NOTE: the message was sent from a room.broadcast channel. See Run() in the ws/room.go file.
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message); err != nil {
				fmt.Println("ws", "Client", "writePump", "message delivered", "WriteMessage", "err", err) // TODO: logger
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				fmt.Println("ws", "Client", "writePump", "ping time over", "WriteMessage", "err", err) // TODO: logger
				return
			}
		}
	}
}

// Subscribe ...
func (c *Client) Subscribe(done <-chan struct{}) {
	fmt.Println("ws", "Client", "Subscribe start") // TODO: logger

	pubsub := c.room.rc.Subscribe(c.room.channel)

	ch := pubsub.Channel()
	go c.handleMsg(ch)

	<-done // NOTE: a pubsub channel for a client will be closed when struct{}{} be passed to the done channel
	err := pubsub.Close()
	fmt.Println("ws", "Client", "Subscribe", "closed pubsub channel", "err", err) // TODO: logger
}

func (c *Client) broadcast(msg []byte) {
	c.room.broadcast <- msg
}

func (c *Client) handleMsg(ch <-chan *redis.Message) {
	for msg := range ch {
		if msg == nil {
			fmt.Println("ws", "Client", "handleMsg", "no more messages in the channel", c.room.channel) // TODO: logger
			return
		}

		c.pipe <- []byte(msg.Payload)
		fmt.Println("ws", "Client", "handleMsg", "client", c.nickname, "received a message", msg.Payload) // TODO: logger
	}
}
