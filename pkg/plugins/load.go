package plugins

import (
	"fmt"
	"plugin"
	"time"

	log "github.com/sirupsen/logrus"
)

func Load(path string, ctx *PluginContext) error {
	start := time.Now()
	p, err := plugin.Open(path)
	log.WithField("took", time.Since(start)).Debug("loaded plugin ", path)
	if err != nil {
		return fmt.Errorf("failed to open plugin")
	}

	sym, err := p.Lookup("RegisterPlugin")
	if err != nil {
		return fmt.Errorf("function 'RegisterPlugin'")
	}

	registerFunc, ok := sym.(func(*PluginContext))
	if !ok {
		return fmt.Errorf("function 'RegisterPlugin' has wrong type")
	}

	registerFunc(ctx)
	return nil
}
