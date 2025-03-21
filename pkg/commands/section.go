package commands

func Section(args []string) CommandResult {
	return CommandResult{Type: "section", Args: args}
}
