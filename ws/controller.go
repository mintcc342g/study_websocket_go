package ws

import (
	"fmt"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var (
	upgrader = websocket.Upgrader{}
)

// WebSocketHandler ...
type WebSocketHandler struct {
	rc    *redis.Client
	rooms map[string]*Room
}

// NewWebSocketHandler ...
func NewWebSocketHandler(rc *redis.Client) *WebSocketHandler {
	return &WebSocketHandler{
		rc:    rc,
		rooms: make(map[string]*Room),
	}
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

	roomID := c.Param("roomID")
	if roomID == "" {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"result_message": "invalid room_id"}) // TODO: set response format
	}

	userName := c.QueryParam("user_name")
	if userName == "" {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{"result_message": "need a user_name"})
	}

	// TODO: need to check the roomID and userName whether they exist in a db or not

	var room *Room
	var ok bool
	if room, ok = h.rooms[roomID]; !ok { // TODO: prevent that clients would create the same room at the same time
		room = MakeRoom(roomID, h.rc)
		h.rooms = map[string]*Room{roomID: room} // TODO: delete room when the room deleted from a db or all clients left the room
	}

	client := NewClient(wsConn, room, userName)
	client.Enter()

	return nil
}
