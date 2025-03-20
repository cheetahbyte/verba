package commands

import "fmt"

// Margin parses the margin command but does not apply it
func Margin(args []string) (CommandResult, error) {
	if len(args) != 4 {
		return CommandResult{}, fmt.Errorf("Fehler: ::margin{left, right, top, bottom} erfordert 4 Werte")
	}

	return CommandResult{Type: "margin", Args: args}, nil
}
