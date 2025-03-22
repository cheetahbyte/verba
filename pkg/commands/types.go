package commands

import (
	"github.com/cheetahbyte/verba/pkg/documents"
	"github.com/jung-kurt/gofpdf"
)

type CommandResult struct {
	Type string
	Args []string
}

type InlineCommand interface {
	InlineText(doc *documents.Document) string
}

type Command interface {
	Execute(pdf *gofpdf.Fpdf, y *float64, doc *documents.Document) error
}
