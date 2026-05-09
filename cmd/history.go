package cmd

import (
	"context"

	"github.com/Abhiram86/echotune/internal/models"
	"github.com/urfave/cli/v3"
)

func History(ctx context.Context, c *cli.Command, storage *models.Storage) error {
	<-ctx.Done()
	return ctx.Err()
}
