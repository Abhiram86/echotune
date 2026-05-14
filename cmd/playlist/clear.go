package playlist

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Abhiram86/echotune/internal/models"
	"github.com/urfave/cli/v3"
)

func Clear(ctx context.Context, c *cli.Command, storage *models.Storage) error {
	reader := bufio.NewReader(os.Stdin)

	confirm := func(message string) bool {
		fmt.Printf("%s (y/N): ", message)
		input, _ := reader.ReadString('\n')
		input = strings.ToLower(strings.TrimSpace(input))
		return input == "y" || input == "yes"
	}

	if confirm("This will delete all playlists. Are you sure?") {
		if err := storage.Playlists.ClearAll(); err != nil {
			return err
		}
		fmt.Println("Successfully cleared playlists.")
	}

	return nil
}
