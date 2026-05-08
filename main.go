package main

import (
	"context"
	"log"
	"os"

	"github.com/Abhiram86/echotune/cmd"
)

func main() {
	if err := cmd.New().Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
