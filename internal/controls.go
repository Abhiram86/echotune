package internal

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Abhiram86/echotune/internal/models"
)

func Controls(ctx context.Context, player *models.Player, reader *bufio.Reader) error {
	fmt.Println("\nControls:")
	fmt.Print("q - quit\tp - play/pause\tf - forward\tb - backward\t\n")

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
		default:
		}

		input, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		input = strings.TrimSpace(input)

		switch input {
		case "q":
			if player.Cmd.Process != nil {
				player.Cmd.Process.Signal(os.Interrupt)
			}
			return nil
		case "p":
			err = player.TogglePlay()
			if err != nil {
				return err
			}
		case "f":
			err = player.Seek(5)
			if err != nil {
				return err
			}
		case "b":
			err = player.Seek(-5)
			if err != nil {
				return err
			}
		}
	}
}
