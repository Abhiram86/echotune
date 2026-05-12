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
	"github.com/Abhiram86/echotune/internal/manual"
	"github.com/Abhiram86/echotune/internal/models"
	"github.com/Abhiram86/echotune/internal/operations"
	"github.com/urfave/cli/v3"
)

func playSong(ctx context.Context, storage *models.Storage, song models.Download, repeat int) error {
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

	reader := bufio.NewReader(os.Stdin)

	for i := range repeat {
		if repeat > 1 {
			fmt.Printf("Playing song (Session %d/%d): %s\n", i+1, repeat, song.Title)
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

		err := internal.Controls(ctx, &player, storage, reader)
		if err != nil {
			return err
		}
	}

	return nil
}

func playSongByIndex(
	ctx context.Context,
	storage *models.Storage,
	idx int,
	repeat int,
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

	return playSong(ctx, storage, song, repeat)
}

func playSongByTitle(
	ctx context.Context,
	storage *models.Storage,
	query string,
	repeat int,
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

	return playSong(ctx, storage, *bestSong, repeat)
}

func Play(
	ctx context.Context,
	c *cli.Command,
	storage *models.Storage,
	args []string,
) error {
	query := c.Args().First()
	repeat := max(c.Int("repeat"), 1)

	orderedArgs := manual.OrderedArgParse(args)
	fmt.Println(orderedArgs)

	if len(query) > 0 {
		if idx, err := strconv.Atoi(query); err == nil {
			return playSongByIndex(ctx, storage, idx, repeat)
		}

		err := playSongByTitle(ctx, storage, query, repeat)
		if err != nil {
			return err
		}
	}

	songs := make([]models.Download, 0, len(storage.Downloads.Songs))
	for _, song := range storage.Downloads.Songs {
		songs = append(songs, song)
	}

	for _, arg := range orderedArgs {
		switch arg {
		case "shuffle":
			songs = operations.Shuffle(songs)
		case "limit":
			songs = operations.Limit(songs, min(c.Int("limit"), len(songs)))
		}
	}

	for idx := range repeat {
		if repeat > 1 {
			fmt.Printf("Playing song (Session %d/%d): %s\n", idx+1, repeat, songs[idx].Title)
		}

		for _, song := range songs {
			if err := playSong(ctx, storage, song, 1); err != nil {
				if err.Error() == "interrupted, user quit" {
					return nil
				}
				return err
			}
		}
	}

	return nil
}
