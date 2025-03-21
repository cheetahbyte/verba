// pkg/commands/section.go
package commands

import (
	"strings"

	"github.com/cheetahbyte/verba/pkg/documents"
	"github.com/jung-kurt/gofpdf"
)

type SectionCommand struct {
	Args []string
}

func (s SectionCommand) Execute(pdf *gofpdf.Fpdf, y *float64, doc *documents.Document) error {
	if *y > doc.Margin.Top {
		*y += 4
	}

	pdf.SetFont("CMUSerif", "B", 14)
	pdf.SetXY(doc.Margin.Left, *y)
	pdf.MultiCell(doc.TextWidth(), 7, strings.Join(s.Args, " "), "", "L", false)
	*y = pdf.GetY() + 1
	pdf.SetFont("CMUSerif", "", 11)

	return nil
}
