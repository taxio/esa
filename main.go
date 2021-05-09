package main

import (
	"context"
	"os"

	"github.com/taxio/esa/log"
)

func main() {
	rootCmd := NewRootCmd()
	if err := rootCmd.ExecuteContext(context.Background()); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}
