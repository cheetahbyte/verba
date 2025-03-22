package commands

import (
	"fmt"
	"strings"

	"github.com/cheetahbyte/verba/pkg/bibstyle"
	"github.com/cheetahbyte/verba/pkg/documents"
	"github.com/jung-kurt/gofpdf"
)

type BibstyleCommand struct {
	Args []string
}

func (b BibstyleCommand) Execute(pdf *gofpdf.Fpdf, y *float64, doc *documents.Document) error {
	if len(b.Args) == 0 {
		return fmt.Errorf("::bibstyle expects a style name (e.g. apa, ieee)")
	}

	style := strings.ToLower(b.Args[0])

	switch style {
	case "apa":
		doc.Bibliography.Style = bibstyle.APAFormatter{}
	case "ieee":
		doc.Bibliography.Style = bibstyle.IEEEFormatter{}
	default:
		return fmt.Errorf("unknown bibliography style: %s", style)
	}

	return nil
}
