// pkg/commands/section.go
package commands

import "github.com/cheetahbyte/verba/pkg/context"

type SectionCommand struct {
	Args []string
}

func (s SectionCommand) Execute(ctx *context.CommandContext) error {
	ctx.PDF.SetFont("CMUSerif", "B", 14)
	ctx.PDF.SetY(*ctx.Y)
	ctx.PDF.CellFormat(0, 10, s.Args[0], "", 1, "L", false, 0, "")
	*ctx.Y += 10 // ❗️Y verschieben!
	ctx.PDF.SetFont("CMUSerif", "", 11)
	return nil
}

func (c *SectionCommand) SetArgs(args []string) {
	c.Args = args
}
