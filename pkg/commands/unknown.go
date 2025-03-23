// pkg/commands/unknown.go
package commands

import (
	"fmt"

	"github.com/cheetahbyte/verba/pkg/context"
)

type UnknownCommand struct {
	Raw string
}

func (u UnknownCommand) Execute(ctx *context.CommandContext) error {
	text := fmt.Sprintf("[Unbekannter Befehl: %s]", u.Raw)

	ctx.PDF.SetXY(ctx.Document.Margin.Left, *ctx.Y)
	ctx.PDF.MultiCell(ctx.Document.TextWidth(), 5, text, "", "L", false)
	*ctx.Y = ctx.PDF.GetY() + 2

	return nil
}

func (c *UnknownCommand) SetArgs(args []string) {
	c.Raw = args[0]
}
