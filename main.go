package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cheetahbyte/verba/pkg/commands"
	"github.com/cheetahbyte/verba/pkg/context"
	"github.com/cheetahbyte/verba/pkg/documents"
	"github.com/cheetahbyte/verba/pkg/parser"
	"github.com/cheetahbyte/verba/pkg/plugins"
	"github.com/cheetahbyte/verba/pkg/registries"
	"github.com/jung-kurt/gofpdf"
)

func main() {
	start := time.Now()

	blockCmdRegistry := registries.NewRegistry[context.BlockCommand]()
	inlineCmdRegistry := registries.NewRegistry[context.InlineCommand]()
	cmdReg := context.CommandRegistry{
		Inline: inlineCmdRegistry,
		Block:  blockCmdRegistry,
	}

	commands.RegisterAll(cmdReg.Block, cmdReg.Inline)

	doc := documents.NewDocument()
	// Fonts und PDF vorbereiten
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(doc.Margin.Left, doc.Margin.Top, doc.Margin.Right)
	pdf.SetAutoPageBreak(true, doc.Margin.Bottom)
	pdf.AddPage()
	pdf.SetXY(doc.Margin.Left, doc.Margin.Top)
	fontRegular := "fonts/cmunrm.ttf"
	fontItalic := "fonts/cmunti.ttf"
	fontBold := "fonts/cmunbx.ttf"
	fontBoldItalic := "fonts/cmunbi.ttf"

	pdf.AddUTF8Font("CMUSerif", "", fontRegular)
	pdf.AddUTF8Font("CMUSerif", "I", fontItalic)
	pdf.AddUTF8Font("CMUSerif", "B", fontBold)
	pdf.AddUTF8Font("CMUSerif", "BI", fontBoldItalic)
	pdf.SetFont("CMUSerif", "", 11)

	y := doc.Margin.Top

	ctx := &context.CommandContext{
		PDF:         pdf,
		Y:           &y,
		Document:    doc,
		CmdRegistry: cmdReg,
	}

	pluginCtx := &plugins.PluginContext{
		Commands: &context.CommandRegistry{
			Block:  blockCmdRegistry,
			Inline: inlineCmdRegistry,
		},
	}

	files, err := os.ReadDir("plugins")
	if err != nil {
		log.Fatalf("Plugins konnten nicht geladen werden: %v", err)
	}

	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".so") {
			path := filepath.Join("plugins", file.Name())
			log.Printf("Lade Plugin: %s\n", path)
			err := plugins.Load(path, pluginCtx)
			if err != nil {
				log.Printf("Fehler beim Laden von Plugin %s: %v", file.Name(), err)
			}
		}
	}

	// Datei parsen
	commandList, err := parser.ParseFile("thesis.verba", ctx.CmdRegistry)
	if err != nil {
		log.Fatalln("Fehler beim Parsen:", err)
	}

	fmt.Printf("Gelesene Kommandos: %d\n", len(commandList))
	for i, cmd := range commandList {
		fmt.Printf("[%d] -> %T\n", i, cmd)
	}

	// Kommandos ausführen
	for _, cmd := range commandList {
		switch c := cmd.(type) {
		case *commands.ParagraphCommand:
			if err := c.ExecuteInline(ctx); err != nil {
				log.Println("Paragraph render error:", err)
			}
		case context.BlockCommand:
			if err := c.Execute(ctx); err != nil {
				log.Println("Block render error:", err)
			}
		default:
			log.Println("Unbekannter Command-Typ:", c)
		}
	}

	// PDF speichern
	if err := pdf.OutputFileAndClose("output.pdf"); err != nil {
		log.Fatalln("Fehler beim Speichern der PDF:", err)
	}

	fmt.Println("PDF erfolgreich erstellt: output.pdf")
	fmt.Printf("Verarbeitung abgeschlossen in %s\n", time.Since(start))
}
