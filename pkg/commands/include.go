package commands

import (
	"path/filepath"

	"github.com/cheetahbyte/verba/pkg/context"
	"github.com/cheetahbyte/verba/pkg/plugins"
	log "github.com/sirupsen/logrus"
)

type IncludeCommand struct {
	PluginName string
}

func (c *IncludeCommand) Execute(ctx *context.CommandContext) error {
	if _, ok := ctx.Environment["plugins"]; !ok {
		ctx.Environment["plugins"] = []string{c.PluginName}
	} else {
		plugins, _ := ctx.Environment["plugins"].([]string)
		ctx.Environment["plugins"] = append(plugins, c.PluginName)
	}
	err := plugins.Load(filepath.Join("plugins", c.PluginName+".so"), ctx.Environment["pluginCtx"].(*plugins.PluginContext))
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func (c *IncludeCommand) SetArgs(args []string) {
	c.PluginName = args[0]
}
