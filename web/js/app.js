class DashboardManager {
    constructor() {
        this.ws = null;
        this.displays = {
            system: document.querySelector('.status-display'),
            orch: document.querySelector('.orch-display'),
            memory: document.querySelector('.memory-display'),
            emotion: document.querySelector('.emotion-display'),
            evolution: document.querySelector('.evolution-display'),
            exploration: document.getElementById('exploration-log')
        };
        this.initializeWebSocket();
    }

    initializeWebSocket() {
        this.ws = new WebSocket(`ws://${window.location.host}/ws`);
        
        this.ws.onopen = () => {
            console.log('WebSocket connection established');
            this.updateSystemStatus('Connected');
        };

        this.ws.onclose = () => {
            console.log('WebSocket connection closed');
            this.updateSystemStatus('Disconnected');
            // Attempt to reconnect after 5 seconds
            setTimeout(() => this.initializeWebSocket(), 5000);
        };

        this.ws.onerror = (error) => {
            console.error('WebSocket error:', error);
            this.updateSystemStatus('Error');
        };

        this.ws.onmessage = (event) => {
            try {
                const data = JSON.parse(event.data);
                this.handleUpdate(data);
            } catch (error) {
                console.error('Error parsing message:', error);
            }
        };
    }

    handleUpdate(data) {
        switch (data.type) {
            case 'system':
                this.updateSystemStatus(data.status);
                break;
            case 'orch':
                this.updateOrchVisualization(data.agents);
                break;
            case 'memory':
                this.updateMemoryState(data.state);
                break;
            case 'emotion':
                this.updateEmotionMetrics(data.metrics);
                break;
            case 'evolution':
                this.updateEvolutionStatus(data.status);
                break;
            case 'exploration':
                this.updateExplorationLog(data.message);
                break;
            default:
                console.warn('Unknown update type:', data.type);
        }
    }

    updateDisplay(element, content) {
        if (element) {
            element.innerHTML = content;
            element.classList.add('update-flash');
            setTimeout(() => element.classList.remove('update-flash'), 500);
        }
    }

    updateSystemStatus(status) {
        this.updateDisplay(this.displays.system, `
            <div class="status-indicator ${status.toLowerCase()}">
                <h3>Current Status: ${status}</h3>
                <p>Last Updated: ${new Date().toLocaleTimeString()}</p>
            </div>
        `);
    }

    updateOrchVisualization(agents) {
        const agentList = agents.map(agent => `
            <div class="agent-card">
                <h4>Agent ${agent.id}</h4>
                <p>Status: ${agent.status}</p>
                <p>Tasks: ${agent.taskCount}</p>
            </div>
        `).join('');

        this.updateDisplay(this.displays.orch, agentList);
    }

    updateMemoryState(state) {
        this.updateDisplay(this.displays.memory, `
            <div class="memory-stats">
                <p>Total Entries: ${state.totalEntries}</p>
                <p>Active Connections: ${state.activeConnections}</p>
                <p>Cache Hit Rate: ${state.cacheHitRate}%</p>
            </div>
        `);
    }

    updateEmotionMetrics(metrics) {
        this.updateDisplay(this.displays.emotion, `
            <div class="emotion-metrics">
                <p>Current Tone: ${metrics.tone}</p>
                <p>Pulse Rate: ${metrics.pulseRate}</p>
                <p>Response Style: ${metrics.responseStyle}</p>
            </div>
        `);
    }

    updateEvolutionStatus(status) {
        this.updateDisplay(this.displays.evolution, `
            <div class="evolution-stats">
                <p>Generation: ${status.generation}</p>
                <p>Population Size: ${status.populationSize}</p>
                <p>Fitness Score: ${status.fitnessScore}</p>
            </div>
        `);
    }

    updateExplorationLog(message) {
        if (this.displays.exploration) {
            const timestamp = new Date().toLocaleTimeString();
            const entry = `[${timestamp}] ${message}`;
            this.displays.exploration.innerHTML += '<br>' + entry;
            
            // Auto-scroll to bottom
            this.displays.exploration.scrollTop = this.displays.exploration.scrollHeight;
            
            // Limit log entries to last 50
            const lines = this.displays.exploration.innerHTML.split('<br>');
            if (lines.length > 50) {
                this.displays.exploration.innerHTML = lines.slice(-50).join('<br>');
            }
        }
    }
}

// Initialize dashboard when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    window.dashboard = new DashboardManager();
});