package commands

func Include(args []string) (CommandResult, error) {
	return CommandResult{Type: "include", Args: args}, nil
}
