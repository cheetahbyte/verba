// pkg/commands/figure.go
package commands

import (
	"github.com/cheetahbyte/verba/pkg/documents"
	"github.com/jung-kurt/gofpdf"
)

type FigureCommand struct {
	Args []string
}

func (f FigureCommand) Execute(pdf *gofpdf.Fpdf, y *float64, doc *documents.Document) error {
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

	*y += spacingAbove

	imgWidth := doc.TextWidth()
	imgHeight := 0.0 // automatische Höhe

	pdf.ImageOptions(imagePath, doc.Margin.Left, *y, imgWidth, imgHeight, false, gofpdf.ImageOptions{
		ImageType: "",
		ReadDpi:   true,
	}, 0, "")

	*y += imgWidth * 0.6 // einfache heuristische Höhe (kann später verfeinert werden)

	if caption != "" {
		*y += spacingImageToCaption
		pdf.SetFont("CMUSerif", "I", 10)
		pdf.SetXY(doc.Margin.Left, *y)
		pdf.MultiCell(doc.TextWidth(), 5, caption, "", "C", false)
		*y = pdf.GetY() + spacingBelow
		pdf.SetFont("CMUSerif", "", 11)
	}

	return nil
}
