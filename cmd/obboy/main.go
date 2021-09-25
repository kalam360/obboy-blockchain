package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var obboyCmd = &cobra.Command {
		Use: "obboy",
		Short: "Obboy Blockchain Cli",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	err := obboyCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}