package webserver

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"project_for_tmk_04_06/web"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type TranslationEvent struct {
	Source      string `json:"source"`
	Translation string `json:"translation"`
}

var (
	clients   = make(map[*websocket.Conn]bool)
	clientsMu sync.Mutex
	EventBus  = make(chan TranslationEvent, 10)
	upgrader  = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func StartServer(port int) {
	go processEvents()

	http.Handle("/", http.FileServer(http.FS(web.Assets)))

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			slog.Error("WS Upgrade error", "err", err)
			return
		}

		clientsMu.Lock()
		clients[ws] = true
		clientsMu.Unlock()

		ws.SetPongHandler(func(appData string) error {
			_ = ws.SetReadDeadline(time.Now().Add(60 * time.Second))
			return nil
		})

		_ = ws.WriteJSON(TranslationEvent{Source: "Agent System", Translation: "Connected. Ready to translate."})
	})

	addr := fmt.Sprintf("0.0.0.0:%d", port)
	slog.Info(fmt.Sprintf("Web UI running at http://localhost:%d", port))
	if err := http.ListenAndServe(addr, nil); err != nil {
		slog.Error("ListenAndServe failed", "err", err)
	}
}

func processEvents() {
	for {
		event := <-EventBus
		msg, _ := json.Marshal(event)

		clientsMu.Lock()
		for client := range clients {
			if err := client.WriteMessage(websocket.TextMessage, msg); err != nil {
				client.Close()
				delete(clients, client)
			}
		}
		clientsMu.Unlock()
	}
}
