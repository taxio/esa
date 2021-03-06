package main

import (
	"context"
	"github.com/taxio/esa/api"
	"log"
	"os"
)

func main() {
	ctx := context.Background()

	cfg, err := LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	client, err := api.NewClient(cfg.AccessToken, cfg.TeamName)
	if err != nil {
		log.Fatal(err)
	}

	rootCmd := NewRootCmd(client)
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
