package commands

import (
	"github.com/cheetahbyte/verba/pkg/citation"
	"github.com/cheetahbyte/verba/pkg/documents"
	"github.com/jung-kurt/gofpdf"
)

type BibfileCommand struct {
	Args []string
}

func (b BibfileCommand) Execute(_ *gofpdf.Fpdf, _ *float64, doc *documents.Document) error {
	if len(b.Args) != 1 {
		return nil
	}
	return citation.LoadBibFile(b.Args[0], doc)
}
