package downloads

import (
	"context"
	"fmt"

	"github.com/Abhiram86/echotune/internal/models"
	"github.com/Abhiram86/echotune/internal/operations"
	"github.com/urfave/cli/v3"
)

func List(ctx context.Context, c *cli.Command, storage *models.Storage) error {
	songs := operations.ToSortedSlice(storage.Downloads.Songs, func(a, b *models.Download) bool {
		if c.Bool("sort") {
			return a.Timestamp.Before(b.Timestamp)
		}
		return a.Metadata.Title < b.Metadata.Title
	})

	if c.Int("limit") > 0 {
		songs = operations.Limit(songs, int(c.Int("limit")))
	}

	for i, song := range operations.Reverse(songs) {
		fmt.Printf("%d. %s\n", i+1, song.Title)
	}

	return nil
}
