package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cheetahbyte/verba/pkg/commands"
	"github.com/cheetahbyte/verba/pkg/documents"
	"github.com/cheetahbyte/verba/pkg/parser"
	"github.com/jung-kurt/gofpdf"
)

func main() {
	start := time.Now()
	file, err := os.Open("text.verba")
	if err != nil {
		fmt.Println("Fehler beim Öffnen der Datei:", err)
		return
	}
	defer file.Close()

	documentSettings := documents.Document{
		Margin: documents.DocumentMargin{
			Left:   10.0,
			Right:  10.0,
			Top:    10.0,
			Bottom: 10.0,
		},
		PageWidth: 210.0,
	}

	var commandsList []commands.CommandResult
	commandsList, _ = parser.ParseFile("text.verba")

	pdf := gofpdf.New("P", "mm", "A4", "")
	y := 0.0
	fontLoaded := false
	pdfInitialized := false

	fontRegular := "fonts/cmunrm.ttf"
	fontItalic := "fonts/cmunti.ttf"
	fontBold := "fonts/cmunbx.ttf"
	fontBoldItalic := "fonts/cmunbi.ttf"

	pdf.AddUTF8Font("CMUSerif", "", fontRegular)
	pdf.AddUTF8Font("CMUSerif", "I", fontItalic)
	pdf.AddUTF8Font("CMUSerif", "B", fontBold)
	pdf.AddUTF8Font("CMUSerif", "BI", fontBoldItalic)
	pdf.SetFont("CMUSerif", "", 11)

	for _, cmd := range commandsList {
		if !pdfInitialized {
			if cmd.Type == "class" {
				documentSettings = documents.DocumentClasses[cmd.Args[0]]
			}

			pdf.SetMargins(documentSettings.Margin.Left, documentSettings.Margin.Top, documentSettings.Margin.Right)
			pdf.SetAutoPageBreak(true, documentSettings.Margin.Bottom)
			pdf.AddPage()
			pdf.SetXY(documentSettings.Margin.Left, documentSettings.Margin.Top)
			y = documentSettings.Margin.Top
			pdfInitialized = true

			// Fonts laden
			if !fontLoaded {
				pdf.AddUTF8Font("CMUSerif", "", "fonts/cmunrm.ttf")
				pdf.AddUTF8Font("CMUSerif", "I", "fonts/cmunti.ttf")
				pdf.AddUTF8Font("CMUSerif", "B", "fonts/cmunbx.ttf")
				pdf.AddUTF8Font("CMUSerif", "BI", "fonts/cmunbi.ttf")
				pdf.SetFont("CMUSerif", "", 11)
				fontLoaded = true
			}

			// Bei class sofort return, da er oben schon verarbeitet wurde
			if cmd.Type == "class" {
				continue
			}
		}
		switch cmd.Type {
		case "class":
			documentSettings = documents.DocumentClasses[cmd.Args[0]]
		case "margin":
			if len(cmd.Args) == 4 {
				newLeft, _ := strconv.ParseFloat(cmd.Args[0], 64)
				newRight, _ := strconv.ParseFloat(cmd.Args[1], 64)
				newTop, _ := strconv.ParseFloat(cmd.Args[2], 64)
				newBottom, _ := strconv.ParseFloat(cmd.Args[3], 64)
				documentSettings.Margin.Right = newRight
				documentSettings.Margin.Left = newLeft
				documentSettings.Margin.Top = newTop
				documentSettings.Margin.Bottom = newBottom
				pdf.SetMargins(documentSettings.Margin.Left, documentSettings.Margin.Top, documentSettings.Margin.Right)
				pdf.SetAutoPageBreak(true, documentSettings.Margin.Bottom)
				pdf.SetXY(documentSettings.Margin.Left, documentSettings.Margin.Top)
				y = documentSettings.Margin.Top // Reset text position
			}

		case "subsection":
			if y > documentSettings.Margin.Top {
				y += 4
			}
			pdf.SetFont("CMUSerif", "B", 12)
			pdf.SetXY(documentSettings.Margin.Left, y)
			pdf.MultiCell(documentSettings.TextWidth(), 7, cmd.Args[0], "", "L", false)
			y = pdf.GetY() + 1
			pdf.SetFont("CMUSerif", "", 11)

		case "section":
			if y > documentSettings.Margin.Top {
				y += 4
			}
			pdf.SetFont("CMUSerif", "B", 14)
			pdf.SetXY(documentSettings.Margin.Left, y)
			pdf.MultiCell(documentSettings.TextWidth(), 7, cmd.Args[0], "", "L", false)
			y = pdf.GetY() + 1
			pdf.SetFont("CMUSerif", "", 11)

		case "bold":
			pdf.SetFont("CMUSerif", "B", 11)
			// **✅ Fix: Handle inline bold text correctly**
			text := strings.Join(cmd.Args, " ")
			if strings.HasSuffix(text, "\\") {
				text = strings.TrimSuffix(text, "\\")
				pdf.MultiCell(documentSettings.TextWidth(), 5, text, "", "L", false)
				y = pdf.GetY() + 5 // **New paragraph because of "\\"**
			} else {
				pdf.MultiCell(documentSettings.TextWidth(), 5, text, "", "L", false)
				y = pdf.GetY() // **Inline bold (no extra spacing)**
			}
			pdf.SetFont("CMUSerif", "", 11)

		case "text":
			pdf.SetXY(documentSettings.Margin.Left, y)
			text := strings.Join(cmd.Args, " ")
			if strings.HasSuffix(text, "\\") {
				text = strings.TrimSuffix(text, "\\")
				pdf.MultiCell(documentSettings.TextWidth(), 5, text, "", "J", false)
				y = pdf.GetY() + 5
			} else {
				pdf.MultiCell(documentSettings.TextWidth(), 5, text, "", "J", false)
				y = pdf.GetY()
			}

		case "figure":
			if len(cmd.Args) < 1 {
				break
			}

			imagePath := cmd.Args[0]
			caption := ""
			if len(cmd.Args) >= 2 {
				caption = cmd.Args[1]
			}

			spacingAbove := 8.0
			spacingBelow := 4.0
			spacingImageToCaption := 1.0

			y += spacingAbove // Abstand nach vorherigem Text

			imgWidth := documentSettings.TextWidth()
			imgHeight := 0.0 // Höhe automatisch skalieren

			// Bild zeichnen
			pdf.ImageOptions(imagePath, documentSettings.Margin.Left, y, imgWidth, imgHeight, false, gofpdf.ImageOptions{
				ImageType: "",
				ReadDpi:   true,
			}, 0, "")

			// Manuelle Y-Anpassung nach dem Bild
			y += imgWidth * 0.6 // Falls Bildhöhe nicht fix ist

			if caption != "" {
				y += spacingImageToCaption // Abstand zwischen Bild und Caption
				pdf.SetFont("CMUSerif", "I", 10)
				pdf.SetXY(documentSettings.Margin.Left, y)
				pdf.MultiCell(documentSettings.TextWidth(), 5, caption, "", "C", false)
				y = pdf.GetY() + spacingBelow // Abstand nach Caption
				pdf.SetFont("CMUSerif", "", 11)
			}

		case "unknown":
			pdf.SetXY(documentSettings.Margin.Left, y)
			pdf.MultiCell(documentSettings.TextWidth(), 5, "[Unbekannter Befehl: "+cmd.Args[0]+"]", "", "L", false)
			y = pdf.GetY() + 2
		}

		// stop overlap
		if y > (documentSettings.PageWidth - documentSettings.Margin.Bottom - 10) {
			pdf.AddPage()
			y = documentSettings.Margin.Top
		}

	}

	processTime := time.Since(start)

	err = pdf.OutputFileAndClose("output.pdf")
	if err != nil {
		fmt.Println("Fehler beim Speichern der PDF:", err)
	} else {
		fmt.Println("PDF erfolgreich erstellt: output.pdf")
		fmt.Printf("Verarbeitung abgeschlossen in %s ms\n", processTime)
	}

}
