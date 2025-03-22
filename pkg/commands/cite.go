package commands

import (
	"fmt"
	"strings"

	"github.com/cheetahbyte/verba/pkg/documents"
	"github.com/jung-kurt/gofpdf"
)

type CiteCommand struct {
	Args []string
}

func (c CiteCommand) InlineText(doc *documents.Document) string {
	labels := []string{}
	fmt.Println("cite command")
	for _, key := range c.Args {
		fmt.Println("citing", key)
		doc.Cite(key)
		id := doc.Bibliography.CitationIDs[key]
		fmt.Println(id)
		labels = append(labels, fmt.Sprintf("[%d]", id))
	}
	return strings.Join(labels, ", ")
}

func (c CiteCommand) Execute(pdf *gofpdf.Fpdf, y *float64, doc *documents.Document) error {
	return nil
}
