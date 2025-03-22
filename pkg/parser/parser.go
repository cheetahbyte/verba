package parser

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/cheetahbyte/verba/pkg/commands"
)

var commandRegex = regexp.MustCompile(`::(\w+)\{(.*?)\}`)

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

func ProcessCommand(line string) ([]commands.Command, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return []commands.Command{}, nil
	}

	// ➤ Sonderfall: Zeile ist ein reiner Blockbefehl wie ::section{...}
	if strings.HasPrefix(line, "::") && commandRegex.MatchString(line) {
		matches := commandRegex.FindAllStringSubmatch(line, -1)
		if len(matches) == 1 && matches[0][0] == line {
			command := matches[0][1]
			args := splitArgs(matches[0][2])

			switch command {
			case "class":
				return []commands.Command{commands.ClassCommand{Args: args}}, nil
			case "margin":
				return []commands.Command{commands.MarginCommand{Args: args}}, nil
			case "section":
				return []commands.Command{commands.SectionCommand{Args: args}}, nil
			case "subsection":
				return []commands.Command{commands.SubsectionCommand{Args: args}}, nil
			case "bibfile":
				return []commands.Command{commands.BibfileCommand{Args: args}}, nil
			case "bibliography":
				return []commands.Command{commands.BibliographyCommand{Args: args}}, nil
			case "bibstyle":
				return []commands.Command{commands.BibstyleCommand{Args: args}}, nil
			}

		}
	}

	var paragraph commands.ParagraphCommand
	lastIndex := 0
	matches := commandRegex.FindAllStringSubmatchIndex(line, -1)

	for _, match := range matches {
		if match[0] > lastIndex {
			paragraph.Elements = append(paragraph.Elements,
				commands.TextCommand{Args: []string{line[lastIndex:match[0]]}})
		}

		cmdName := line[match[2]:match[3]]
		args := splitArgs(line[match[4]:match[5]])
		switch cmdName {
		case "cite":
			paragraph.Elements = append(paragraph.Elements, commands.CiteCommand{Args: args})
		case "bold":
			paragraph.Elements = append(paragraph.Elements, commands.BoldCommand{Args: args})
		case "italic":
			paragraph.Elements = append(paragraph.Elements, commands.ItalicCommand{Args: args})

		default:
			paragraph.Elements = append(paragraph.Elements,
				commands.TextCommand{Args: []string{"[UNKNOWN:" + cmdName + "]"}})
		}

		lastIndex = match[1]
	}

	if lastIndex < len(line) {
		paragraph.Elements = append(paragraph.Elements,
			commands.TextCommand{Args: []string{line[lastIndex:]}})
	}

	return []commands.Command{paragraph}, nil
}

func splitArgs(argStr string) []string {
	args := strings.Split(argStr, ",")
	for i := range args {
		args[i] = strings.TrimSpace(args[i])
	}
	return args
}
