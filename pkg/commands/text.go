// pkg/commands/text.go
package commands

import (
	"strings"

	"github.com/cheetahbyte/verba/pkg/documents"
	"github.com/jung-kurt/gofpdf"
)

// this is not a real command used by the user. this is only used internally
type TextCommand struct {
	Args []string
}

func (t TextCommand) Execute(pdf *gofpdf.Fpdf, y *float64, doc *documents.Document) error {
	text := strings.Join(t.Args, " ")

	pdf.SetXY(doc.Margin.Left, *y)

	if strings.HasSuffix(text, "\\") {
		text = strings.TrimSuffix(text, "\\")
		pdf.MultiCell(doc.TextWidth(), 5, text, "", "J", false)
		*y = pdf.GetY() + 5
	} else {
		pdf.MultiCell(doc.TextWidth(), 5, text, "", "J", false)
		*y = pdf.GetY()
	}

	return nil
}
