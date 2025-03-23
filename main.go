package main

import (
	"flag"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/cheetahbyte/verba/pkg/commands"
	"github.com/cheetahbyte/verba/pkg/context"
	"github.com/cheetahbyte/verba/pkg/documents"
	"github.com/cheetahbyte/verba/pkg/parser"
	"github.com/cheetahbyte/verba/pkg/plugins"
	"github.com/cheetahbyte/verba/pkg/registries"
	"github.com/jung-kurt/gofpdf"
)

var debug = flag.Bool("debug", false, "enable debug logging")

func main() {
	start := time.Now()
	flag.Parse()
	if *debug {
		log.SetLevel(log.DebugLevel)
	}
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:    true,
		DisableColors:    false,
		QuoteEmptyFields: true,
		DisableSorting:   true,
		// Wenn du möchtest:
		DisableLevelTruncation: true,
	})

	log.SetOutput(os.Stdout)

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
		Environment: make(map[string]any),
	}

	ctx.Environment["debug"] = debug

	pluginCtx := &plugins.PluginContext{
		Commands: &context.CommandRegistry{
			Block:  blockCmdRegistry,
			Inline: inlineCmdRegistry,
		},
	}

	ctx.Environment["pluginCtx"] = pluginCtx

	commandList, err := parser.ParseFile("thesis.verba", ctx.CmdRegistry)
	if err != nil {
		log.WithField("file", err.Error()).Error("failed to parse")
	}

	for _, cmd := range commandList {
		switch c := cmd.(type) {
		case *commands.ParagraphCommand:
			if err := c.ExecuteInline(ctx); err != nil {
				log.Error("Fehler beim Ausführen des ParagraphCommand")
			}
		case *commands.DeferredCommand:
			if err := c.Execute(ctx); err != nil {
				log.Errorf("Fehler beim Ausführen von DeferredCommand %q: %v", c.CommandName, err)
			}
		case *commands.DeferredInlineCommand:
			c.ExecuteInline(ctx)
		case context.BlockCommand:
			if err := c.Execute(ctx); err != nil {
				log.Error("Fehler beim Ausführen des BlockCommand")
			}
		default:
			log.Warnf("Unbekannter Command-Typ: %T", cmd)
		}
	}

	// PDF speichern
	if err := pdf.OutputFileAndClose("output.pdf"); err != nil {
		log.Error("cannot save file")
	}

	log.Infof("created output.pdf in %s", time.Since(start))
}
