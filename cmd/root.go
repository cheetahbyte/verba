package cmd

import (
	"errors"
	"log"

	verba "github.com/cheetahbyte/verba/pkg"
	"github.com/spf13/cobra"
)

var RootCMD = &cobra.Command{
	Use:   "verba [file]",
	Short: "Verba: easiest typeset",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		file := ""
		if len(args) > 0 {
			file = args[0]
		} else {
			return errors.New("provide a file")
		}

		// Flags lesen
		debug, _ := cmd.Flags().GetBool("debug")
		// outPath, _ := cmd.Flags().GetString("out")

		if debug {
			log.Println("Debugmodus aktiviert")
		}

		return verba.Verba(file)
	},
}
