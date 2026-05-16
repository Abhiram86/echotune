package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/Abhiram86/echotune/internal"
	"github.com/Abhiram86/echotune/internal/input"
	"github.com/Abhiram86/echotune/internal/models"
	"github.com/Abhiram86/echotune/internal/ui"
	"github.com/urfave/cli/v3"
)

func Search(ctx context.Context, c *cli.Command, storage *models.Storage) error {
	query := c.Args().First()
	repeat := max(c.Int("repeat"), 1)

	if query == "" {
		return fmt.Errorf("no query provided")
	}

	var maxLimit = "10"
	if c.Int("limit") > 0 && c.Int("limit") < 101 {
		maxLimit = strconv.Itoa(c.Int("limit"))
	}
	cached, searchList, err := internal.SearchQuery(ctx, query, storage, maxLimit)

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

	if !c.Bool("auto") {
		err = ui.PrintSearchResults(ctx, searchList.Results)
		if err != nil {
			return err
		}
	}

	var songIdx int

	reader := bufio.NewReader(os.Stdin)
	if c.Bool("auto") {
		songIdx, err = input.SelectBestSong(searchList)
		if err != nil {
			return err
		}
	} else {
		songIdx, err = input.ReadSelection(reader, len(searchList.Results))
	}

	song := searchList.Results[songIdx]

	app := internal.NewPlaybackSession(storage, []models.Download{
		{
			Title:    song.Title,
			Path:     "__SEARCHED__",
			Metadata: song,
		},
	})

	downloaded, exists := storage.Downloads.Songs[song.ID]
	if exists {
		fmt.Printf("Using downloaded song for %s\n", song.Title)
		app.Queue.Songs = []models.Download{downloaded}
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

	for i := range repeat {
		if repeat > 1 {
			fmt.Printf("Playing song (Session %d/%d): %s\n", i+1, repeat, searchList.Results[songIdx].Title)
		}

		status := app.PlayALL(ctx, storage, "d", "a")
		if status == models.Stopped {
			break
		}
	}

	return nil
}
