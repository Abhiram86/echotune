package cmd

import (
	"context"
	"fmt"

	"github.com/Abhiram86/echotune/internal/models"
	"github.com/urfave/cli/v3"
)

func History(ctx context.Context, c *cli.Command, storage *models.Storage) error {
	for idx, song := range storage.History.Songs {
		fmt.Printf("%d. %s\t%s\n", idx+1, song.Title, song.Channel)
	}
	return nil
}
