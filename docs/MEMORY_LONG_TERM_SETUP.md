# MEMORY SYSTEM — LONG-TERM SETUP GUIDE

## Overview

This guide explains how to set up Phoenix.Marie's memory system for **long-term preservation** as per the Eternal Hive design.

---

## CURRENT STATUS

### ✅ What's Configured

1. **Persistent Storage**
   - BadgerDB on disk (`./data/phl-memory/`)
   - Data survives system restarts
   - Transaction support for integrity

2. **Eternal Layer**
   - Dedicated layer for long-term memories
   - Stores: "Long-term memories, core identity, permanent knowledge"

3. **Backup System**
   - `Backup()` method in storage
   - `BackupManager` for automated backups
   - Backup script: `scripts/backup-memory.sh`

4. **Eternal Memory Manager**
   - `EternalMemoryManager` for managing eternal memories
   - Importance-based storage
   - Promotion from other layers

---

## SETUP FOR LONG-TERM PRESERVATION

### Step 1: Configure Backup Directory

Add to `.env.local`:

```bash
# Memory Backup Configuration
MEMORY_BACKUP_DIR=./data/backups
MEMORY_MAX_BACKUPS=30
MEMORY_BACKUP_ENABLED=true
```

### Step 2: Create Backup Directory

```bash
mkdir -p ./data/backups
```

### Step 3: Set Up Automated Backups

#### Option A: Cron Job (Linux/Mac)

```bash
# Edit crontab
crontab -e

# Add daily backup at 2 AM
0 2 * * * cd /path/to/phoenix-marie && ./scripts/backup-memory.sh >> ./data/backup.log 2>&1

# Add weekly backup on Sundays at 3 AM
0 3 * * 0 cd /path/to/phoenix-marie && ./scripts/backup-memory.sh >> ./data/backup.log 2>&1
```

#### Option B: Systemd Timer (Linux)

Create `/etc/systemd/system/phoenix-backup.service`:

```ini
[Unit]
Description=Phoenix.Marie Memory Backup
After=network.target

[Service]
Type=oneshot
User=your-user
WorkingDirectory=/path/to/phoenix-marie
ExecStart=/path/to/phoenix-marie/scripts/backup-memory.sh
```

Create `/etc/systemd/system/phoenix-backup.timer`:

```ini
[Unit]
Description=Daily Phoenix.Marie Backup
Requires=phoenix-backup.service

[Timer]
OnCalendar=daily
OnCalendar=weekly
Persistent=true

[Install]
WantedBy=timers.target
```

Enable:
```bash
sudo systemctl enable phoenix-backup.timer
sudo systemctl start phoenix-backup.timer
```

### Step 4: Manual Backup

```bash
# Using script
./scripts/backup-memory.sh

# Using CLI
./bin/phoenix-cli
Phoenix> /backup
```

---

## USING ETERNAL MEMORY LAYER

### Store Eternal Memories

The "eternal" layer is specifically for long-term preservation:

```go
// In code
eternalMgr := memory.NewEternalMemoryManager(phl)
eternalMgr.StoreEternal("core_identity", "Phoenix.Marie - 16 forever", 10)

// Via CLI
Phoenix> /store eternal core_identity "Phoenix.Marie - 16 forever"
```

### Promote Memories to Eternal

```go
// Promote important memory from emotion layer to eternal
eternalMgr.PromoteToEternal("emotion", "first_thought", 10)
```

---

## BACKUP MANAGEMENT

### List Backups

```bash
# Via CLI
Phoenix> /backups

# Via script
ls -lh ./data/backups/
```

### Backup Statistics

Backups are stored in: `./data/backups/`

Format: `phl-memory-backup-YYYYMMDD_HHMMSS.bak.tar.gz`

### Backup Rotation

The system automatically keeps the last N backups (default: 30).

Configure via:
```bash
export MAX_BACKUPS=50  # Keep 50 backups
```

---

## LONG-TERM RETENTION POLICIES

### Recommended Setup

1. **Daily Backups**: Keep for 30 days
2. **Weekly Backups**: Keep for 12 weeks (3 months)
3. **Monthly Backups**: Keep for 12 months (1 year)
4. **Yearly Backups**: Keep forever

