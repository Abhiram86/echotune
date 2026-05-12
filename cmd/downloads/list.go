package downloads

import (
	"context"
	"fmt"

	"github.com/Abhiram86/echotune/internal/models"
	"github.com/Abhiram86/echotune/internal/operations"
	"github.com/urfave/cli/v3"
)

func List(ctx context.Context, c *cli.Command, storage *models.Storage) error {
	songs := make([]models.Download, 0, len(storage.Downloads.Songs))

	for _, song := range storage.Downloads.Songs {
		songs = append(songs, song)
	}

	if c.Bool("sort") {
		songs = operations.Sort(songs, func(a, b models.Download) bool {
			return a.Timestamp.Before(b.Timestamp)
		})
	} else if c.Bool("sortt") {
		songs = operations.Sort(songs, func(a, b models.Download) bool {
			return a.Metadata.Title < b.Metadata.Title
		})
	} else {
		// Default to sorting by title to ensure consistent order
		songs = operations.Sort(songs, func(a, b models.Download) bool {
			return a.Metadata.Title < b.Metadata.Title
		})
	}

	if c.Int("limit") > 0 {
		songs = operations.Limit(songs, int(c.Int("limit")))
	}

	idx := 1
	for i := len(songs) - 1; i >= 0; i-- {
		fmt.Printf("%d. %s\n", idx, songs[i].Title)
		idx++
	}

	return nil
}
