package main

import (
	"fmt"
	"time"

	"github.com/cheetahbyte/verba/pkg/documents"
	"github.com/cheetahbyte/verba/pkg/parser"
	"github.com/jung-kurt/gofpdf"
)

func main() {
	start := time.Now()

	documentSettings := documents.Document{
		Margin: documents.DocumentMargin{
			Left:   10.0,
			Right:  10.0,
			Top:    10.0,
			Bottom: 10.0,
		},
		PageWidth: 210.0,
	}

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

	pdf.SetMargins(documentSettings.Margin.Left, documentSettings.Margin.Top, documentSettings.Margin.Right)
	pdf.SetAutoPageBreak(true, documentSettings.Margin.Bottom)
	pdf.AddPage()
	pdf.SetXY(documentSettings.Margin.Left, documentSettings.Margin.Top)
	y = documentSettings.Margin.Top

	// Datei parsen
	commandList, err := parser.ParseFile("thesis.verba")
	if err != nil {
		fmt.Println("Fehler beim Parsen:", err)
		return
	}

	// Alle Commands ausführen
	for _, cmd := range commandList {
		err := cmd.Execute(pdf, &y, &documentSettings)
		if err != nil {
			fmt.Println("Fehler beim Ausführen eines Befehls:", err)
			continue
		}

		// Seitenüberlauf prüfen
		if y > (documentSettings.PageWidth - documentSettings.Margin.Bottom - 10) {
			pdf.AddPage()
			y = documentSettings.Margin.Top
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
