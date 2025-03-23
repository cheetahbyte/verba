package commands

import (
	"fmt"

	"github.com/cheetahbyte/verba/pkg/context"
)

type DeferredInlineCommand struct {
	CommandName string
	Args        []string
}

func (c *DeferredInlineCommand) ExecuteInline(ctx *context.CommandContext) string {
	if realCmd, ok := ctx.CmdRegistry.Inline.Get(c.CommandName); ok {
		if setter, ok := realCmd.(context.ArgSetter); ok {
			setter.SetArgs(c.Args)
		}
		return realCmd.InlineText(ctx)
	}

	return fmt.Sprintf("inline command %q not found at execution time", c.CommandName)
}
