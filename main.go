package main

import (
	"fmt"
	"log"
	"time"

	"github.com/cheetahbyte/verba/pkg/commands"
	"github.com/cheetahbyte/verba/pkg/documents"
	"github.com/cheetahbyte/verba/pkg/parser"
	"github.com/jung-kurt/gofpdf"
)

func main() {
	start := time.Now()

	doc := documents.NewDocument()
	// PDF initialisieren
	pdf := gofpdf.New("P", "mm", "A4", "")
	y := 0.0

	// Fonts laden
	fontRegular := "fonts/cmunrm.ttf"
	fontItalic := "fonts/cmunti.ttf"
	fontBold := "fonts/cmunbx.ttf"
	fontBoldItalic := "fonts/cmunbi.ttf"

	pdf.AddUTF8Font("CMUSerif", "", fontRegular)
	pdf.AddUTF8Font("CMUSerif", "I", fontItalic)
	pdf.AddUTF8Font("CMUSerif", "B", fontBold)
	pdf.AddUTF8Font("CMUSerif", "BI", fontBoldItalic)
	pdf.SetFont("CMUSerif", "", 11)

	pdf.SetMargins(doc.Margin.Left, doc.Margin.Top, doc.Margin.Right)
	pdf.SetAutoPageBreak(true, doc.Margin.Bottom)
	pdf.AddPage()
	pdf.SetXY(doc.Margin.Left, doc.Margin.Top)
	y = doc.Margin.Top

	// Datei parsen
	commandList, err := parser.ParseFile("thesis.verba")
	if err != nil {
		fmt.Println("Fehler beim Parsen:", err)
		return
	}

	// Alle Commands ausführen
	for _, cmd := range commandList {
		switch c := cmd.(type) {
		case commands.ParagraphCommand:
			err := c.ExecuteInline(pdf, &y, doc)
			if err != nil {
				log.Println("Paragraph render error:", err)
			}
		default:
			c.Execute(pdf, &y, doc)
		}
	}

	// PDF speichern
	err = pdf.OutputFileAndClose("output.pdf")
	if err != nil {
		fmt.Println("Fehler beim Speichern der PDF:", err)
	} else {
		fmt.Println("PDF erfolgreich erstellt: output.pdf")
		fmt.Printf("Verarbeitung abgeschlossen in %s ms\n", time.Since(start))
	}
}
