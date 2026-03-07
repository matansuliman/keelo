package modules

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-getter"
)

// Downloader manages the fetching and caching of remote modules.
type Downloader struct {
	cacheDir string
}

// NewDownloader creates a new Downloader.
func NewDownloader(cacheDir string) *Downloader {
	return &Downloader{
		cacheDir: cacheDir,
	}
}

// Download fetches a remote module into the local cache and returns the path to it.
func (d *Downloader) Download(source string) (string, error) {
	// Generate a stable hash using the source URL
	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(source)))[:12]
	dst := filepath.Join(d.cacheDir, hash)

	// Verify if the module is already in cache
	if _, err := os.Stat(dst); err == nil {
		return dst, nil
	}

	// Ensure the base cache directory exists
	if err := os.MkdirAll(d.cacheDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create cache directory: %w", err)
	}

	// Use hashicorp/go-getter to fetch the module
	client := &getter.Client{
		Src:  source,
		Dst:  dst,
		Mode: getter.ClientModeDir,
	}

	if err := client.Get(); err != nil {
		return "", fmt.Errorf("failed to download remote module from %s: %w", source, err)
	}

	return dst, nil
}
