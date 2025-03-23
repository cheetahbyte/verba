package commands

import (
	"fmt"
	"strconv"

	"github.com/cheetahbyte/verba/pkg/context"
)

type MarginCommand struct {
	Args []string
}

func (m MarginCommand) Execute(ctx *context.CommandContext) error {
	if len(m.Args) != 4 {
		return fmt.Errorf("Fehler: ::margin{left, right, top, bottom} erfordert 4 Werte")
	}

	left, err1 := strconv.ParseFloat(m.Args[0], 64)
	right, err2 := strconv.ParseFloat(m.Args[1], 64)
	top, err3 := strconv.ParseFloat(m.Args[2], 64)
	bottom, err4 := strconv.ParseFloat(m.Args[3], 64)

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		return fmt.Errorf("Fehler: Ungültige Werte für ::margin – erwartet wurden Zahlen")
	}

	ctx.Document.Margin.Left = left
	ctx.Document.Margin.Right = right
	ctx.Document.Margin.Top = top
	ctx.Document.Margin.Bottom = bottom

	ctx.PDF.SetMargins(left, top, right)
	ctx.PDF.SetAutoPageBreak(true, bottom)
	ctx.PDF.SetXY(left, top)

	*ctx.Y = top

	return nil
}

func (c *MarginCommand) SetArgs(args []string) {
	c.Args = args
}
