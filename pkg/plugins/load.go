package plugins

import (
	"fmt"
	"plugin"
)

func Load(path string, ctx *PluginContext) error {
	p, err := plugin.Open(path)
	if err != nil {
		return fmt.Errorf("Fehler beim Öffnen des Plugins: %w", err)
	}

	sym, err := p.Lookup("RegisterPlugin")
	if err != nil {
		return fmt.Errorf("Symbol 'RegisterPlugin' nicht gefunden: %w", err)
	}

	registerFunc, ok := sym.(func(*PluginContext))
	if !ok {
		return fmt.Errorf("Symbol 'RegisterPlugin' hat falschen Typ")
	}

	registerFunc(ctx)
	return nil
}
