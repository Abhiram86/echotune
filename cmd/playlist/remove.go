package playlist

import (
	"context"
	"fmt"

	"github.com/Abhiram86/echotune/internal"
	"github.com/Abhiram86/echotune/internal/models"
	"github.com/urfave/cli/v3"
)

func Remove(ctx context.Context, c *cli.Command, storage *models.Storage) error {
	if c.Args().Len() == 0 {
		return fmt.Errorf("no playlist specified")
	}

	playlistTitle := c.Args().First()

	playlist, err := searchByQuery(ctx, storage, playlistTitle)
	if err != nil {
		return err
	}

	if c.Args().Len() < 2 {
		if internal.Confirm(fmt.Sprintf("do you want to remove %s?", playlist.Title)) {
			return storage.Playlists.RemovePlayList(*playlist)
		}
		return nil
	}

	bestMatchedSong, found := internal.FindBestMatch(playlist.Songs, func(s models.Download) string {
		return s.Title
	}, c.Args().Get(1))

	if !found {
		return fmt.Errorf("no matches found for '%s'", c.Args().Get(1))
	}

	if internal.Confirm(fmt.Sprintf("do you want to remove %s from %s?", bestMatchedSong.Title, playlist.Title)) {
		return storage.Playlists.RemoveSong(playlist.Title, bestMatchedSong)
	}

	return nil
}
