package types

// LockFile represents the structure of the keelo.lock file.
type LockFile struct {
	Modules []LockedModule `yaml:"modules"`
}

// LockedModule represents a fixed version of a module in the lock file.
type LockedModule struct {
	Name     string `yaml:"name"`
	Source   string `yaml:"source"`
	Resolved string `yaml:"resolved"` // The local hash or version
}
