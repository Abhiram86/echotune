package downloads

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/Abhiram86/echotune/internal"
	"github.com/Abhiram86/echotune/internal/models"
	"github.com/urfave/cli/v3"
)

func playSong(ctx context.Context, storage *models.Storage, song models.Download) error {
	player := models.Player{
		Song: song.Metadata,
	}

	fmt.Printf("Playing %s\n", song.Title)

	entries, err := os.ReadDir(song.Path)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", song.Path, err)
	}

	var audioFile string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".opus") {
			audioFile = filepath.Join(song.Path, entry.Name())
			break
		}
	}

	if audioFile == "" {
		return fmt.Errorf("no audio file found in %s", song.Path)
	}

	if err := internal.PlaySong(
		ctx,
		&player,
		models.Playable{
			URL: audioFile,
		},
	); err != nil {
		return err
	}

	reader := bufio.NewReader(os.Stdin)
	return internal.Controls(ctx, &player, storage, reader)
}

func playSongByIndex(
	ctx context.Context,
	storage *models.Storage,
	idx int,
) error {
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

	return playSong(ctx, storage, song)
}

func playSongByTitle(
	ctx context.Context,
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

	return playSong(ctx, storage, *bestSong)
}

func Play(
	ctx context.Context,
	c *cli.Command,
	storage *models.Storage,
) error {
	query := c.Args().First()

	if idx, err := strconv.Atoi(query); err == nil {
		return playSongByIndex(ctx, storage, idx)
	}

	return playSongByTitle(ctx, storage, query)
}
