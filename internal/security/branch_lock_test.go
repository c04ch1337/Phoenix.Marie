package security

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestBranchLock(t *testing.T) {
	// Setup test environment
	testBranch := "test-branch"
	cleanup := setupTestEnvironment(t)
	defer cleanup()

	// Test creation of new branch lock
	t.Run("NewBranchLock", func(t *testing.T) {
		bl, err := NewBranchLock(testBranch)
		if err != nil {
			t.Fatalf("Failed to create branch lock: %v", err)
		}
		if bl.state.BranchName != testBranch {
			t.Errorf("Expected branch name %s, got %s", testBranch, bl.state.BranchName)
		}
		if bl.state.IsLocked {
			t.Error("New branch lock should not be locked by default")
		}
	})

	// Test lock operation
	t.Run("Lock", func(t *testing.T) {
		bl, _ := NewBranchLock(testBranch)

		if err := bl.Lock(); err != nil {
			t.Fatalf("Failed to lock branch: %v", err)
		}

		if !bl.state.IsLocked {
			t.Error("Branch should be locked after Lock() operation")
		}

		// Test double lock
		if err := bl.Lock(); err == nil {
			t.Error("Expected error when locking already locked branch")
		}
	})

	// Test unlock operation
	t.Run("Unlock", func(t *testing.T) {
		bl, _ := NewBranchLock(testBranch)

		// Lock first
		bl.Lock()

		if err := bl.Unlock(); err != nil {
			t.Fatalf("Failed to unlock branch: %v", err)
		}

		if bl.state.IsLocked {
			t.Error("Branch should be unlocked after Unlock() operation")
		}

		// Test double unlock
		if err := bl.Unlock(); err == nil {
			t.Error("Expected error when unlocking already unlocked branch")
		}
	})

	// Test state verification
	t.Run("VerifyState", func(t *testing.T) {
		bl, _ := NewBranchLock(testBranch)

		if err := bl.VerifyState(); err != nil {
			t.Fatalf("State verification failed: %v", err)
		}

		// Tamper with state
		bl.state.Hash = "invalid_hash"

		if err := bl.VerifyState(); err == nil {
			t.Error("Expected error when verifying tampered state")
		}
	})

	// Test emergency unlock
	t.Run("EmergencyUnlock", func(t *testing.T) {
		bl, _ := NewBranchLock(testBranch)
		bl.Lock()

		// Try with invalid key
		if err := bl.EmergencyUnlock("invalid_key"); err == nil {
			t.Error("Expected error with invalid emergency key")
		}

		// Try with correct key
		validKey := bl.state.EmergencyKey
		if err := bl.EmergencyUnlock(validKey); err != nil {
			t.Fatalf("Emergency unlock failed: %v", err)
		}

		if bl.state.IsLocked {
			t.Error("Branch should be unlocked after emergency unlock")
		}

		// Verify new emergency key was generated
		if bl.state.EmergencyKey == validKey {
			t.Error("Emergency key should be regenerated after use")
		}
	})

	// Test backup system
	t.Run("BackupSystem", func(t *testing.T) {
		bl, _ := NewBranchLock(testBranch)

		// Perform operations to trigger backups
		bl.Lock()
		time.Sleep(time.Millisecond) // Ensure unique timestamp
		bl.Unlock()

		// Check if backup files exist
		backupFiles, err := filepath.Glob(filepath.Join(bl.backupDir, "*.backup"))
		if err != nil {
			t.Fatalf("Failed to list backup files: %v", err)
		}

		if len(backupFiles) < 2 {
			t.Error("Expected at least 2 backup files")
		}
	})
}

// Helper function to setup test environment
func setupTestEnvironment(t *testing.T) func() {
	// Create temporary directories
	tmpDir := filepath.Join(os.TempDir(), "branch_lock_test")
	gitDir := filepath.Join(tmpDir, ".git", "branch_locks")
	backupDir := filepath.Join(gitDir, "backups")

	dirs := []string{tmpDir, gitDir, backupDir}
	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create test directory %s: %v", dir, err)
		}
	}

	// Return cleanup function
	return func() {
		if err := os.RemoveAll(tmpDir); err != nil {
			t.Logf("Warning: Failed to cleanup test directory: %v", err)
		}
	}
}
