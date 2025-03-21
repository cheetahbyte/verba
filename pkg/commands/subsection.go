package commands

func Subsection(args []string) CommandResult {
	return CommandResult{Type: "subsection", Args: args}
}
