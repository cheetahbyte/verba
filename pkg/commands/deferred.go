package commands

import (
	"fmt"

	"github.com/cheetahbyte/verba/pkg/context"
)

// this command will never be used by the user.
type DeferredCommand struct {
	CommandName string
	Args        []string
}

func (c *DeferredCommand) Execute(ctx *context.CommandContext) error {
	if realCmd, ok := ctx.CmdRegistry.Block.Get(c.CommandName); ok {
		if setter, ok := realCmd.(context.ArgSetter); ok {
			setter.SetArgs(c.Args)
		}
		return realCmd.Execute(ctx)
	}

	return fmt.Errorf("command %q not found at execution time", c.CommandName)
}
