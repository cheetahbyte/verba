// pkg/commands/class.go
package commands

import (
	"fmt"

	"github.com/cheetahbyte/verba/pkg/context"
	"github.com/cheetahbyte/verba/pkg/styles"
)

type ClassCommand struct {
	Args []string
}

func (c ClassCommand) Execute(ctx *context.CommandContext) error {
	if len(c.Args) != 1 {
		return fmt.Errorf("::class erwartet genau ein Argument")
	}

	className := c.Args[0]
	class, ok := styles.DocumentStyleClasses[className]
	if !ok {
		return fmt.Errorf("Unbekannte Dokumentklasse: %s", className)
	}

	// apply class settings
	ctx.Document.Margin = class

	ctx.PDF.SetMargins(ctx.Document.Margin.Left, ctx.Document.Margin.Top, ctx.Document.Margin.Right)
	ctx.PDF.SetAutoPageBreak(true, ctx.Document.Margin.Bottom)
	ctx.PDF.SetXY(ctx.Document.Margin.Left, ctx.Document.Margin.Top)
	*ctx.Y = ctx.Document.Margin.Top

	return nil
}

func (c *ClassCommand) SetArgs(args []string) {
	c.Args = args
}
