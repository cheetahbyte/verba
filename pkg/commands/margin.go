package commands

import (
	"fmt"
	"strconv"

	"github.com/cheetahbyte/verba/pkg/documents"
	"github.com/jung-kurt/gofpdf"
)

type MarginCommand struct {
	Args []string
}

func (m MarginCommand) Execute(pdf *gofpdf.Fpdf, y *float64, doc *documents.Document) error {
	if len(m.Args) != 4 {
		return fmt.Errorf("Fehler: ::margin{left, right, top, bottom} erfordert 4 Werte")
	}

	left, err1 := strconv.ParseFloat(m.Args[0], 64)
	right, err2 := strconv.ParseFloat(m.Args[1], 64)
	top, err3 := strconv.ParseFloat(m.Args[2], 64)
	bottom, err4 := strconv.ParseFloat(m.Args[3], 64)

	if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
		return fmt.Errorf("Fehler: Ungültige Werte für ::margin – erwartet wurden Zahlen")
	}

	doc.Margin.Left = left
	doc.Margin.Right = right
	doc.Margin.Top = top
	doc.Margin.Bottom = bottom

	pdf.SetMargins(left, top, right)
	pdf.SetAutoPageBreak(true, bottom)
	pdf.SetXY(left, top)

	*y = top

	return nil
}
