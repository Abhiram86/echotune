package cmd

import (
	"context"

	"github.com/urfave/cli/v3"
)

func New() *cli.Command {
	return &cli.Command{
		Name:  "echotune",
		Usage: "echoes audio to the terminal!",

		Commands: []*cli.Command{
			{
				Name:  "search",
				Usage: "search for a song",
				Action: func(ctx context.Context, c *cli.Command) error {
					return Search(c)
				},
			},
			{
				Name:  "history",
				Usage: "show the history of played songs",
				Action: func(ctx context.Context, c *cli.Command) error {
					return History(c)
				},
			},
		},
	}
}
