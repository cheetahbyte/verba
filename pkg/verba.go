package pkg

import (
	"time"

	"github.com/cheetahbyte/verba/pkg/commands"
	"github.com/cheetahbyte/verba/pkg/context"
	"github.com/cheetahbyte/verba/pkg/documents"
	"github.com/cheetahbyte/verba/pkg/parser"
	"github.com/cheetahbyte/verba/pkg/plugins"
	"github.com/cheetahbyte/verba/pkg/registries"
	"github.com/jung-kurt/gofpdf"
	log "github.com/sirupsen/logrus"
)

func NewPDF(doc *documents.Document) *gofpdf.Fpdf {
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
	return pdf
}

func NewContext() *context.CommandContext {
	doc := documents.NewDocument()
	pdf := NewPDF(doc)
	cmdReg := NewCommandRegistry()
	return &context.CommandContext{
		Document:    doc,
		PDF:         pdf,
		Y:           &doc.Margin.Top,
		Environment: make(map[string]any),
		CmdRegistry: cmdReg,
	}
}

func NewPluginContext(ctx *context.CommandContext) *plugins.PluginContext {
	return &plugins.PluginContext{
		Commands: &context.CommandRegistry{
			Block:  ctx.CmdRegistry.Block,
			Inline: ctx.CmdRegistry.Inline,
		},
	}
}

func NewCommandRegistry() *context.CommandRegistry {
	blockCmdRegistry := registries.NewRegistry[context.BlockCommand]()
	inlineCmdRegistry := registries.NewRegistry[context.InlineCommand]()
	return &context.CommandRegistry{
		Inline: inlineCmdRegistry,
		Block:  blockCmdRegistry,
	}
}

func Verba(file string) error {
	start := time.Now()
	ctx := NewContext()
	commands.RegisterAll(ctx.CmdRegistry.Block, ctx.CmdRegistry.Inline)
	ctx.Environment["pluginCtx"] = NewPluginContext(ctx)
	commandList, err := parser.ParseFile(file, ctx)
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
	if err := ctx.PDF.OutputFileAndClose("output.pdf"); err != nil {
		log.Error("cannot save file")
	}

	log.Infof("created output.pdf in %s", time.Since(start))
	return nil
}
