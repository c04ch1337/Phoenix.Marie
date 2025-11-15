package api

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // In production, this should be more restrictive
	},
}

type Server struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan []byte
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	mu         sync.Mutex
}

func NewServer() *Server {
	return &Server{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}
}

func (s *Server) Start() {
	go s.run()
}

func (s *Server) run() {
	for {
		select {
		case client := <-s.register:
			s.mu.Lock()
			s.clients[client] = true
			s.mu.Unlock()

		case client := <-s.unregister:
			s.mu.Lock()
			if _, ok := s.clients[client]; ok {
				delete(s.clients, client)
				client.Close()
			}
			s.mu.Unlock()

		case message := <-s.broadcast:
			s.mu.Lock()
			for client := range s.clients {
				if err := client.WriteMessage(websocket.TextMessage, message); err != nil {
					client.Close()
					delete(s.clients, client)
				}
			}
			s.mu.Unlock()
		}
	}
}

func (s *Server) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not upgrade connection", http.StatusInternalServerError)
		return
	}

	s.register <- conn

	// Clean up on disconnect
	defer func() {
		s.unregister <- conn
		conn.Close()
	}()

	// Keep connection alive
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

// REST Endpoints

func (s *Server) HandleSystemStatus(w http.ResponseWriter, r *http.Request) {
	status := map[string]interface{}{
		"status": "operational",
		"time":   time.Now(),
	}
	json.NewEncoder(w).Encode(status)
}

func (s *Server) HandleOrchMetrics(w http.ResponseWriter, r *http.Request) {
	metrics := map[string]interface{}{
		"agents": []map[string]interface{}{
			{
				"id":        "agent-1",
				"status":    "active",
				"taskCount": 5,
			},
		},
	}
	json.NewEncoder(w).Encode(metrics)
}

func (s *Server) HandleMemoryState(w http.ResponseWriter, r *http.Request) {
	state := map[string]interface{}{
		"totalEntries":      100,
		"activeConnections": 5,
		"cacheHitRate":      95.5,
	}
	json.NewEncoder(w).Encode(state)
}

func (s *Server) HandleEmotionData(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{
		"tone":          "calm",
		"pulseRate":     5,
		"responseStyle": "direct",
	}
	json.NewEncoder(w).Encode(data)
}

func (s *Server) HandleEvolutionStats(w http.ResponseWriter, r *http.Request) {
	stats := map[string]interface{}{
		"generation":     10,
		"populationSize": 100,
		"fitnessScore":   0.85,
	}
	json.NewEncoder(w).Encode(stats)
}

func (s *Server) SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	// Serve static files
	mux.Handle("/", http.FileServer(http.Dir("web")))

	// WebSocket endpoint
	mux.HandleFunc("/ws", s.HandleWebSocket)

	// REST API endpoints
	mux.HandleFunc("/api/system/status", s.HandleSystemStatus)
	mux.HandleFunc("/api/orch/metrics", s.HandleOrchMetrics)
	mux.HandleFunc("/api/memory/state", s.HandleMemoryState)
	mux.HandleFunc("/api/emotion/data", s.HandleEmotionData)
	mux.HandleFunc("/api/evolution/stats", s.HandleEvolutionStats)

	return mux
}
