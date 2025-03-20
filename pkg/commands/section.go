package commands

// Subsection processes a subsection command
func Subsection(args []string) CommandResult {
	return CommandResult{Type: "subsection", Args: args}
}
