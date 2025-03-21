package commands

import (
	"fmt"

	"github.com/cheetahbyte/verba/pkg/documents"
)

func Class(args []string) (CommandResult, error) {
	// first sanity check wether document class exists
	_, ok := documents.DocumentClasses[args[0]]
	if !ok {
		return CommandResult{}, fmt.Errorf("unknown document type: %s", args[0])
	}
	return CommandResult{Type: "class", Args: args}, nil
}
