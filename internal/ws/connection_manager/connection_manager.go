package connectionmanager

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

type ConnectionManager struct {
	connections map[int]*websocket.Conn
	mu sync.RWMutex
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[int]*websocket.Conn),
	}
}

func (cm *ConnectionManager) Add(userId int, conn *websocket.Conn) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.connections[userId] = conn
	fmt.Printf("connections: %v\n", cm.connections)
}

func (cm *ConnectionManager) Get(userId int) (*websocket.Conn, bool) {
	cm.mu.RLock()
	defer cm.mu.RUnlock()

	conn, ok := cm.connections[userId]
	fmt.Printf("connections: %v\n", cm.connections)
	return conn, ok
}

func (cm *ConnectionManager) Remove(userId int) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	delete(cm.connections, userId)
	fmt.Printf("connections: %v\n", cm.connections)
}
