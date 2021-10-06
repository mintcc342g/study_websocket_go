package ws

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var (
	upgrader = websocket.Upgrader{}
)

// WebSocketHandler ...
type WebSocketHandler struct {
	RedisClient *redis.Client
}

// ListenBroadcast ...
func (h *WebSocketHandler) ListenBroadcast(c echo.Context) error {

	// TODO: authentication

	wsConn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		fmt.Println("ws", "WebSocketHandler", "Upgrade", "Error", err) // TODO: logger
		return err
	}
	defer wsConn.Close()

	// TODO: need to create ping/pong methods

	roomID := c.Param("roomID")
	if roomID == "" {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"result_message": err.Error()})
	}

	// TODO: needs an entering method
	pubsub := h.RedisClient.Subscribe(roomID)

	_, err = pubsub.Receive()
	if err != nil {
		fmt.Println("ws", "WebSocketHandler", "ListenBroadcast", "Receive", "Error", err) // TODO: logger
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"result_message": err.Error()})
	}

	ch := pubsub.Channel()

	// TODO: needs a writing method
	err = h.RedisClient.Publish(roomID, "hello").Err()
	if err != nil {
		fmt.Println("ws", "WebSocketHandler", "ListenBroadcast", "Publish", "Error", err) // TODO: logger
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"result_message": err.Error()})
	}

	time.AfterFunc(time.Second, func() {
		_ = pubsub.Close()
	})

	// TODO: needs a reading method
	for msg := range ch {
		fmt.Println(msg.Channel, msg.Payload)
	}

	return nil
}
