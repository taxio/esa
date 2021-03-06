package main

import (
	"fmt"
	"github.com/srvc/fail/v4"
	"os"
	"os/exec"
)

func execEditor(editor, filePath string) error {
	c := exec.Command("sh", "-c", fmt.Sprintf("%s %s", editor, filePath))
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout
	c.Stdin = os.Stdin
	if err := c.Run(); err != nil {
		return fail.Wrap(err)
	}
	return nil
}
