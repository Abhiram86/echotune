package platform

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

type AppPaths struct {
	CacheFile        string
	HistoryFile      string
	DownloadFile     string
	DownloadMediaDir string
	PlaylistDir      string
	SocketFile       string
}

func NewAppPaths() (*AppPaths, error) {
	// 1. Resolve base directories once
	baseDataDir, err := getDataDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get data directory: %w", err)
	}

	baseCacheDir, err := os.UserCacheDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get cache directory: %w", err)
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	// 2. Define the specific paths
	paths := &AppPaths{
		CacheFile:    filepath.Join(baseCacheDir, "echotune", "cache.json"),
		HistoryFile:  filepath.Join(baseDataDir, "history.json"),
		DownloadFile: filepath.Join(baseDataDir, "downloads.json"),
		DownloadMediaDir: filepath.Join(
			home,
			"Music",
			"echotune",
		),
		PlaylistDir: filepath.Join(baseDataDir, "playlists"),
		SocketFile:  getSocketPath(), // Use OS-specific IPC path
	}

	return paths, nil
}

func getSocketPath() string {
	if runtime.GOOS == "windows" {
		// mpv on Windows uses Named Pipes for IPC
		return `\\.\pipe\echotune_ipc`
	}

	if xdg := os.Getenv("XDG_RUNTIME_DIR"); xdg != "" {
		return filepath.Join(xdg, "echotune.sock")
	}

	return filepath.Join(os.TempDir(), "echotune.sock")
}

// getDataDir handles OS-specific data directory resolution.
func getDataDir() (string, error) {
	if runtime.GOOS == "linux" {
		if xdg := os.Getenv("XDG_DATA_HOME"); xdg != "" {
			return filepath.Join(xdg, "echotune"), nil
		}
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, ".local", "share", "echotune"), nil
	}

	// Windows: AppData\Roaming, macOS: Library/Application Support
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, "echotune"), nil
}
