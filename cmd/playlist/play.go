package playlist

import (
	"context"
	"fmt"

	"github.com/Abhiram86/echotune/internal"
	"github.com/Abhiram86/echotune/internal/manual"
	"github.com/Abhiram86/echotune/internal/models"
	"github.com/Abhiram86/echotune/internal/operations"
	"github.com/urfave/cli/v3"
)

func searchByQuery(_ context.Context, storage *models.Storage, query string) (*models.Playlist, error) {
	return internal.FindPlaylistByTitle(storage, query)
}

func Play(ctx context.Context, c *cli.Command, storage *models.Storage, args []string) error {
	query := c.Args().First()
	repeat := max(c.Int("repeat"), 1)
	orderedArgs := manual.OrderedArgParse(args)

	if query == "" {
		return fmt.Errorf("no query provided")
	}

	playlist, err := searchByQuery(ctx, storage, query)
	if err != nil {
		return err
	}

	app := internal.NewPlaybackSession(storage, []models.Download{})

	for _, song := range playlist.Songs {
		resolved := internal.ResolvedDownload(storage, &song)
		app.Queue.Songs = append(app.Queue.Songs, *resolved)
	}

	for i := range orderedArgs {
		switch orderedArgs[i] {
		case "shuffle":
			app.Queue.Songs = operations.Shuffle(app.Queue.Songs)
		case "limit":
			app.Queue.Songs = operations.Limit(app.Queue.Songs, min(c.Int("limit"), len(app.Queue.Songs)))
		}
	}

	for range repeat {
		status := app.PlayALL(ctx, storage, "n", "z", "x")
		if status == models.Stopped && app.Queue.CurrentIndex >= len(app.Queue.Songs) {
			app.Queue.CurrentIndex = 0
		} else if status == models.Stopped {
			break
		}
	}

	return nil
}
