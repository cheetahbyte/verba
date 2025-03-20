package commands

// Bold processes a bold text command
func Bold(args []string) CommandResult {
	return CommandResult{Type: "bold", Args: args}
}
