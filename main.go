package main

import (
	"context"
	"log"
	"os"

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

	if err := cmd.New(&storage).Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
