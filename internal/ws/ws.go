package ws

import (
	connectionmanager "connect/internal/ws/connection_manager"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	cm = connectionmanager.NewConnectionManager()
)

func Upgrader() *websocket.Upgrader {
	return &upgrader
}

func ConnectionManager() *connectionmanager.ConnectionManager {
	return cm
}
