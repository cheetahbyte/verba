// pkg/commands/unknown.go
package commands

import (
	"fmt"

	"github.com/cheetahbyte/verba/pkg/documents"
	"github.com/jung-kurt/gofpdf"
)

type UnknownCommand struct {
	Raw string
}

func (u UnknownCommand) Execute(pdf *gofpdf.Fpdf, y *float64, doc *documents.Document) error {
	text := fmt.Sprintf("[Unbekannter Befehl: %s]", u.Raw)

	pdf.SetXY(doc.Margin.Left, *y)
	pdf.MultiCell(doc.TextWidth(), 5, text, "", "L", false)
	*y = pdf.GetY() + 2

	return nil
}
