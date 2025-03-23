package commands

import (
	"github.com/cheetahbyte/verba/pkg/context"
	"github.com/cheetahbyte/verba/pkg/registries"
)

func RegisterAll(br *registries.Registry[context.BlockCommand], ir *registries.Registry[context.InlineCommand]) {
	registerInline(ir)
	registerBlock(br)

}

func registerInline(r *registries.Registry[context.InlineCommand]) {
	r.Register("bold", &BoldCommand{})

	r.Register("italic", &ItalicCommand{})
	r.Register("paragraph", &ParagraphCommand{})

}

func registerBlock(r *registries.Registry[context.BlockCommand]) {
	r.Register("class", &ClassCommand{})
	r.Register("figure", &FigureCommand{})
	r.Register("margin", &MarginCommand{})
	r.Register("section", &SectionCommand{})
	r.Register("subsection", &SubsectionCommand{})
	r.Register("unknown", &UnknownCommand{})
	r.Register("text", &TextCommand{})
	r.Register("paragraph", &ParagraphCommand{})
}