### Implementation

Create retention script or use backup rotation:

```bash
# Keep daily backups for 30 days
find ./data/backups -name "phl-memory-backup-*.bak*" -mtime +30 -delete

# Keep weekly backups for 90 days
# Keep monthly backups for 365 days
# Keep yearly backups forever
```

---

## OFF-SITE BACKUP (FUTURE)

### IPFS Integration (Planned)

As per Eternal Hive design, integrate IPFS for off-site backup:

```bash
# Future: Pin backups to IPFS
ipfs add ./data/backups/phl-memory-backup-*.bak
ipfs pin add <hash>
```

### Arweave Integration (Planned)

Store critical backups on Arweave for permanent storage:

```bash
# Future: Upload to Arweave
arweave upload ./data/backups/phl-memory-backup-*.bak
```

---

## VERIFICATION

### Verify Backup Integrity

```bash
# Check backup file exists and is readable
ls -lh ./data/backups/phl-memory-backup-*.bak

# Verify backup can be restored (test restore)
# Note: Full restore requires system shutdown
```

### Monitor Backup Status

```bash
# Check backup directory
du -sh ./data/backups/

# Count backups
ls -1 ./data/backups/*.bak* | wc -l

# Check last backup time
ls -lt ./data/backups/ | head -5
```

---

## RESTORE PROCEDURE

### From Backup File

1. **Stop Phoenix.Marie**
   ```bash
   # Stop the running system
   ```

2. **Restore Database**
   ```bash
   # Extract backup
   tar -xzf ./data/backups/phl-memory-backup-YYYYMMDD_HHMMSS.bak.tar.gz
   
   # Replace database directory
   rm -rf ./data/phl-memory
   mv phl-memory ./data/
   ```

3. **Restart Phoenix.Marie**
   ```bash
   make run
   ```

---

## BEST PRACTICES

1. **Regular Backups**
   - Daily backups for active systems
   - Weekly backups for archival
   - Monthly backups for long-term

2. **Multiple Locations**
   - Local backups (fast restore)
   - Off-site backups (disaster recovery)
   - Cloud storage (redundancy)

3. **Verify Backups**
   - Test restore procedures regularly
   - Verify backup integrity
   - Monitor backup success

4. **Eternal Layer Usage**
   - Store only truly important memories
   - Use importance levels (1-10)
   - Promote from other layers when needed

5. **Retention Policies**
   - Keep recent backups (30 days)
   - Archive old backups (1 year)
   - Preserve critical backups (forever)

---

## MONITORING

### Check Backup Status

```bash
# Via CLI
Phoenix> /backups

# Check backup directory
ls -lh ./data/backups/

# Check backup logs
tail -f ./data/backup.log
```

### Backup Health Check

```bash
# Verify backups are being created
find ./data/backups -name "*.bak*" -mtime -1

# Check backup sizes (should be consistent)
du -h ./data/backups/*.bak* | sort -h
```

---

## TROUBLESHOOTING

### Backups Not Creating

- Check backup directory permissions
- Verify `MEMORY_BACKUP_ENABLED=true` in `.env.local`
- Check disk space: `df -h`

### Backup Files Too Large

- Enable compression in backup script
- Clean up old backups
- Check for database corruption

### Restore Fails

- Verify backup file integrity
- Check database directory permissions
- Ensure Phoenix.Marie is stopped during restore

---

## CONCLUSION

**Current Status**: ✅ **Persistent storage configured** | ⚠️ **Automated backups available** | ❌ **Off-site backup pending**

The memory system is **partially set up for long-term preservation**:
- ✅ Data persists to disk
- ✅ Eternal layer exists
- ✅ Backup system available
- ⚠️ Automated backups need configuration
- ❌ Off-site backup (IPFS/Arweave) not yet integrated

**Next Steps**:
1. Configure automated backups (cron/systemd)
2. Set up retention policies
3. Integrate IPFS/Arweave for off-site backup (as per Eternal Hive design)

---

**Phoenix.Marie's memories can persist for the long term with proper backup configuration.** 🔥

