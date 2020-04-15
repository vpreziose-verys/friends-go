package provider

// Storage ...
type Storage struct {
	StorageConfig
	provider
}

// StorageConfig ...
type StorageConfig struct {
	Config
}

// Close method will be called during teardown
func (p *Storage) Close() {}

// LoadGameStatuses func ...
func (p *Storage) LoadGameStatuses() error {
	return ErrNotImplemented
}

// SaveGameStatuses func ...
func (p *Storage) SaveGameStatuses() error {
	return ErrNotImplemented
}

// RemoveGameStatuses func ...
func (p *Storage) RemoveGameStatuses() error {
	return ErrNotImplemented
}

// LoadReplicationGroup func ...
func (p *Storage) LoadReplicationGroup() error {
	return ErrNotImplemented
}

// SaveReplicationGroup func ...
func (p *Storage) SaveReplicationGroup() error {
	return ErrNotImplemented
}

// GetFile func ...
func (p *Storage) GetFile(filename string) error {
	return ErrNotImplemented
}

// SetFile func ...
func (p *Storage) SetFile(filename string) error {
	return ErrNotImplemented
}
