package cmd

import (
	"context"

	"github.com/Abhiram86/echotune/internal/models"
	"github.com/urfave/cli/v3"
)

func New(storage *models.Storage) *cli.Command {
	return &cli.Command{
		Name:  "echotune",
		Usage: "echoes audio to the terminal!",

		Commands: []*cli.Command{
			{
				Name:  "search",
				Usage: "search for a song",
				Action: func(ctx context.Context, c *cli.Command) error {
					return Search(ctx, c, storage)
				},
			},
			{
				Name:  "history",
				Usage: "show the history of played songs",
				Action: func(ctx context.Context, c *cli.Command) error {
					return History(ctx, c, storage)
				},
			},
			{
				Name:  "clear",
				Usage: "clear the cache",
				Action: func(ctx context.Context, c *cli.Command) error {
					if c.Args().First() == "all" {
						return storage.Clear()
					}
					return nil
				},
			},
		},
	}
}
