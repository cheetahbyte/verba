package commands

import (
	"strings"

	"github.com/cheetahbyte/verba/pkg/documents"
	"github.com/jung-kurt/gofpdf"
)

type ItalicCommand struct {
	Args []string
}

func (b ItalicCommand) Execute(pdf *gofpdf.Fpdf, y *float64, doc *documents.Document) error {
	pdf.SetFont("CMUSerif", "I", 11)
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

func (b ItalicCommand) ExecuteInline(doc *documents.Document, pdf *gofpdf.Fpdf) {
	text := strings.Join(b.Args, ", ")
	pdf.SetFontStyle("I")
	pdf.Write(5, text)
	pdf.SetFontStyle("") // zurück zu normal
}
