package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Abhiram86/echotune/internal"
	"github.com/Abhiram86/echotune/internal/models"
	"github.com/urfave/cli/v3"
)

func Search(c *cli.Command) error {
	query := c.Args().First()
	if query == "" {
		return fmt.Errorf("no query provided")
	}

	searchList, err := internal.SearchQuery(query)

	if err != nil {
		return err
	}

	for idx, result := range searchList.Results {
		fmt.Printf("%d. %s\t%s\n", idx+1, result.Title, result.Channel)
		fmt.Println(result.URL)
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\n\nEnter the number of the song you want to play: ")

	input, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	input = strings.TrimSpace(input)

	songIdx, err := strconv.Atoi(input)
	if err != nil {
		return fmt.Errorf("invalid song index, must be a number")
	}

	if songIdx < 1 || songIdx > len(searchList.Results) {
		return fmt.Errorf("invalid song index")
	}

	fmt.Printf("Playing song %s\n", searchList.Results[songIdx-1].Title)

	player := models.Player{}

	err = internal.PlaySong(&player, searchList.Results[songIdx-1].URL)
	if err != nil {
		return err
	}

	fmt.Println("\nControls:")
	fmt.Print("q - quit\tp - play/pause\tf - forward\tb - backward\t\n")

	for {
		select {
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
