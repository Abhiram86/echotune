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

func Controls(ctx context.Context, player *models.Player, storage *models.Storage, reader *bufio.Reader) error {
	fmt.Println("\nControls:")
	fmt.Print("q - quit\tp - play/pause\tf - forward\tb - backward\td - download\n")
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
				return nil
			case "p":
				err := player.TogglePlay()
				if err != nil {
					return err
				}
			case "d":
				err := DownloadSong(ctx, storage, player.Song, &mgr)
				if err != nil {
					log.Printf("download skipped: %v", err)
				}
			case "f":
				err := player.Seek(5)
				if err != nil {
					return err
				}
			case "b":
				err := player.Seek(-5)
				if err != nil {
					return err
				}
			}
		}
	}
}
