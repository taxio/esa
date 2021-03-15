package main

import (
	"io"
	"os"

	"github.com/taxio/esa/log"
)

func mustClose(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Println(err)
		os.Exit(1)
	}
}
