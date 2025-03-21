// pkg/commands/class.go
package commands

import (
	"fmt"

	"github.com/cheetahbyte/verba/pkg/documents"
	"github.com/jung-kurt/gofpdf"
)

type ClassCommand struct {
	Args []string
}

func (c ClassCommand) Execute(pdf *gofpdf.Fpdf, y *float64, doc *documents.Document) error {
	if len(c.Args) != 1 {
		return fmt.Errorf("::class erwartet genau ein Argument")
	}

	className := c.Args[0]
	class, ok := documents.DocumentClasses[className]
	if !ok {
		return fmt.Errorf("Unbekannte Dokumentklasse: %s", className)
	}

	// apply class settings
	*doc = class

	pdf.SetMargins(doc.Margin.Left, doc.Margin.Top, doc.Margin.Right)
	pdf.SetAutoPageBreak(true, doc.Margin.Bottom)
	pdf.SetXY(doc.Margin.Left, doc.Margin.Top)
	*y = doc.Margin.Top

	return nil
}
