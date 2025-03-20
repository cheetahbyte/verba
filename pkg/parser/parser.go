package parser

import (
	"regexp"
	"strings"

	"github.com/cheetahbyte/verba/pkg/commands"
)

var commandRegex = regexp.MustCompile(`^::(\w+)\{(.*?)\}$`)

// ProcessCommand parses a line and returns structured data
func ProcessCommand(line string) (commands.CommandResult, error) {
	line = strings.TrimSpace(line)

	// Ignore empty lines (no error)
	if line == "" {
		return commands.CommandResult{Type: "text", Args: []string{""}}, nil
	}

	// Check if the line matches the command format
	matches := commandRegex.FindStringSubmatch(line)
	if len(matches) == 3 {
		command := matches[1]
		args := strings.Split(matches[2], ",") // Split arguments by comma
		for i := range args {
			args[i] = strings.TrimSpace(args[i]) // Trim spaces from arguments
		}

		switch command {
		case "margin":
			return commands.Margin(args)
		case "subsection":
			return commands.Subsection(args), nil
		case "bold":
			return commands.Bold(args), nil
		default:
			return commands.CommandResult{Type: "unknown", Args: []string{command}}, nil
		}
	}

	// If it's not a command, treat it as normal text
	return commands.CommandResult{Type: "text", Args: []string{line}}, nil
}
