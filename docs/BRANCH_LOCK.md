# Branch Lock System Documentation

## Overview
The Branch Lock System is a security mechanism that provides controlled access to Git branches, preventing unauthorized modifications and ensuring code integrity. It implements cryptographic verification, state management, and emergency recovery procedures.

## Features
- Branch-level locking mechanism
- Cryptographic state verification
- Automated backup system
- Emergency recovery procedures
- Pre-commit hook integration
- Comprehensive logging

## Components

### 1. Branch Lock Core (`internal/security/branch_lock.go`)
The core implementation provides:
- Branch state management
- Lock/unlock operations
- State verification
- Emergency recovery
- Backup management
- Logging system

### 2. Pre-commit Hook (`.git/hooks/pre-commit`)
Enforces branch lock rules by:
- Checking branch lock status before commits
- Verifying state integrity
- Validating sensitive file modifications
- Logging all validation attempts

## Usage

### Basic Operations

#### Initialize Branch Lock
```go
bl, err := security.NewBranchLock("feature-branch")
if err != nil {
    log.Fatalf("Failed to initialize branch lock: %v", err)
}
```

#### Lock Branch
```go
if err := bl.Lock(); err != nil {
    log.Printf("Failed to lock branch: %v", err)
}
```

#### Unlock Branch
```go
if err := bl.Unlock(); err != nil {
    log.Printf("Failed to unlock branch: %v", err)
}
```

### Emergency Recovery

In case of system failure or lock corruption:

1. Obtain the emergency key from the lock file
2. Use the emergency unlock procedure:
```go
if err := bl.EmergencyUnlock(emergencyKey); err != nil {
    log.Printf("Emergency unlock failed: %v", err)
}
```

## Security Considerations

### State Verification
- All lock operations generate cryptographic hashes
- State integrity is verified before operations
- Tampering attempts are logged and blocked

### Backup System
- Automatic backups created for all state changes
- Backups stored in `.git/branch_locks/backups`
- Each backup includes full state information

### Pre-commit Validation
- Prevents commits to locked branches
- Validates lock state integrity
- Monitors sensitive file modifications
- Logs all validation attempts

## Recovery Procedures

### Standard Recovery
1. Verify branch lock state:
```go
if err := bl.VerifyState(); err != nil {
    // State is corrupted, proceed with recovery
}
```

2. Check backup files in `.git/branch_locks/backups`
3. Use emergency unlock if necessary

### Emergency Recovery
1. Locate the emergency key in the lock file
2. Execute emergency unlock
3. Verify new state integrity
4. Create new backup

## Logging

Logs are stored in:
- Branch lock operations: `branch_lock.log`
- Pre-commit validations: `.git/hooks/branch-lock.log`

Log entries include:
- Timestamp
- Operation type
- Status
- Error messages (if any)

## Best Practices

1. **Regular Verification**
   - Periodically verify branch lock states
   - Monitor logs for unusual activity
   - Test emergency recovery procedures

2. **Backup Management**
   - Regularly clean old backup files
   - Verify backup integrity
   - Document backup locations

3. **Access Control**
   - Limit access to emergency keys
   - Monitor sensitive file modifications
   - Review logs regularly

4. **Integration**
   - Always use pre-commit hooks
   - Implement in CI/CD pipelines
   - Verify state before deployments

## Error Handling

Common errors and solutions:

1. **Lock Acquisition Failure**
   - Verify branch state
   - Check for existing locks
   - Review logs for conflicts

2. **State Verification Failure**
   - Check for file corruption
   - Verify cryptographic hashes
   - Consider emergency recovery

3. **Backup Failures**
   - Verify directory permissions
   - Check available disk space
   - Monitor backup integrity

## Testing

Run the test suite:
```bash
go test ./internal/security -v
```

Tests cover:
- Lock/unlock operations
- State verification
- Emergency recovery
- Backup system
- Pre-commit hooks

## Maintenance

Regular maintenance tasks:

1. **Log Rotation**
   - Implement log rotation
   - Archive old logs
   - Monitor log sizes

2. **Backup Cleanup**
   - Remove old backups
   - Verify backup integrity
   - Document cleanup procedures

3. **State Verification**
   - Regular integrity checks
   - Validate cryptographic hashes
   - Monitor state changes

## Support

For issues or questions:
1. Check logs for error messages
2. Verify state integrity
3. Review backup files
4. Consider emergency recovery
5. Document and report issues

## License
This Branch Lock System is part of the Phoenix Marie project and is subject to its licensing terms.