// pkg/commands/subsection.go
package commands

import (
	"github.com/cheetahbyte/verba/pkg/documents"
	"github.com/jung-kurt/gofpdf"
)

type SubsectionCommand struct {
	Args []string
}

func (s SubsectionCommand) Execute(pdf *gofpdf.Fpdf, y *float64, doc *documents.Document) error {
	pdf.SetFont("CMUSerif", "B", 12)
	pdf.CellFormat(0, 10, s.Args[0], "", 1, "L", false, 0, "")
	pdf.Ln(4)
	pdf.SetFont("CMUSerif", "", 11)
	return nil
}
