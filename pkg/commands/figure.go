package commands

import (
	"fmt"
)

func Figure(args []string) (CommandResult, error) {
	if len(args) != 2 {
		return CommandResult{}, fmt.Errorf("::figure expects two values. ::figure{<path>, <caption>}")
	}
	return CommandResult{Type: "figure", Args: args}, nil
}
