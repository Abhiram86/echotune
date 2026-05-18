package internal

import (
	"fmt"
	"os"

	"github.com/Abhiram86/echotune/internal/models"
	"github.com/Abhiram86/echotune/internal/operations"
)

func FindDownloadByTitle(storage *models.Storage, query string) (*models.Download, error) {
	storage.LoadDownloads()
	song, found := FindBestMatch(storage.Downloads.Songs, func(d models.Download) string {
		return d.Title
	}, query)

	if !found {
		return nil, fmt.Errorf("no downloads matched '%s'", query)
	}
	return &song, nil
}

func FindPlaylistByTitle(storage *models.Storage, query string) (*models.Playlist, error) {
	storage.LoadPlaylists()
	playlist, found := FindBestMatch(storage.Playlists.Playlists, func(p models.Playlist) string {
		return p.Title
	}, query)

	if !found {
		return nil, fmt.Errorf("no playlist matched '%s'", query)
	}
	return &playlist, nil
}

func SortedDownloads(storage *models.Storage) []models.Download {
	storage.LoadDownloads()
	return operations.ToSortedSlice(storage.Downloads.Songs, func(a, b *models.Download) bool {
		return a.Title < b.Title
	})
}

func ResolveDownloadPath(storage *models.Storage, song *models.Download) bool {
	storage.LoadDownloads()
	if downloaded, exists := storage.Downloads.Songs[song.Metadata.ID]; exists {
		if _, err := os.Stat(downloaded.Path); err == nil {
			return true
		}
	}
	return false
}

func ResolvedDownload(storage *models.Storage, song *models.Download) *models.Download {
	storage.LoadDownloads()
	if downloaded, exists := storage.Downloads.Songs[song.Metadata.ID]; exists {
		if _, err := os.Stat(downloaded.Path); err == nil {
			return &downloaded
		}
	}
	return song
}
