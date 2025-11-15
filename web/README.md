# Phase 3.5 Dashboard

A real-time web dashboard for monitoring the Phoenix Marie system components.

## Features

- Real-time system status monitoring
- ORCH agent visualization
- Memory system state tracking
- Emotion system metrics display
- Evolution system status monitoring
- WebSocket-based real-time updates
- Responsive design for various screen sizes
- Basic security with API key authentication

## Architecture

The dashboard consists of three main components:

1. **Frontend**
   - HTML5 interface with responsive design
   - CSS3 for styling and animations
   - JavaScript for real-time data handling and display
   - WebSocket client for live updates

2. **Backend API**
   - RESTful endpoints for system data
   - WebSocket server for real-time updates
   - Metrics collection service
   - Basic authentication middleware

3. **Integration Layer**
   - Real-time metrics collection
   - System state monitoring
   - Data broadcasting system

## API Endpoints

- `GET /api/system/status` - Current system operational status
- `GET /api/orch/metrics` - ORCH agent metrics and status
- `GET /api/memory/state` - Memory system statistics
- `GET /api/emotion/data` - Emotion system metrics
- `GET /api/evolution/stats` - Evolution system statistics
- `WS /ws` - WebSocket endpoint for real-time updates

## Security

The dashboard implements basic security measures:

- API key authentication for REST endpoints
- Public access only to static files and WebSocket connection
- CORS protection
- Rate limiting (TODO)

## Getting Started

1. Start the dashboard server:
   ```bash
   go run cmd/dashboard/main.go
   ```

2. Access the dashboard:
   - Open `http://localhost:8080` in your browser
   - For API access, include the header `X-API-Key: phoenix-dashboard-key`

## WebSocket Data Format

Real-time updates are sent in the following JSON format:

```json
{
  "timestamp": "2025-11-15T21:20:40.531Z",
  "system": {
    "status": "operational",
    "time": "2025-11-15T21:20:40.531Z"
  },
  "orch": {
    "agents": [
      {
        "id": "agent-1",
        "status": "active",
        "taskCount": 5
      }
    ]
  },
  "memory": {
    "totalEntries": 100,
    "activeConnections": 5,
    "cacheHitRate": 95.5
  },
  "emotion": {
    "tone": "calm",
    "pulseRate": 5,
    "responseStyle": "direct"
  },
  "evolution": {
    "generation": 10,
    "populationSize": 100,
    "fitnessScore": 0.85
  }
}
```

## Future Improvements

- Enhanced security measures
- User authentication system
- Historical data visualization
- Advanced metrics and analytics
- Custom alert configurations
- Performance optimization