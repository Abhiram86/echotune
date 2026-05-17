package internal

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/Abhiram86/echotune/internal/models"
)

func DownloadSong(ctx context.Context, storage *models.Storage, song models.SearchResult, mgr *models.DownloadManager) error {
	if _, exists := storage.Downloads.Songs[song.ID]; exists {
		return fmt.Errorf("song already downloaded")
	}

	mgr.Mu.Lock()
	if mgr.IsDownloading {
		mgr.Mu.Unlock()
		return fmt.Errorf("download is in progress...")
	}
	mgr.IsDownloading = true
	fmt.Println("Downloading song...")
	mgr.Mu.Unlock()

	go func() {
		defer func() {
			mgr.Mu.Lock()
			mgr.IsDownloading = false
			mgr.Mu.Unlock()
		}()

		outputPath := filepath.Join(
			storage.Downloads.MediaPath,
			song.ID,
			"%(title)s.%(ext)s",
		)
		cmd := exec.CommandContext(
			ctx,
			"yt-dlp",
			"--no-progress",
			"-x",
			"--audio-format", "opus",
			"--audio-quality", "0",
			"--embed-thumbnail",
			"--embed-metadata",
			"--convert-thumbnails", "jpg",
			"--restrict-filenames",
			"-o",
			outputPath,
			song.URL,
		)

		// cmd.Stdout = os.Stdout
		// cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			log.Printf("download failed: %v", err)
			return
		}

		downloaded := models.Download{
			Title:     song.Title,
			Path:      filepath.Join(storage.Downloads.MediaPath, song.ID),
			Metadata:  song,
			Timestamp: time.Now(),
		}

		err := storage.Downloads.Add(downloaded)
		if err != nil {
			log.Printf("failed to add downloaded song: %v", err)
		}

		fmt.Println("Download finished on " + storage.Downloads.MediaPath)
	}()
	return nil
}
