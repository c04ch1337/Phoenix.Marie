package memory

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// BackupManager handles memory backups
type BackupManager struct {
	storage     *Storage
	backupDir   string
	maxBackups  int
	enabled     bool
}

// BackupConfig holds backup configuration
type BackupConfig struct {
	Enabled        bool
	BackupDir      string
	MaxBackups     int
	ScheduleDaily  bool
	ScheduleWeekly bool
	ScheduleMonthly bool
}

// NewBackupManager creates a new backup manager
func NewBackupManager(storage *Storage, config BackupConfig) *BackupManager {
	if config.BackupDir == "" {
		config.BackupDir = "./data/backups"
	}
	if config.MaxBackups == 0 {
		config.MaxBackups = 30 // Keep 30 backups by default
	}

	bm := &BackupManager{
		storage:    storage,
		backupDir:  config.BackupDir,
		maxBackups: config.MaxBackups,
		enabled:    config.Enabled,
	}

	// Create backup directory
	os.MkdirAll(bm.backupDir, 0755)

	return bm
}

// CreateBackup creates a backup of the memory database
func (bm *BackupManager) CreateBackup() (string, error) {
	if !bm.enabled {
		return "", fmt.Errorf("backup system is disabled")
	}

	timestamp := time.Now().Format("20060102_150405")
	backupPath := filepath.Join(bm.backupDir, fmt.Sprintf("phl-memory-backup-%s.bak", timestamp))

	// Use BadgerDB's native backup
	file, err := os.Create(backupPath)
	if err != nil {
		return "", fmt.Errorf("failed to create backup file: %w", err)
	}
	defer file.Close()

	_, err = bm.storage.db.Backup(file, 0)
	if err != nil {
		return "", fmt.Errorf("failed to create backup: %w", err)
	}

	// Rotate old backups
	bm.rotateBackups()

	return backupPath, nil
}

// RestoreBackup restores memory from a backup file
func (bm *BackupManager) RestoreBackup(backupPath string) error {
	if !bm.enabled {
		return fmt.Errorf("backup system is disabled")
	}

	// Close current database
	if err := bm.storage.Close(); err != nil {
		return fmt.Errorf("failed to close current database: %w", err)
	}

	// Open backup file
	file, err := os.Open(backupPath)
	if err != nil {
		return fmt.Errorf("failed to open backup file: %w", err)
	}
	defer file.Close()

	// Note: Full restore requires closing and reopening the database
	// This is a simplified version - full restore should be done
	// when the system is not running
	return fmt.Errorf("restore requires system shutdown - use backup file manually")
}

// ListBackups returns a list of available backups
func (bm *BackupManager) ListBackups() ([]BackupInfo, error) {
	files, err := os.ReadDir(bm.backupDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read backup directory: %w", err)
	}

	var backups []BackupInfo
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".bak" {
			info, err := file.Info()
			if err != nil {
				continue
			}

			backups = append(backups, BackupInfo{
				Path:     filepath.Join(bm.backupDir, file.Name()),
				Size:     info.Size(),
				Created:  info.ModTime(),
			})
		}
	}

	return backups, nil
}

// rotateBackups removes old backups beyond maxBackups limit
func (bm *BackupManager) rotateBackups() {
	backups, err := bm.ListBackups()
	if err != nil {
		return
	}

	if len(backups) <= bm.maxBackups {
		return
	}

	// Sort by creation time (oldest first)
	// Remove oldest backups
	toRemove := len(backups) - bm.maxBackups
	for i := 0; i < toRemove; i++ {
		os.Remove(backups[i].Path)
	}
}

// BackupInfo contains information about a backup
type BackupInfo struct {
	Path    string
	Size    int64
	Created time.Time
}

// GetBackupStats returns statistics about backups
func (bm *BackupManager) GetBackupStats() BackupStats {
	backups, _ := bm.ListBackups()

	totalSize := int64(0)
	for _, backup := range backups {
		totalSize += backup.Size
	}

	var oldest, newest time.Time
	if len(backups) > 0 {
		oldest = backups[0].Created
		newest = backups[0].Created
		for _, backup := range backups {
			if backup.Created.Before(oldest) {
				oldest = backup.Created
			}
			if backup.Created.After(newest) {
				newest = backup.Created
			}
		}
	}

	return BackupStats{
		TotalBackups: len(backups),
		TotalSize:    totalSize,
		OldestBackup: oldest,
		NewestBackup: newest,
		MaxBackups:   bm.maxBackups,
		Enabled:      bm.enabled,
	}
}

// BackupStats contains backup statistics
type BackupStats struct {
	TotalBackups int
	TotalSize    int64
	OldestBackup time.Time
	NewestBackup time.Time
	MaxBackups   int
	Enabled      bool
}

