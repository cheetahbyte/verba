package commands

import (
	"github.com/cheetahbyte/verba/pkg/documents"
	"github.com/jung-kurt/gofpdf"
)

type Command interface {
	Execute(pdf *gofpdf.Fpdf, y *float64, doc *documents.Document) error
}
