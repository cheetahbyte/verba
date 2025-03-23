package commands

import (
	"strings"

	"github.com/cheetahbyte/verba/pkg/context"
)

type TextCommand struct {
	Args []string
}

func (t TextCommand) Execute(ctx *context.CommandContext) error {
	text := strings.Join(t.Args, " ")
	ctx.PDF.MultiCell(ctx.Document.TextWidth(), 5, text, "", "J", false)
	*ctx.Y = ctx.PDF.GetY()

	return nil
}

func (c *TextCommand) SetArgs(args []string) {
	c.Args = args
}
