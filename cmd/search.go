package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/Abhiram86/echotune/internal"
	"github.com/Abhiram86/echotune/internal/input"
	"github.com/Abhiram86/echotune/internal/models"
	"github.com/Abhiram86/echotune/internal/ui"
	"github.com/urfave/cli/v3"
)

func Search(ctx context.Context, c *cli.Command, storage *models.Storage) error {
	query := c.Args().First()
	if query == "" {
		return fmt.Errorf("no query provided")
	}

	cached, searchList, err := internal.SearchQuery(ctx, query, storage)

	if err != nil {
		return err
	}

	if cached != nil {
		searchList = &models.SearchList{
			Query:   query,
			Results: cached.Results,
		}
		fmt.Printf("Using cached results for %s\n", query)
	}

	err = ui.PrintSearchResults(ctx, searchList.Results)
	if err != nil {
		return err
	}

	reader := bufio.NewReader(os.Stdin)
	songIdx, err := input.ReadSelection(reader, len(searchList.Results))

	fmt.Printf("Playing song %s\n", searchList.Results[songIdx].Title)

	player := models.Player{}

	err = internal.PlaySong(ctx, &player, searchList.Results[songIdx])
	if err != nil {
		return err
	}

	if cached == nil {
		err = storage.Cache.Add(*searchList, songIdx)
		if err != nil {
			return err
		}
	}
	err = storage.History.Add(searchList.Results[songIdx])
	if err != nil {
		return err
	}

	err = internal.Controls(ctx, &player, reader)
	if err != nil {
		return err
	}

	return nil
}
