package router

import (
	"log/slog"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins in development
		},
	}
	clients      = make(map[*websocket.Conn]bool)
	clientsMutex sync.RWMutex
)

// HotReloadHandler handles WebSocket connections for hot reload
func HotReloadHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		slog.Error("failed to upgrade websocket connection", "error", err)
		return
	}
	defer conn.Close()

	// Register client
	clientsMutex.Lock()
	clients[conn] = true
	clientsMutex.Unlock()

	slog.Debug("hot reload client connected", "remote_addr", r.RemoteAddr)

	// Keep connection alive and handle disconnect
	defer func() {
		clientsMutex.Lock()
		delete(clients, conn)
		clientsMutex.Unlock()
		slog.Debug("hot reload client disconnected", "remote_addr", r.RemoteAddr)
	}()

	// Read messages (mainly to detect disconnect)
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

// NotifyReload sends reload notification to all connected clients
func NotifyReload() {
	clientsMutex.RLock()
	defer clientsMutex.RUnlock()

	slog.Info("notifying clients to reload", "count", len(clients))

	for conn := range clients {
		err := conn.WriteMessage(websocket.TextMessage, []byte("reload"))
		if err != nil {
			slog.Error("failed to send reload message", "error", err)
		}
	}
}

// StartReloadWatcher watches for reload signal file
func StartReloadWatcher() {
	go func() {
		// This will be triggered by Air's post_cmd
		// For now, it's a placeholder that could watch a signal file
		slog.Debug("hot reload watcher started")
	}()
}

// TriggerReloadIfDev triggers a reload if in development mode
func TriggerReloadIfDev() {
	if os.Getenv("GO_ENV") != "production" {
		NotifyReload()
	}
}
