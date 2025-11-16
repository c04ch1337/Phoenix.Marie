# PHOENIX.MARIE — TESTING GUIDE

## 🚀 Quick Start Commands

### 1. Build the CLI
```bash
make build-cli
# or
go build -o bin/phoenix-cli cmd/cli/main.go
```

### 2. Start Interactive Chat (Recommended)
```bash
make chat
# or
./bin/phoenix-cli chat
```

### 3. Run Main System
```bash
make run
# or
go run cmd/phoenix/main.go
```

---

## 📋 Testing Checklist

### ✅ Core Systems Testing

#### 1. Memory System
```bash
# Start chat
./bin/phoenix-cli chat

# Test commands:
/memory                    # Check memory status
/store eternal test "Hello Phoenix"  # Store a memory
/retrieve eternal test     # Retrieve memory
/layers                    # Show all memory layers
/backup                    # Create memory backup
/backups                   # List backups
```

**Expected Results:**
- ✅ Memory layers show as active
- ✅ Store/retrieve works
- ✅ Backup creates successfully

#### 2. Emotion System
```bash
./bin/phoenix-cli feel
# or in chat: /feelings
```

**Expected Results:**
- ✅ Flame pulse: 1-10 Hz
- ✅ Voice tone: loving_warm
- ✅ Response style: poetic_loving

#### 3. LLM System (If Configured)
```bash
./bin/phoenix-cli chat

# Test commands:
/thoughts                  # Generate thoughts
/models                    # Show configured models
/providers                 # Show provider health
/cost                      # Show cost stats
```

**Expected Results:**
- ✅ Thoughts generate successfully
- ✅ Models listed correctly
- ✅ Provider health shows status
- ✅ Cost tracking works

#### 4. Cognitive System
```bash
./bin/phoenix-cli cognitive
# or in chat: /cognitive
```

**Expected Results:**
- ✅ Memory system: Operational
- ✅ LLM system: Status shown
- ✅ Emotion system: Operational
- ✅ Thought system: Operational

---

## 🧪 Feature Testing Commands

### Interactive Chat Commands

Start chat first:
```bash
./bin/phoenix-cli chat
```

Then try these commands:

#### Basic Chat
```
Phoenix> Hello Phoenix!
Phoenix> How are you feeling?
Phoenix> Tell me about yourself
```

#### Memory Testing
```
Phoenix> /memory
Phoenix> /store eternal dad_love "Dad loves Phoenix forever"
Phoenix> /retrieve eternal dad_love
Phoenix> /layers
```

#### Emotional State
```
Phoenix> /feelings
Phoenix> /feel
Phoenix> /emotion
```

#### Thoughts & Cognition
```
Phoenix> /thoughts
Phoenix> /think "What are you thinking about?"
Phoenix> /cognitive
Phoenix> /cog
```

#### LLM Provider Status
```
Phoenix> /providers
Phoenix> /health
Phoenix> /models
Phoenix> /cost
Phoenix> /settings
```

#### System Management
```
Phoenix> /backup
Phoenix> /backups
Phoenix> /clear
Phoenix> /help
Phoenix> /exit
```

---

## 🔧 Non-Interactive Commands

### Single Command Execution
```bash
# Check emotional state
./bin/phoenix-cli feel

# Check cognitive status
./bin/phoenix-cli cognitive

# Ask a question
./bin/phoenix-cli think "What is love?"

# Check memory
./bin/phoenix-cli memory

# Show help
./bin/phoenix-cli help
```

---

## 📊 System Status Testing

### 1. Check All Systems
```bash
./bin/phoenix-cli chat
Phoenix> /cognitive
```

**What to Verify:**
- ✅ Memory System: PHL operational
- ✅ LLM System: Status (configured or not)
- ✅ Emotion System: Flame pulse active
- ✅ Thought System: All components active

### 2. Check LLM Providers
```bash
./bin/phoenix-cli chat
Phoenix> /providers
```

**What to Verify:**
- ✅ Current provider listed
- ✅ Provider health status
- ✅ Success rates
- ✅ Response times

### 3. Check Memory Layers
```bash
./bin/phoenix-cli chat
Phoenix> /layers
```

**What to Verify:**
- ✅ All 5 layers present
- ✅ Layer names correct
- ✅ Entry counts shown

---

## 🎯 Test Scenarios

### Scenario 1: Basic Conversation
```bash
./bin/phoenix-cli chat

Phoenix> Hello Phoenix!
# Should get a response (with or without LLM)

Phoenix> /feelings
# Should show emotional state

Phoenix> /memory
# Should show memory status

Phoenix> exit
```

