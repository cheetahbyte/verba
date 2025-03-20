package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cheetahbyte/verba/pkg/commands"
	"github.com/cheetahbyte/verba/pkg/parser"
	"github.com/jung-kurt/gofpdf"
)

func main() {
	start := time.Now()
	// Datei öffnen
	file, err := os.Open("text.verba")
	if err != nil {
		fmt.Println("Fehler beim Öffnen der Datei:", err)
		return
	}
	defer file.Close()

	// Standardmargen setzen
	left, right, top, bottom := 20.0, 20.0, 30.0, 30.0

	// PDF erstellen
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(left, top, right)
	pdf.SetAutoPageBreak(true, bottom)
	pdf.AddPage()
	pageWidth, _ := pdf.GetPageSize()
	textWidth := pageWidth - left - right // Berechnet verfügbare Textbreite

	// CMU Serif Font einbinden
	fontRegular := "fonts/cmunrm.ttf"
	fontItalic := "fonts/cmunti.ttf"
	fontBold := "fonts/cmunbx.ttf"
	fontBoldItalic := "fonts/cmunbi.ttf"

	pdf.AddUTF8Font("CMUSerif", "", fontRegular)
	pdf.AddUTF8Font("CMUSerif", "I", fontItalic)
	pdf.AddUTF8Font("CMUSerif", "B", fontBold)
	pdf.AddUTF8Font("CMUSerif", "BI", fontBoldItalic)
	pdf.SetFont("CMUSerif", "", 11)

	var commandsList []commands.CommandResult
	y := top
	scanner := bufio.NewScanner(file)

	// **Lese und speichere Befehle**
	for scanner.Scan() {
		line := scanner.Text()
		cmd, err := parser.ProcessCommand(line)
		if err != nil {
			fmt.Println("Fehler:", err)
			continue
		}
		commandsList = append(commandsList, cmd)
	}

	// **Verarbeite Befehle und Texte**
	for _, cmd := range commandsList {
		switch cmd.Type {
		case "margin":
			if len(cmd.Args) == 4 {
				newLeft, _ := strconv.ParseFloat(cmd.Args[0], 64)
				newRight, _ := strconv.ParseFloat(cmd.Args[1], 64)
				newTop, _ := strconv.ParseFloat(cmd.Args[2], 64)
				newBottom, _ := strconv.ParseFloat(cmd.Args[3], 64)
				left, right, top, bottom = newLeft, newRight, newTop, newBottom
				textWidth = pageWidth - left - right
				pdf.SetMargins(left, top, right)
				pdf.SetAutoPageBreak(true, bottom)
				pdf.SetXY(left, top)
				y = top // Reset text position
			}

		case "subsection":
			if y > top {
				y += 6 // Mehr Abstand vor Unterabschnitt
			}
			pdf.SetFont("CMUSerif", "B", 14)
			pdf.SetXY(left, y)
			pdf.MultiCell(textWidth, 7, cmd.Args[0], "", "L", false)
			y = pdf.GetY() + 4 // Update y nach MultiCell

			pdf.SetFont("CMUSerif", "", 11) // Zurück auf Standard

		case "bold":
			pdf.SetFont("CMUSerif", "B", 11)
			// **✅ Fix: Handle inline bold text correctly**
			text := strings.Join(cmd.Args, " ")
			if strings.HasSuffix(text, "\\") {
				text = strings.TrimSuffix(text, "\\")
				pdf.MultiCell(textWidth, 5, text, "", "L", false)
				y = pdf.GetY() + 5 // **New paragraph because of "\\"**
			} else {
				pdf.MultiCell(textWidth, 5, text, "", "L", false)
				y = pdf.GetY() // **Inline bold (no extra spacing)**
			}
			pdf.SetFont("CMUSerif", "", 11)

		case "text":
			pdf.SetXY(left, y)
			text := strings.Join(cmd.Args, " ")
			if strings.HasSuffix(text, "\\") {
				text = strings.TrimSuffix(text, "\\")
				pdf.MultiCell(textWidth, 5, text, "", "J", false)
				y = pdf.GetY() + 5
			} else {
				pdf.MultiCell(textWidth, 5, text, "", "J", false)
				y = pdf.GetY()
			}

		case "unknown":
			pdf.SetXY(left, y)
			pdf.MultiCell(textWidth, 5, "[Unbekannter Befehl: "+cmd.Args[0]+"]", "", "L", false)
			y = pdf.GetY() + 2
		}

		// stop overlap
		if y > (pageWidth - bottom - 10) {
			pdf.AddPage()
			y = top
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
