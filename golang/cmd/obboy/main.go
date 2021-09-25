package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/kalam360/obboy-blockchain/database"
)

func main() {
	// create the cobra command 
	var obboyCmd = &cobra.Command {
		Use: "obboy",
		Short: "Obboy Blockchain Cli",
		// it will run the following function when called
		Run: func(cmd *cobra.Command, args []string) {
			database.ReadGenesis()
		},
	}

	// This will execute the command in the terminal and catch any error and print it. 
	err := obboyCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}