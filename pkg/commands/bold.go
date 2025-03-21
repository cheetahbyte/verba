package commands

import (
	"strings"

	"github.com/cheetahbyte/verba/pkg/documents"
	"github.com/jung-kurt/gofpdf"
)

type BoldCommand struct {
	Args []string
}

func (b BoldCommand) Execute(pdf *gofpdf.Fpdf, y *float64, doc *documents.Document) error {
	pdf.SetFont("CMUSerif", "B", 11)
	text := strings.Join(b.Args, " ")
	if strings.HasSuffix(text, "\\") {
		text = strings.TrimSuffix(text, "\\")
		pdf.MultiCell(doc.TextWidth(), 5, text, "", "L", false)
		*y = pdf.GetY() + 5
	} else {
		pdf.MultiCell(doc.TextWidth(), 5, text, "", "L", false)
		*y = pdf.GetY()
	}
	pdf.SetFont("CMUSerif", "", 11)
	return nil
}
