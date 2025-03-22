package commands

import (
	"github.com/cheetahbyte/verba/pkg/documents"
	"github.com/jung-kurt/gofpdf"
)

type ParagraphCommand struct {
	Elements []Command // Mixed: TextCommand, CiteCommand, BoldCommand, etc.
}

func (p ParagraphCommand) ExecuteInline(pdf *gofpdf.Fpdf, y *float64, doc *documents.Document) error {
	for _, elem := range p.Elements {
		switch v := elem.(type) {
		case TextCommand:
			pdf.Write(5, v.Args[0])
		case CiteCommand:
			for _, key := range v.Args {
				doc.Cite(key)
				pdf.Write(5, elem.(CiteCommand).InlineText(doc))
			}
		case ItalicCommand:
			elem.(ItalicCommand).ExecuteInline(doc, pdf)
		case BoldCommand:
			elem.(BoldCommand).ExecuteInline(doc, pdf)
		}

	}
	pdf.Ln(6)
	return nil
}

func (p ParagraphCommand) Execute(pdf *gofpdf.Fpdf, y *float64, doc *documents.Document) error {
	return p.ExecuteInline(pdf, y, doc)
}
