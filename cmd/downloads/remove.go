package downloads

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/Abhiram86/echotune/internal"
	"github.com/Abhiram86/echotune/internal/models"
	"github.com/urfave/cli/v3"
)

func removeByIndex(ctx context.Context, c *cli.Command, storage *models.Storage, idx int) error {
	if idx < 1 || idx > len(storage.Downloads.Songs) {
		return fmt.Errorf("index out of range")
	}

	songs := make([]models.Download, 0, len(storage.Downloads.Songs))

	for _, song := range storage.Downloads.Songs {
		songs = append(songs, song)
	}

	sort.Slice(songs, func(i, j int) bool {
		return songs[i].Title < songs[j].Title
	})

	song := songs[idx-1]

	fmt.Printf("do you want to uninstall %s? [y/N]: ", song.Title)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.ToLower(strings.TrimSpace(input))
	if input != "y" && input != "yes" {
		return nil
	}

	err := storage.Downloads.Remove(song)
	if err != nil {
		return err
	}

	fmt.Printf("\nRemoved %s\n", song.Title)
	return nil
}

func RemoveByTitle(
	ctx context.Context,
	c *cli.Command,
	storage *models.Storage,
	query string,
) error {
	var bestSong *models.Download
	highScore := -1

	// Iterate map values
	for _, song := range storage.Downloads.Songs {
		s := internal.Score(song.Title, query)
		if s > highScore && s > 0 {
			highScore = s
			// We need to take a local copy or pointer to the current song
			currentSong := song
			bestSong = &currentSong
		}
	}

	if bestSong == nil {
		return fmt.Errorf("no matches found for '%s'", query)
	}

	fmt.Printf("do you want to uninstall %s? [y/N]: ", bestSong.Title)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.ToLower(strings.TrimSpace(input))
	if input != "y" && input != "yes" {
		return nil
	}

	return storage.Downloads.Remove(*bestSong)
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
