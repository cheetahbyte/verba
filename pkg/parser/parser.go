package parser

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/cheetahbyte/verba/pkg/commands"
)

var commandRegex = regexp.MustCompile(`^::(\w+)\{(.*?)\}$`)

func ParseFile(path string) ([]commands.CommandResult, error) {
	var result []commands.CommandResult
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		cmd, err := ProcessCommand(line)
		if err != nil {
			fmt.Println("Fehler:", err)
			continue
		}

		if cmd.Type == "include" && len(cmd.Args) == 1 {
			subCommands, err := ParseFile(cmd.Args[0])
			if err != nil {
				fmt.Println("Fehler beim Parsen von Include-Datei:", err)
				continue
			}
			result = append(result, subCommands...)
		} else {
			result = append(result, cmd)
		}
	}

	return result, nil
}

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
		// control instructions
		case "margin":
			return commands.Margin(args)
		case "class":
			return commands.Class(args)
		case "include":
			return commands.Include(args)
		// specials
		case "figure":
			return commands.Figure(args)
		// text formatting
		case "subsection":
			return commands.Subsection(args), nil
		case "section":
			return commands.Section(args), nil
		case "bold":
			return commands.Bold(args), nil
		default:
			return commands.CommandResult{Type: "unknown", Args: []string{command}}, nil
		}
	}

	// If it's not a command, treat it as normal text
	return commands.CommandResult{Type: "text", Args: []string{line}}, nil
}
