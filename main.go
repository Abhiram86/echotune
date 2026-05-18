package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/Abhiram86/echotune/cmd"
	"github.com/Abhiram86/echotune/internal/models"
)

func main() {
	storage := models.Storage{}
	err := storage.Mount()
	if err != nil {
		log.Fatal(err)
	}

	// defer func() {
	// 	if err := storage.Save(); err != nil {
	// 		log.Fatal(err)
	// 	}
	// }()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	if err := cmd.New(&storage).Run(ctx, os.Args); err != nil {
		log.Fatal(err)
	}
}
