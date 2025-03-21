// pkg/parser/parser.go
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

func ParseFile(path string) ([]commands.Command, error) {
	var result []commands.Command
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		cmds, err := ProcessCommand(line)
		if err != nil {
			fmt.Println("Fehler:", err)
			continue
		}
		result = append(result, cmds...)
	}

	return result, nil
}

// in parser.go
func ProcessCommand(line string) ([]commands.Command, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return []commands.Command{commands.TextCommand{Args: []string{""}}}, nil
	}

	matches := commandRegex.FindStringSubmatch(line)
	if len(matches) == 3 {
		command := matches[1]
		args := splitArgs(matches[2])

		switch command {
		case "class":
			return []commands.Command{commands.ClassCommand{Args: args}}, nil
		case "bold":
			return []commands.Command{commands.BoldCommand{Args: args}}, nil
		case "margin":
			return []commands.Command{commands.MarginCommand{Args: args}}, nil
		case "section":
			return []commands.Command{commands.SectionCommand{Args: args}}, nil
		case "subsection":
			return []commands.Command{commands.SubsectionCommand{Args: args}}, nil
		case "figure":
			return []commands.Command{commands.FigureCommand{Args: args}}, nil
		case "include":
			if len(args) != 1 {
				return nil, fmt.Errorf("::include erwartet genau ein Argument")
			}
			return ParseFile(args[0]) // direkt Commands aus Datei zurückgeben
		default:
			return []commands.Command{commands.UnknownCommand{Raw: command}}, nil
		}
	}

	return []commands.Command{commands.TextCommand{Args: []string{line}}}, nil
}

func splitArgs(argStr string) []string {
	args := strings.Split(argStr, ",")
	for i := range args {
		args[i] = strings.TrimSpace(args[i])
	}
	return args
}
