package internal

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/Abhiram86/echotune/internal/models"
)

type Control string

const (
	Download           Control = "d"
	Next               Control = "n"
	Previous           Control = "z"
	Repeat             Control = "r"
	AddToPlaylist      Control = "a"
	RemoveFromPlaylist Control = "x"
)

var (
	globalInputChan chan string
	inputOnce       sync.Once
)
var pendingAction string

func Controls(
	ctx context.Context,
	app *PlaybackSession,
	storage *models.Storage,
	reader *bufio.Reader,
	extraControls ...Control,
) error {
	player := app.Player

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
			controlText = append(controlText, "z - previous")

		case Repeat:
			controlText = append(controlText, "r - repeat")

		case AddToPlaylist:
			controlText = append(controlText, "a - add to playlist")

		case RemoveFromPlaylist:
			controlText = append(controlText, "x - remove from playlist")
		}
	}

	fmt.Println("\nControls:")
	fmt.Println(strings.Join(controlText, "\t"))

	mgr := models.DownloadManager{}

	inputOnce.Do(func() {
		globalInputChan = make(chan string)
		go func() {
			for {
				input, err := reader.ReadString('\n')
				if err != nil {
					return
				}
				globalInputChan <- strings.TrimSpace(input)
			}
		}()
	})

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

		case input := <-globalInputChan:

			if pendingAction != "" {
				switch pendingAction {

				case "add_playlist":
					err := app.AddToPlaylist(ctx, storage, input)
					if err != nil {
						fmt.Printf("failed to add song '%s' to playlist '%s': %v\n", app.CurrentSong().Title, input, err)
						return err
					}
					fmt.Printf("added song '%s' to playlist '%s'\n", app.CurrentSong().Title, input)

				case "remove_playlist":
					err := app.RemoveFromPlaylist(ctx, storage, input)
					if err != nil {
						return err
					}
					fmt.Printf("removed song '%s' from playlist '%s'\n", app.CurrentSong().Title, input)
				}

				pendingAction = ""
				continue
			}

			switch input {

			case "q":
				if player.Cmd.Process != nil {
					player.Stop()
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
					err := DownloadSong(ctx, storage, app.CurrentSong().Metadata, &mgr)
					if err != nil {
						log.Printf("download skipped: %v", err)
					}
				}

			case "n":
				if enabled["n"] {
					err := app.Next(ctx)
					if err != nil {
						return err
					}
				}

			case "z":
				if enabled["z"] {
					err := app.Previous(ctx)
					if err != nil {
						return err
					}
				}

			case "r":
				if enabled["r"] {
					fmt.Println("repeat toggled")
				}

			case "a":
				if enabled["a"] {
					fmt.Print("Playlist title: ")
					pendingAction = "add_playlist"
				}

			case "x":
				if enabled["x"] {
					fmt.Print("Playlist title: ")
					pendingAction = "remove_playlist"
				}

			default:
				fmt.Printf("unknown command: %s\n", input)
			}
		}
	}
}
