package commands

import (
	"strings"

	"github.com/cheetahbyte/verba/pkg/documents"
	"github.com/jung-kurt/gofpdf"
)

type CiteCommand struct {
	Args []string
}

func (c CiteCommand) InlineText(doc *documents.Document) string {
	labels := []string{}
	for _, key := range c.Args {
		entry := doc.Cite(key) // registriert das Zitat
		label := doc.Bibliography.Style.CiteLabel(entry)
		labels = append(labels, label)
	}
	return strings.Join(labels, ", ")
}

func (c CiteCommand) Execute(pdf *gofpdf.Fpdf, y *float64, doc *documents.Document) error {
	return nil
}
