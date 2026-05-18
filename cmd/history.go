package cmd

import (
	"context"
	"fmt"

	"github.com/Abhiram86/echotune/internal/models"
	"github.com/Abhiram86/echotune/internal/operations"
	"github.com/urfave/cli/v3"
)

func History(ctx context.Context, c *cli.Command, storage *models.Storage) error {
	maxLimit := models.MaxHistory

	if c.Int("limit") > 0 {
		maxLimit = c.Int("limit")
	}

	if err := storage.LoadHistory(); err != nil {
		return err
	}

	songs := storage.History.Songs
	if len(songs) == 0 {
		fmt.Println("No history found.")
		return nil
	}

	if maxLimit > len(songs) {
		maxLimit = len(songs)
	}

	displayList := append([]models.SearchResult(nil), songs...)

	displayList = operations.Reverse(displayList)
	displayList = operations.Limit(displayList, maxLimit)

	for i, song := range displayList {
		fmt.Printf("%d. %s\n", i+1, song.Title)
	}

	return nil
}
