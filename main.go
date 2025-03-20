package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/jung-kurt/gofpdf"
)

func main() {
	// Datei öffnen
	file, err := os.Open("file.v")
	if err != nil {
		fmt.Println("Fehler beim Öffnen der Datei:", err)
		return
	}
	defer file.Close()

	// PDF erstellen
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetFont("Arial", "", 12)
	pdf.AddPage()

	// Datei zeilenweise lesen und in PDF schreiben
	y := 10.0 // Startposition für den Text
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		pdf.Text(10, y, scanner.Text())
		y += 10 // Abstand zur nächsten Zeile
	}

	// PDF speichern
	err = pdf.OutputFileAndClose("output.pdf")
	if err != nil {
		fmt.Println("Fehler beim Speichern der PDF:", err)
	}
}
