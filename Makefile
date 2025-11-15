.PHONY: build run lock branch2

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
