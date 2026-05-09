package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Abhiram86/echotune/internal"
	"github.com/Abhiram86/echotune/internal/models"
	"github.com/urfave/cli/v3"
)

func Search(ctx context.Context, c *cli.Command, storage *models.Storage) error {
	query := c.Args().First()
	if query == "" {
		return fmt.Errorf("no query provided")
	}

	cached, searchList, err := internal.SearchQuery(ctx, query, storage)

	if err != nil {
		return err
	}

	if cached != nil {
		searchList = &models.SearchList{
			Query:   query,
			Results: cached.Results,
		}
		fmt.Printf("Using cached results for %s\n", query)
	}

	for idx, result := range searchList.Results {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			fmt.Printf("%d. %s\t%s\n", idx+1, result.Title, result.Channel)
			fmt.Println(result.URL)
		}
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

	err = internal.PlaySong(ctx, &player, searchList.Results[songIdx-1])
	if err != nil {
		return err
	}

	if cached == nil {
		err = storage.Cache.Add(*searchList, songIdx-1)
		if err != nil {
			return err
		}
	}
	err = storage.History.Add(searchList.Results[songIdx-1])
	if err != nil {
		return err
	}

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
