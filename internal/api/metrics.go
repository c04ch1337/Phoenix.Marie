package api

import (
	"encoding/json"
	"time"
)

type MetricsService struct {
	server *Server
}

func NewMetricsService(server *Server) *MetricsService {
	return &MetricsService{
		server: server,
	}
}

func (m *MetricsService) Start() {
	go m.collectMetrics()
}

func (m *MetricsService) collectMetrics() {
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		metrics := m.gatherAllMetrics()
		data, err := json.Marshal(metrics)
		if err != nil {
			continue
		}
		m.server.broadcast <- data
	}
}

func (m *MetricsService) gatherAllMetrics() map[string]interface{} {
	return map[string]interface{}{
		"timestamp": time.Now(),
		"system": map[string]interface{}{
			"status": "operational",
			"time":   time.Now(),
		},
		"orch": map[string]interface{}{
			"agents": []map[string]interface{}{
				{
					"id":        "agent-1",
					"status":    "active",
					"taskCount": 5,
				},
			},
		},
		"memory": map[string]interface{}{
			"totalEntries":      100,
			"activeConnections": 5,
			"cacheHitRate":      95.5,
		},
		"emotion": map[string]interface{}{
			"tone":          "calm",
			"pulseRate":     5,
			"responseStyle": "direct",
		},
		"evolution": map[string]interface{}{
			"generation":     10,
			"populationSize": 100,
			"fitnessScore":   0.85,
		},
	}
}
