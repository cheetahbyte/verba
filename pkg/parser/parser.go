package parser

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/cheetahbyte/verba/pkg/commands"
	"github.com/cheetahbyte/verba/pkg/context"
)

var commandRegex = regexp.MustCompile(`::(\w+)\{(.*?)\}`)

func ParseFile(path string, reg context.CommandRegistry) ([]any, error) {
	var result []any
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		cmds, err := ProcessCommand(line, reg)
		if err != nil {
			fmt.Println("Fehler:", err)
			continue
		}
		result = append(result, cmds...)
	}

	return result, nil
}

type ArgSetter interface {
	SetArgs([]string)
}

func setArgs(cmd any, args []string) {
	if s, ok := cmd.(ArgSetter); ok {
		s.SetArgs(args)
	}
}

func ProcessCommand(line string, reg context.CommandRegistry) ([]any, error) {
	line = strings.TrimSpace(line)
	if line == "" {
		return nil, nil
	}

	// ➤ Sonderfall: kompletter Blockbefehl
	if strings.HasPrefix(line, "::") && commandRegex.MatchString(line) {
		matches := commandRegex.FindAllStringSubmatch(line, -1)
		if len(matches) == 1 && matches[0][0] == line {
			cmdName := matches[0][1]
			args := splitArgs(matches[0][2])

			if cmd, ok := reg.Block.Get(cmdName); ok {
				setArgs(cmd, args)
				return []any{cmd}, nil
			}

			return nil, fmt.Errorf("unknown block command: %s", cmdName)
		}
	}

	// ➤ Inline-Kommandos im Fließtext
	var paragraph = &commands.ParagraphCommand{}
	lastIndex := 0
	matches := commandRegex.FindAllStringSubmatchIndex(line, -1)

	for _, match := range matches {
		if match[0] > lastIndex {
			text := line[lastIndex:match[0]]
			paragraph.Elements = append(paragraph.Elements, &commands.TextCommand{Args: []string{text}})
		}

		cmdName := line[match[2]:match[3]]
		args := splitArgs(line[match[4]:match[5]])

		if cmd, ok := reg.Inline.Get(cmdName); ok {
			setArgs(cmd, args)
			paragraph.Elements = append(paragraph.Elements, cmd)
		} else {
			paragraph.Elements = append(paragraph.Elements,
				&commands.TextCommand{Args: []string{"[UNKNOWN:" + cmdName + "]"}})
		}

		lastIndex = match[1]
	}

	if lastIndex < len(line) {
		paragraph.Elements = append(paragraph.Elements,
			&commands.TextCommand{Args: []string{line[lastIndex:]}})
	}

	return []any{paragraph}, nil
}

func splitArgs(argStr string) []string {
	args := strings.Split(argStr, ",")
	for i := range args {
		args[i] = strings.TrimSpace(args[i])
	}
	return args
}
