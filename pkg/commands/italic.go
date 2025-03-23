package commands

import (
	"strings"

	"github.com/cheetahbyte/verba/pkg/context"
)

type ItalicCommand struct {
	Args []string
}

func (b ItalicCommand) Execute(ctx *context.CommandContext) error {
	ctx.PDF.SetFont("CMUSerif", "I", 11)
	text := strings.Join(b.Args, " ")
	if strings.HasSuffix(text, "\\") {
		text = strings.TrimSuffix(text, "\\")
		ctx.PDF.MultiCell(ctx.Document.TextWidth(), 5, text, "", "L", false)
		*ctx.Y = ctx.PDF.GetY() + 5
	} else {
		ctx.PDF.MultiCell(ctx.Document.TextWidth(), 5, text, "", "L", false)
		*ctx.Y = ctx.PDF.GetY()
	}
	ctx.PDF.SetFont("CMUSerif", "", 11)
	return nil
}

func (b ItalicCommand) InlineText(ctx *context.CommandContext) string {
	text := strings.Join(b.Args, ", ")
	ctx.PDF.SetFontStyle("I")
	ctx.PDF.Write(5, text)
	ctx.PDF.SetFontStyle("") // zurück zu normal
	return ""
}

func (c *ItalicCommand) SetArgs(args []string) {
	c.Args = args
}
