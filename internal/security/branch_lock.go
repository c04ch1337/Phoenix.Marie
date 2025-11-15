package security

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// BranchLockState represents the current state of a branch lock
type BranchLockState struct {
	BranchName   string    `json:"branch_name"`
	LockID       string    `json:"lock_id"`
	CreatedAt    time.Time `json:"created_at"`
	LastVerified time.Time `json:"last_verified"`
	Hash         string    `json:"hash"`
	Version      int       `json:"version"`
	BackupPath   string    `json:"backup_path"`
	IsLocked     bool      `json:"is_locked"`
	EmergencyKey string    `json:"emergency_key"`
}

// BranchLock manages branch locking operations
type BranchLock struct {
	mu           sync.RWMutex
	state        BranchLockState
	lockFilePath string
	backupDir    string
	logger       *Logger
}

// Logger handles logging operations
type Logger struct {
	logFile *os.File
}

// NewBranchLock creates a new branch lock instance
func NewBranchLock(branchName string) (*BranchLock, error) {
	lockID := generateLockID(branchName)
	logger, err := newLogger("branch_lock.log")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %v", err)
	}

	bl := &BranchLock{
		lockFilePath: filepath.Join(".git", "branch_locks", fmt.Sprintf("%s.lock", branchName)),
		backupDir:    filepath.Join(".git", "branch_locks", "backups"),
		logger:       logger,
	}

	bl.state = BranchLockState{
		BranchName:   branchName,
		LockID:       lockID,
		CreatedAt:    time.Now(),
		LastVerified: time.Now(),
		Version:      1,
		IsLocked:     false,
		EmergencyKey: generateEmergencyKey(),
	}

	if err := bl.initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize branch lock: %v", err)
	}

	return bl, nil
}

// Lock acquires a lock on the branch
func (bl *BranchLock) Lock() error {
	bl.mu.Lock()
	defer bl.mu.Unlock()

	if bl.state.IsLocked {
		return fmt.Errorf("branch %s is already locked", bl.state.BranchName)
	}

	bl.state.IsLocked = true
	bl.state.LastVerified = time.Now()
	bl.state.Hash = bl.calculateStateHash()

	if err := bl.saveState(); err != nil {
		return fmt.Errorf("failed to save lock state: %v", err)
	}

	bl.logger.Log("info", fmt.Sprintf("Branch %s locked successfully", bl.state.BranchName))
	return nil
}

// Unlock releases the lock on the branch
func (bl *BranchLock) Unlock() error {
	bl.mu.Lock()
	defer bl.mu.Unlock()

	if !bl.state.IsLocked {
		return fmt.Errorf("branch %s is not locked", bl.state.BranchName)
	}

	bl.state.IsLocked = false
	bl.state.LastVerified = time.Now()
	bl.state.Hash = bl.calculateStateHash()

	if err := bl.saveState(); err != nil {
		return fmt.Errorf("failed to save lock state: %v", err)
	}

	bl.logger.Log("info", fmt.Sprintf("Branch %s unlocked successfully", bl.state.BranchName))
	return nil
}

// VerifyState checks the integrity of the lock state
func (bl *BranchLock) VerifyState() error {
	bl.mu.RLock()
	defer bl.mu.RUnlock()

	currentHash := bl.calculateStateHash()
	if currentHash != bl.state.Hash {
		return fmt.Errorf("state integrity check failed")
	}

	return nil
}

// EmergencyUnlock performs an emergency unlock using the emergency key
func (bl *BranchLock) EmergencyUnlock(emergencyKey string) error {
	if emergencyKey != bl.state.EmergencyKey {
		return fmt.Errorf("invalid emergency key")
	}

	bl.mu.Lock()
	defer bl.mu.Unlock()

	bl.state.IsLocked = false
	bl.state.LastVerified = time.Now()
	bl.state.Hash = bl.calculateStateHash()
	bl.state.EmergencyKey = generateEmergencyKey() // Generate new emergency key

	if err := bl.saveState(); err != nil {
		return fmt.Errorf("failed to save state after emergency unlock: %v", err)
	}

	bl.logger.Log("warning", fmt.Sprintf("Emergency unlock performed on branch %s", bl.state.BranchName))
	return nil
}

// Internal helper functions

func (bl *BranchLock) initialize() error {
	if err := os.MkdirAll(filepath.Dir(bl.lockFilePath), 0755); err != nil {
		return fmt.Errorf("failed to create lock directory: %v", err)
	}

	if err := os.MkdirAll(bl.backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %v", err)
	}

	return bl.saveState()
}

func (bl *BranchLock) saveState() error {
	data, err := json.Marshal(bl.state)
	if err != nil {
		return fmt.Errorf("failed to marshal state: %v", err)
	}

	if err := os.WriteFile(bl.lockFilePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write lock file: %v", err)
	}

	// Create backup
	backupPath := filepath.Join(bl.backupDir,
		fmt.Sprintf("%s_%d.backup", bl.state.BranchName, time.Now().Unix()))
	if err := os.WriteFile(backupPath, data, 0644); err != nil {
		return fmt.Errorf("failed to create backup: %v", err)
	}

	bl.state.BackupPath = backupPath
	return nil
}

func (bl *BranchLock) calculateStateHash() string {
	data := fmt.Sprintf("%s:%s:%s:%t:%d",
		bl.state.BranchName,
		bl.state.LockID,
		bl.state.CreatedAt.String(),
		bl.state.IsLocked,
		bl.state.Version)

	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:])
}

func generateLockID(branchName string) string {
	timestamp := time.Now().UnixNano()
	data := fmt.Sprintf("%s:%d", branchName, timestamp)
	hash := sha256.Sum256([]byte(data))
	return hex.EncodeToString(hash[:8])
}

func generateEmergencyKey() string {
	timestamp := time.Now().UnixNano()
	randomData := fmt.Sprintf("%d:%d", timestamp, time.Now().UnixMicro())
	hash := sha256.Sum256([]byte(randomData))
	return hex.EncodeToString(hash[:16])
}

func newLogger(filename string) (*Logger, error) {
	logFile, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open log file: %v", err)
	}

	return &Logger{logFile: logFile}, nil
}

func (l *Logger) Log(level string, message string) {
	timestamp := time.Now().Format(time.RFC3339)
	logEntry := fmt.Sprintf("[%s] %s: %s\n", timestamp, level, message)
	l.logFile.WriteString(logEntry)
}
