package commands

import (
	"fmt"

	"github.com/cheetahbyte/verba/pkg/context"
)

type ParagraphCommand struct {
	Elements []any // Mixed: TextCommand, CiteCommand, BoldCommand, etc.
}

func (p ParagraphCommand) InlineText(ctx *context.CommandContext) string {
	return ""
}

func (p ParagraphCommand) ExecuteInline(ctx *context.CommandContext) error {

	ctx.PDF.SetY(*ctx.Y) // 🟢 Y-Position setzen

	for _, elem := range p.Elements {
		switch v := elem.(type) {
		case *TextCommand:
			ctx.PDF.Write(5, v.Args[0])
		case context.InlineCommand:
			ctx.PDF.Write(5, v.InlineText(ctx))
		case *DeferredInlineCommand:
			ctx.PDF.Write(5, v.ExecuteInline(ctx))
		default:
			fmt.Println("Unbekannter Inline-Typ:", fmt.Sprintf("%T", v))
		}
	}

	*ctx.Y += 6   // 🟢 nächste Y-Zeile vorbereiten
	ctx.PDF.Ln(6) // 🟢 neue Zeile im PDF

	return nil
}

func (p ParagraphCommand) Execute(ctx *context.CommandContext) error {
	return p.ExecuteInline(ctx)
}

func (c *ParagraphCommand) SetArgs(args []string) {
}
