package cmd

import (
	"context"
	"fmt"

	"github.com/Abhiram86/echotune/internal/models"
	"github.com/urfave/cli/v3"
)

func History(ctx context.Context, c *cli.Command, storage *models.Storage) error {
	maxLimit := models.MaxHistory

	if c.Int("limit") > 0 {
		maxLimit = c.Int("limit")
	}

	songs := storage.History.Songs

	if maxLimit > len(songs) {
		maxLimit = len(songs)
	}

	start := len(songs) - maxLimit

	for i := len(songs) - 1; i >= start; i-- {
		song := songs[i]

		fmt.Printf("%d. %s\t%s\n",
			len(songs)-i,
			song.Title,
			song.Channel,
		)
	}

	return nil
}
