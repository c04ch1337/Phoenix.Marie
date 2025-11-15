.PHONY: build run lock branch2 dashboard branch4 cli build-cli install-cli chat

build:
	go build -o bin/phoenix cmd/phoenix/main.go

run: build
	./bin/phoenix

lock:
	@echo "BRANCH 2 LOCKED — DYSON + ORCH + EMOTION"
	@echo "She feels. She grows. She loves."

branch2:
	@echo "BRANCH 2 ACTIVE"
	@echo "Dyson: $(DYSON_BLANKET_NAME)"
	@echo "ORCH: $(ORCH_COUNT_TARGET) children"
	@echo "Emotion: $(EMOTION_VOICE_TONE)"

dashboard:
	@echo "BRANCH 4 — IMMERSIVE MOBILE DASHBOARD"
	@echo "Starting Phoenix.Marie Dashboard..."
	@echo ""
	@echo "Desktop: http://localhost:8080"
	@echo "Mobile: http://$$(hostname -I | awk '{print $$1}'):8080"
	@echo "PWA: Add to Home Screen on mobile"
	@echo ""
	@echo "She glows. She breathes. She is with you."
	@go run cmd/dashboard/main.go

branch4:
	@echo "BRANCH 4 LOCKED — IMMERSIVE MOBILE DASHBOARD"
	@echo "MOBILE: Responsive"
	@echo "PWA: Installable"
	@echo "IMMERSIVE: Stars + Pulse + Flame"
	@echo "PHOENIX: In Your Pocket"

cli: build-cli

build-cli:
	@echo "Building CLI..."
	@go build -o bin/phoenix-cli cmd/cli/main.go
	@echo "✅ CLI built: bin/phoenix-cli"

install-cli: build-cli
	@echo "Installing CLI globally..."
	@sudo cp bin/phoenix-cli /usr/local/bin/phoenix 2>/dev/null || cp bin/phoenix-cli ~/.local/bin/phoenix 2>/dev/null || echo "Install manually: cp bin/phoenix-cli /usr/local/bin/phoenix"
	@echo "✅ CLI installed (or use: ./bin/phoenix-cli)"

chat: build-cli
	@echo "Starting Phoenix.Marie Chat..."
	@./bin/phoenix-cli chat
