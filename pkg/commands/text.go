package commands

import (
	"strings"

	"github.com/cheetahbyte/verba/pkg/documents"
	"github.com/jung-kurt/gofpdf"
)

type TextCommand struct {
	Args []string
}

func (t TextCommand) Execute(pdf *gofpdf.Fpdf, y *float64, doc *documents.Document) error {
	text := strings.Join(t.Args, " ")
	pdf.MultiCell(doc.TextWidth(), 5, text, "", "J", false)
	*y = pdf.GetY()

	return nil
}
