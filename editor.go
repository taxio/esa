package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/srvc/fail/v4"
)

type Editor interface {
	Exec(ctx context.Context, filePath string) error
	SetEditor(name string)
}

func NewEditor(editor string) Editor {
	return &editorImpl{editor: editor}
}

type editorImpl struct {
	editor string
}

func (e *editorImpl) Exec(ctx context.Context, filePath string) error {
	c := exec.CommandContext(ctx, "sh", "-c", fmt.Sprintf("%s %s", e.editor, filePath))
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout
	c.Stdin = os.Stdin
	if err := c.Run(); err != nil {
		return fail.Wrap(err)
	}
	return nil
}

func (e *editorImpl) SetEditor(name string) {
	e.editor = name
}
