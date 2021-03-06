package main

import (
	"context"
	"fmt"
	"os"
)

func main() {
	ctx := context.Background()
	rootCmd := NewRootCmd()
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(0)
}
