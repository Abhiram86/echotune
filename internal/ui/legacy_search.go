package ui

import (
	"context"
	"fmt"

	"github.com/Abhiram86/echotune/internal/models"
)

func PrintSearchResults(
	ctx context.Context,
	results []models.SearchResult,
) error {
	for idx, result := range results {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			fmt.Printf("%d. %s\t%s\n",
				idx+1,
				result.Title,
				result.Channel,
			)
			fmt.Println(result.URL)
		}
	}

	return nil
}
