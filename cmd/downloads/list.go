package downloads

import (
	"context"
	"fmt"
	"sort"

	"github.com/Abhiram86/echotune/internal/models"
	"github.com/urfave/cli/v3"
)

func List(ctx context.Context, c *cli.Command, storage *models.Storage) error {
	songs := make([]models.Download, 0, len(storage.Downloads.Songs))

	for _, song := range storage.Downloads.Songs {
		songs = append(songs, song)
	}

	sort.Slice(songs, func(i, j int) bool {
		return songs[i].Title < songs[j].Title
	})

	idx := 1
	for i := len(songs) - 1; i >= 0; i-- {
		fmt.Printf("%d. %s\n", idx, songs[i].Title)
		idx++
	}

	return nil
}
