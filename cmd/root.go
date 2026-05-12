package cmd

import (
	"context"
	"os"

	"github.com/Abhiram86/echotune/cmd/downloads"
	"github.com/Abhiram86/echotune/internal/models"
	"github.com/urfave/cli/v3"
)

func New(storage *models.Storage) *cli.Command {
	return &cli.Command{
		Name:                  "echotune",
		Usage:                 "echoes audio to the terminal!",
		EnableShellCompletion: true,

		Commands: []*cli.Command{
			{
				Name:  "search",
				Usage: "search and play a song",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:    "auto",
						Aliases: []string{"a"},
						Usage:   "automatically play the relevant result",
					},
					&cli.IntFlag{
						Name:    "limit",
						Aliases: []string{"l"},
						Usage:   "limit the number of results",
					},
					&cli.IntFlag{
						Name:    "repeat",
						Aliases: []string{"r"},
						Usage:   "repeat the search",
					},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					return Search(ctx, c, storage)
				},
			},
			{
				Name:  "history",
				Usage: "show the history of played songs",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "limit",
						Aliases: []string{"l"},
						Usage:   "limit the number of results",
					},
				},
				Action: func(ctx context.Context, c *cli.Command) error {
					return History(ctx, c, storage)
				},
			},
			{
				Name:  "downloads",
				Usage: "manage the downloaded songs",

				Commands: []*cli.Command{
					{
						Name:  "list",
						Usage: "list downloaded songs",
						Flags: []cli.Flag{
							&cli.IntFlag{
								Name:    "limit",
								Aliases: []string{"l"},
								Usage:   "limit the number of results",
							},
							&cli.BoolFlag{
								Name:    "sort",
								Aliases: []string{"s"},
								Usage:   "sort by date",
							},
							&cli.BoolFlag{
								Name:    "sortt",
								Aliases: []string{"st"},
								Usage:   "sort by title",
							},
						},
						Action: func(ctx context.Context, c *cli.Command) error {
							return downloads.List(ctx, c, storage)
						},
					},
					{
						Name:  "play",
						Usage: "play a downloaded song",
						Flags: []cli.Flag{
							&cli.IntFlag{
								Name:    "repeat",
								Aliases: []string{"r"},
								Usage:   "repeat the search",
							},
							&cli.IntFlag{
								Name:    "limit",
								Aliases: []string{"l"},
								Usage:   "play latest n songs",
							},
							&cli.BoolFlag{
								Name:    "shuffle",
								Aliases: []string{"sh"},
								Usage:   "play songs in random order",
							},
						},
						Action: func(ctx context.Context, c *cli.Command) error {
							return downloads.Play(ctx, c, storage, os.Args[2:])
						},
					},
					{
						Name:  "remove",
						Usage: "remove a downloaded song",
						Action: func(ctx context.Context, c *cli.Command) error {
							return downloads.Remove(ctx, c, storage)
						},
					},
				},
			},
			{
				Name:  "clear",
				Usage: "clear cache, history, or all data",
				Action: func(ctx context.Context, c *cli.Command) error {
					return Clear(ctx, c, storage)
				},
			},
		},
	}
}
