package main

import (
	"github.com/cheetahbyte/verba/cmd"
)

func main() {
	cmd.RootCMD.Flags().Bool("debug", false, "Debugmodus aktivieren")
	cmd.RootCMD.Flags().String("out", "out.pdf", "Pfad zur Ausgabedatei")
	cmd.RootCMD.Flags().String("bibfile", "", "Pfad zur .bib-Datei")
	cmd.RootCMD.Execute()
}
