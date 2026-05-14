package downloads

import (
	"context"
	"fmt"
	"strconv"

	"github.com/Abhiram86/echotune/internal"
	"github.com/Abhiram86/echotune/internal/models"
	"github.com/urfave/cli/v3"
)

func removeByIndex(ctx context.Context, c *cli.Command, storage *models.Storage, idx int) error {
	songs := getSortedDownloads(storage)
	if idx < 1 || idx > len(songs) {
		return fmt.Errorf("index out of range")
	}

	if internal.Confirm(fmt.Sprintf("do you want to remove %s?", songs[idx-1].Title)) {
		err := storage.Downloads.Remove(songs[idx-1])
		if err != nil {
			return err
		}
		fmt.Printf("\nRemoved %s\n", songs[idx-1].Title)
	}
	return nil
}

func RemoveByTitle(ctx context.Context, c *cli.Command, storage *models.Storage, query string) error {
	downloaded, err := songByQuery(ctx, storage, query)
	if err != nil {
		return err
	}

	if internal.Confirm(fmt.Sprintf("do you want to remove %s?", downloaded.Title)) {
		err := storage.Downloads.Remove(*downloaded)
		if err != nil {
			return err
		}
		fmt.Printf("\nRemoved %s\n", downloaded.Title)
	}
	return nil
}

func Remove(
	ctx context.Context,
	c *cli.Command,
	storage *models.Storage,
) error {
	query := c.Args().First()
	if query == "" {
		return nil
	}

	if idx, err := strconv.Atoi(query); err == nil {
		return removeByIndex(ctx, c, storage, idx)
	}

	return RemoveByTitle(ctx, c, storage, query)
}
