// pkg/commands/figure.go
package commands

import (
	"github.com/cheetahbyte/verba/pkg/context"
	"github.com/jung-kurt/gofpdf"
)

type FigureCommand struct {
	Args []string
}

func (f FigureCommand) Execute(ctx *context.CommandContext) error {
	if len(f.Args) < 1 {
		return nil // Kein Bildpfad übergeben, ignoriere
	}

	imagePath := f.Args[0]
	caption := ""
	if len(f.Args) >= 2 {
		caption = f.Args[1]
	}

	spacingAbove := 8.0
	spacingBelow := 4.0
	spacingImageToCaption := 1.0

	*ctx.Y += spacingAbove

	imgWidth := ctx.Document.TextWidth()
	imgHeight := 0.0 // automatische Höhe

	ctx.PDF.ImageOptions(imagePath, ctx.Document.Margin.Left, *ctx.Y, imgWidth, imgHeight, false, gofpdf.ImageOptions{
		ImageType: "",
		ReadDpi:   true,
	}, 0, "")

	*ctx.Y += imgWidth * 0.6 // einfache heuristische Höhe (kann später verfeinert werden)

	if caption != "" {
		*ctx.Y += spacingImageToCaption
		ctx.PDF.SetFont("CMUSerif", "I", 10)
		ctx.PDF.SetXY(ctx.Document.Margin.Left, *ctx.Y)
		ctx.PDF.MultiCell(ctx.Document.TextWidth(), 5, caption, "", "C", false)
		*ctx.Y = ctx.PDF.GetY() + spacingBelow
		ctx.PDF.SetFont("CMUSerif", "", 11)
	}

	return nil
}

func (c *FigureCommand) SetArgs(args []string) {
	c.Args = args
}