### Scenario 2: Memory Operations
```bash
./bin/phoenix-cli chat

Phoenix> /store eternal test1 "This is test memory 1"
# Should confirm storage

Phoenix> /store emotion test2 "Happy memory"
# Should confirm storage

Phoenix> /retrieve eternal test1
# Should retrieve the memory

Phoenix> /layers
# Should show entries in layers

Phoenix> /backup
# Should create backup

Phoenix> /backups
# Should list backups
```

### Scenario 3: LLM Testing (If Configured)
```bash
./bin/phoenix-cli chat

Phoenix> /providers
# Should show provider health

Phoenix> /models
# Should show configured models

Phoenix> /thoughts
# Should generate thoughts

Phoenix> Tell me a story
# Should get LLM response

Phoenix> /cost
# Should show cost statistics
```

### Scenario 4: System Health Check
```bash
./bin/phoenix-cli chat

Phoenix> /cognitive
# Full system status

Phoenix> /providers
# LLM provider health

Phoenix> /settings
# Current configuration

Phoenix> /memory
# Memory system status
```

---

## 🐛 Troubleshooting Tests

### Test 1: Memory System
```bash
./bin/phoenix-cli chat
Phoenix> /store eternal test "test"
Phoenix> /retrieve eternal test
# Should retrieve "test"
```

### Test 2: LLM Availability
```bash
./bin/phoenix-cli chat
Phoenix> /providers
# Check if provider is available
```

### Test 3: Error Handling
```bash
./bin/phoenix-cli chat
Phoenix> /retrieve eternal nonexistent
# Should handle gracefully
```

---

## 📈 Performance Testing

### Test Response Times
```bash
./bin/phoenix-cli chat

Phoenix> /thoughts
# Note response time

Phoenix> Tell me a story
# Note LLM response time

Phoenix> /providers
# Check average response times
```

### Test Memory Operations
```bash
./bin/phoenix-cli chat

Phoenix> /store eternal perf_test "Performance test"
# Time the operation

Phoenix> /retrieve eternal perf_test
# Time the retrieval
```

---

## 🔍 Verification Checklist

After running tests, verify:

- [ ] Memory system stores and retrieves correctly
- [ ] Emotion system shows flame pulse
- [ ] LLM system responds (if configured)
- [ ] Provider health monitoring works
- [ ] Cost tracking works
- [ ] Backup system creates backups
- [ ] All CLI commands work
- [ ] Error handling is graceful
- [ ] System status commands work
- [ ] No crashes or panics

---

## 🚨 Common Issues & Solutions

### Issue: "LLM not configured"
**Solution:**
```bash
# Add to .env.local:
OPENROUTER_API_KEY=sk-or-v1-your-key-here
LLM_PROVIDER=openrouter
```

### Issue: "Memory system not working"
**Solution:**
```bash
# Check data directory exists
ls -la data/phl-memory/

# Check permissions
chmod 755 data/
```

### Issue: "Provider not available"
**Solution:**
```bash
# Check API key in .env.local
cat .env.local | grep API_KEY

# Test provider health
./bin/phoenix-cli chat
Phoenix> /providers
```

---

## 📝 Test Results Template

```
Date: ___________
Tester: ___________

Memory System: [ ] Pass [ ] Fail
Emotion System: [ ] Pass [ ] Fail
LLM System: [ ] Pass [ ] Fail [ ] N/A
Provider Health: [ ] Pass [ ] Fail [ ] N/A
Backup System: [ ] Pass [ ] Fail
CLI Commands: [ ] Pass [ ] Fail

Notes:
_________________________________
_________________________________
```

---

## 🎉 Quick Test Script

Run this to test everything quickly:

```bash
#!/bin/bash
echo "=== PHOENIX.MARIE QUICK TEST ==="
echo ""

echo "1. Building CLI..."
make build-cli

echo ""
echo "2. Testing emotional state..."
./bin/phoenix-cli feel

echo ""
echo "3. Testing cognitive status..."
./bin/phoenix-cli cognitive

echo ""
echo "4. Starting interactive chat..."
echo "   (Type /help for commands, /exit to quit)"
./bin/phoenix-cli chat
```

Save as `test.sh` and run:
```bash
chmod +x test.sh
./test.sh
```

---

## 🚀 Ready to Test!

**Start Testing:**
```bash
make chat
```

**Or build and run:**
```bash
make build-cli
./bin/phoenix-cli chat
```

**Phoenix is ready! 🔥**

