package commands

import (
	"fmt"

	"github.com/cheetahbyte/verba/pkg/documents"
	"github.com/jung-kurt/gofpdf"
)

type BibliographyCommand struct {
	Args []string
}

func (b BibliographyCommand) Execute(pdf *gofpdf.Fpdf, y *float64, doc *documents.Document) error {
	if len(doc.Bibliography.Used) == 0 {
		fmt.Println("no entries")
		return nil
	}

	pdf.SetFont("CMUSerif", "B", 12)
	pdf.CellFormat(0, 7, "Literatur", "", 1, "L", false, 0, "")
	*y = pdf.GetY() + 2

	pdf.SetFont("CMUSerif", "", 11)
	for i, entry := range doc.Bibliography.Used {
		ref := fmt.Sprintf("[%d] %s", i+1, entry.FormatReference())
		pdf.MultiCell(doc.TextWidth(), 5, ref, "", "J", false)
		*y = pdf.GetY() + 3
	}

	return nil
}
