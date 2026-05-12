package internal

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Abhiram86/echotune/internal/models"
)

type Control string

const (
	Download Control = "d"
	Next     Control = "n"
	Previous Control = "v"
	Repeat   Control = "r"
)

func Controls(
	ctx context.Context,
	player *models.Player,
	storage *models.Storage,
	reader *bufio.Reader,
	extraControls ...Control,
) error {

	enabled := map[string]bool{
		"q": true,
		"p": true,
		"f": true,
		"b": true,
	}

	controlText := []string{
		"q - quit",
		"p - play/pause",
		"f - forward",
		"b - backward",
	}

	for _, ctrl := range extraControls {
		enabled[string(ctrl)] = true

		switch ctrl {
		case Download:
			controlText = append(controlText, "d - download")

		case Next:
			controlText = append(controlText, "n - next")

		case Previous:
			controlText = append(controlText, "v - previous")

		case Repeat:
			controlText = append(controlText, "r - repeat")
		}
	}

	fmt.Println("\nControls:")
	fmt.Println(strings.Join(controlText, "\t"))

	mgr := models.DownloadManager{}

	inputChan := make(chan string)

	go func() {
		for {
			input, _ := reader.ReadString('\n')
			inputChan <- strings.TrimSpace(input)
		}
	}()

	for {
		select {

		case <-ctx.Done():
			if player.Cmd.Process != nil {
				player.Cmd.Process.Signal(os.Kill)
			}
			return ctx.Err()

		case <-player.Done:
			fmt.Println("Player stopped")
			return nil

		case input := <-inputChan:

			switch input {

			case "q":
				if player.Cmd.Process != nil {
					player.Cmd.Process.Signal(os.Interrupt)
				}
				return fmt.Errorf("interrupted, user quit")

			case "p":
				if err := player.TogglePlay(); err != nil {
					return err
				}

			case "f":
				if err := player.Seek(5); err != nil {
					return err
				}

			case "b":
				if err := player.Seek(-5); err != nil {
					return err
				}

			case "d":
				if enabled["d"] {
					err := DownloadSong(ctx, storage, player.Song, &mgr)
					if err != nil {
						log.Printf("download skipped: %v", err)
					}
				}

			case "n":
				if enabled["n"] {
					fmt.Println("next song")
				}

			case "v":
				if enabled["v"] {
					fmt.Println("previous song")
				}

			case "r":
				if enabled["r"] {
					fmt.Println("repeat toggled")
				}
			}
		}
	}
}
