#!/bin/bash
# PHOENIX.MARIE — MEMORY BACKUP SCRIPT
# Automated backup of memory system for long-term preservation

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$SCRIPT_DIR"

BACKUP_DIR="${BACKUP_DIR:-./data/backups}"
MAX_BACKUPS="${MAX_BACKUPS:-30}"
DATA_DIR="${DATA_DIR:-./data}"

# Create backup directory
mkdir -p "$BACKUP_DIR"

# Generate backup filename with timestamp
TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
BACKUP_FILE="$BACKUP_DIR/phl-memory-backup-$TIMESTAMP.bak"

echo "PHOENIX.MARIE — MEMORY BACKUP"
echo "=============================="
echo "Backup directory: $BACKUP_DIR"
echo "Data directory: $DATA_DIR"
echo ""

# Check if memory database exists
if [ ! -d "$DATA_DIR/phl-memory" ]; then
    echo "⚠️  Memory database not found at $DATA_DIR/phl-memory"
    echo "   No backup created."
    exit 0
fi

# Create backup using BadgerDB backup (if system is running)
# For offline backup, we can copy the directory
echo "Creating backup..."

# Method 1: If Phoenix is running, use programmatic backup
# Method 2: Copy database directory (offline backup)
if [ -d "$DATA_DIR/phl-memory" ]; then
    # Create tar archive of database
    tar -czf "$BACKUP_FILE.tar.gz" -C "$DATA_DIR" phl-memory
    BACKUP_FILE="$BACKUP_FILE.tar.gz"
    echo "✅ Backup created: $BACKUP_FILE"
else
    echo "❌ Failed to create backup"
    exit 1
fi

# Rotate old backups
echo "Rotating backups (keeping last $MAX_BACKUPS)..."

# Count backups
BACKUP_COUNT=$(ls -1 "$BACKUP_DIR"/*.bak* 2>/dev/null | wc -l)

if [ "$BACKUP_COUNT" -gt "$MAX_BACKUPS" ]; then
    # Remove oldest backups
    ls -1t "$BACKUP_DIR"/*.bak* 2>/dev/null | tail -n +$((MAX_BACKUPS + 1)) | xargs rm -f
    echo "✅ Removed old backups"
fi

# Show backup info
BACKUP_SIZE=$(du -h "$BACKUP_FILE" | cut -f1)
echo ""
echo "Backup Statistics:"
echo "  File: $BACKUP_FILE"
echo "  Size: $BACKUP_SIZE"
echo "  Total backups: $BACKUP_COUNT"
echo ""
echo "✅ Backup complete"

