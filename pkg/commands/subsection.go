// pkg/commands/subsection.go
package commands

import "github.com/cheetahbyte/verba/pkg/context"

type SubsectionCommand struct {
	Args []string
}

func (s SubsectionCommand) Execute(ctx *context.CommandContext) error {
	ctx.PDF.SetFont("CMUSerif", "B", 12)
	ctx.PDF.CellFormat(0, 10, s.Args[0], "", 1, "L", false, 0, "")
	ctx.PDF.Ln(4)
	ctx.PDF.SetFont("CMUSerif", "", 11)
	return nil
}

func (c *SubsectionCommand) SetArgs(args []string) {
	c.Args = args
}
