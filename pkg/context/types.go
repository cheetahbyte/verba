package context

import (
	"github.com/cheetahbyte/verba/pkg/documents"
	"github.com/cheetahbyte/verba/pkg/registries"
	"github.com/jung-kurt/gofpdf"
)

type InlineCommand interface {
	InlineText(ctx *CommandContext) string
	SetArgs([]string)
}

type BlockCommand interface {
	Execute(ctx *CommandContext) error
	SetArgs([]string)
}

type CommandContext struct {
	PDF         *gofpdf.Fpdf
	Y           *float64
	Document    *documents.Document
	CmdRegistry CommandRegistry
	Environment map[string]any
}

type CommandRegistry struct {
	Inline *registries.Registry[InlineCommand]
	Block  *registries.Registry[BlockCommand]
}

type ArgSetter interface {
	SetArgs([]string)
}
