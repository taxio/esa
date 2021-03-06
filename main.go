package main

import (
	"context"
	"github.com/spf13/afero"
	"log"
	"os"
)

func main() {
	ctx := context.Background()

	cfg, err := LoadConfig(afero.NewOsFs())
	if err != nil {
		log.Fatal(err)
	}

	client, err := NewClient(cfg.AccessToken, cfg.TeamName)
	if err != nil {
		log.Fatal(err)
	}

	rootCmd := NewRootCmd(cfg, client)
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
