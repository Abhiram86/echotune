package cmd

import (
	"context"
	"fmt"

	"github.com/Abhiram86/echotune/internal"
	"github.com/Abhiram86/echotune/internal/models"
	"github.com/urfave/cli/v3"
)

func Clear(ctx context.Context, c *cli.Command, storage *models.Storage) error {
	arg := c.Args().First()
	if arg == "" {
		return fmt.Errorf("please specify what to clear: all, cache, or history")
	}

	switch arg {
	case "all":
		if internal.Confirm("This will delete all cache, history, and downloads and playlists. Are you sure?") {
			if err := storage.ClearAll(); err != nil {
				return err
			}
			fmt.Println("Successfully cleared everything.")
		} else {
			fmt.Println("Aborted.")
		}
	case "cache":
		if err := storage.Cache.Clear(); err != nil {
			return err
		}
		fmt.Println("Successfully cleared cache.")
	case "history":
		if internal.Confirm("Are you sure you want to clear your playback history?") {
			if err := storage.History.Clear(); err != nil {
				return err
			}
			fmt.Println("Successfully cleared history.")
		} else {
			fmt.Println("Aborted.")
		}
	default:
		return fmt.Errorf("invalid clear option: %s (choose all, cache, or history)", arg)
	}

	return nil
}
