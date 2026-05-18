package downloads

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Abhiram86/echotune/internal"
	"github.com/Abhiram86/echotune/internal/manual"
	"github.com/Abhiram86/echotune/internal/models"
	"github.com/Abhiram86/echotune/internal/operations"
	"github.com/urfave/cli/v3"
)

func getSortedDownloads(storage *models.Storage) []models.Download {
	return internal.SortedDownloads(storage)
}

func songByQuery(ctx context.Context, storage *models.Storage, query string) (*models.Download, error) {
	return internal.FindDownloadByTitle(storage, query)
}

func Play(
	ctx context.Context,
	c *cli.Command,
	storage *models.Storage,
	args []string,
) error {
	query := c.Args().First()

	if idx, err := strconv.Atoi(query); err == nil {
		songs := getSortedDownloads(storage)
		if idx < 1 || idx > len(songs) {
			return fmt.Errorf("index out of range")
		}
		song := songs[idx-1]

		app := internal.NewPlaybackSession(storage, []models.Download{song})
		app.PlayALL(ctx, storage, "a")
		return nil
	}

	if query != "" {
		downloaded, err := songByQuery(ctx, storage, query)
		if err != nil {
			return err
		}

		app := internal.NewPlaybackSession(storage, []models.Download{*downloaded})
		app.PlayALL(ctx, storage, "a")
		return nil
	}

	return PlayAll(ctx, c, storage, args)
}

func PlayAll(
	ctx context.Context,
	c *cli.Command,
	storage *models.Storage,
	args []string,
) error {
	repeat := max(c.Int("repeat"), 1)
	orderedArgs := manual.OrderedArgParse(args)

	songs := getSortedDownloads(storage)

	for _, arg := range orderedArgs {
		switch arg {
		case "shuffle":
			songs = operations.Shuffle(songs)
		case "limit":
			songs = operations.Limit(songs, min(c.Int("limit"), len(songs)))
		}
	}

	app := internal.NewPlaybackSession(storage, songs)

	for range repeat {
		status := app.PlayALL(ctx, storage, "n", "z", "a")
		if status == models.Stopped && app.Queue.CurrentIndex >= len(app.Queue.Songs) {
			app.Queue.CurrentIndex = 0
		} else if status == models.Stopped {
			// Manually interrupted or an error occurred
			break
		}
	}

	return nil
}
