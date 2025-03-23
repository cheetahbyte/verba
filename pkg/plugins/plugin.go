package plugins

import "github.com/cheetahbyte/verba/pkg/context"

type Plugin interface {
	Register(ctx *PluginContext)
}

type PluginContext struct {
	Commands    *context.CommandRegistry
	Environment map[string]any
}
