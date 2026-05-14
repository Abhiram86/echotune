package playlist

import (
	"context"
	"fmt"

	"github.com/Abhiram86/echotune/internal/models"
	"github.com/Abhiram86/echotune/internal/operations"
	"github.com/urfave/cli/v3"
)

func List(ctx context.Context, c *cli.Command, storage *models.Storage) error {
	query := c.Args().First()

	if query != "" {
		playlist, err := searchByQuery(ctx, storage, query)
		if err != nil {
			return err
		}

		fmt.Printf("Playlist Title: %s\n\n", playlist.Title)

		songs := operations.ToSortedSlice(playlist.Songs, func(a, b *models.Download) bool {
			return a.Title < b.Title
		})

		if c.Int("limit") > 0 {
			songs = operations.Limit(songs, min(c.Int("limit"), len(songs)))
		}

		for i, song := range operations.Reverse(songs) {
			fmt.Printf("%d. %s\n", i+1, song.Title)
		}

		return nil
	}

	playlistsMap := storage.Playlists.Playlists
	playlists := make([]models.Playlist, 0, len(playlistsMap))

	for _, playlist := range playlistsMap {
		playlists = append(playlists, playlist)
	}

	if c.Bool("sort") {
		playlists = operations.Sort(playlists, func(a, b models.Playlist) bool {
			return a.Timestamp.Before(b.Timestamp)
		})
	} else {
		playlists = operations.Sort(playlists, func(a, b models.Playlist) bool {
			return a.Title < b.Title
		})
	}

	if c.Int("limit") > 0 {
		playlists = operations.Limit(playlists, int(c.Int("limit")))
	}

	for i, playlist := range operations.Reverse(playlists) {
		fmt.Printf("%d. %s\n", i+1, playlist.Title)
	}

	return nil
}
