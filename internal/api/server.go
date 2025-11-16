package api

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"golang.org/x/time/rate"
)

// Secure WebSocket configuration
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Only allow connections from our domain
		origin := r.Header.Get("Origin")
		return origin == "http://localhost:8080" || origin == "https://localhost:8080"
	},
}

type Client struct {
	conn     *websocket.Conn
	limiter  *rate.Limiter
	lastSeen time.Time
	mu       sync.Mutex
}

type Server struct {
	clients    map[*websocket.Conn]*Client
	broadcast  chan []byte
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	mu         sync.Mutex
}

func NewServer() *Server {
	s := &Server{
		clients:    make(map[*websocket.Conn]*Client),
		broadcast:  make(chan []byte),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}

	// Start cleanup routine for inactive clients
	go s.cleanupInactiveClients()
	return s
}

func (s *Server) cleanupInactiveClients() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		s.mu.Lock()
		for conn, client := range s.clients {
			if time.Since(client.lastSeen) > 10*time.Minute {
				s.unregister <- conn
			}
		}
		s.mu.Unlock()
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
			s.clients[client] = &Client{
				conn:     client,
				limiter:  rate.NewLimiter(rate.Every(time.Second), 10), // 10 messages per second
				lastSeen: time.Now(),
			}
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
			for conn, client := range s.clients {
				client.mu.Lock()
				if !client.limiter.Allow() {
					client.mu.Unlock()
					continue
				}
				client.lastSeen = time.Now()
				client.mu.Unlock()

				if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
					conn.Close()
					delete(s.clients, conn)
				}
			}
			s.mu.Unlock()
		}
	}
}

func (s *Server) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Validate JWT token for WebSocket connections
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Authentication token required", http.StatusUnauthorized)
		return
	}

	if _, err := validateToken(token); err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not upgrade connection", http.StatusInternalServerError)
		return
	}

	s.register <- conn

	defer func() {
		s.unregister <- conn
		conn.Close()
	}()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}

		client := s.clients[conn]
		client.mu.Lock()
		client.lastSeen = time.Now()
		client.mu.Unlock()
	}
}

// Secure headers middleware
func secureHeadersMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Content-Security-Policy", "default-src 'self'")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		next.ServeHTTP(w, r)
	})
}

// CORS middleware
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// REST Endpoints with authentication
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

	// Apply security middleware to all routes
	secureHandler := secureHeadersMiddleware(mux)
	corsHandler := corsMiddleware(secureHandler)

	// Serve static files with security headers
	mux.Handle("/", http.FileServer(http.Dir("web")))

	// WebSocket endpoint
	mux.HandleFunc("/ws", s.HandleWebSocket)

	// Protected REST API endpoints
	mux.HandleFunc("/api/system/status", authMiddleware(s.HandleSystemStatus))
	mux.HandleFunc("/api/orch/metrics", authMiddleware(s.HandleOrchMetrics))
	mux.HandleFunc("/api/memory/state", authMiddleware(s.HandleMemoryState))
	mux.HandleFunc("/api/emotion/data", authMiddleware(s.HandleEmotionData))
	mux.HandleFunc("/api/evolution/stats", authMiddleware(s.HandleEvolutionStats))

	return corsHandler
}
