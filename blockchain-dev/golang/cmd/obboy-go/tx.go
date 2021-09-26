package main

import (
	"fmt"
	"os"

	"github.com/kalam360/obboy-blockchain/golang/database"
	"github.com/spf13/cobra"
)

func txCmd() *cobra.Command {
	var txsCmd = &cobra.Command{
		Use: "tx",
		Short: "Interact with txs (add ..l)",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return incorrectUsageErr()
		},
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	txsCmd.AddCommand(txAddCmd())
	return txsCmd
}

func txAddCmd() *cobra.Command {
	var txAddcmd = &cobra.Command {
		Use: "add",
		Short: "Adds new tx to Database",
		Run: func(cmd *cobra.Command, args []string) {
			from, _ := cmd.Flags().GetString(flagFrom)
			to, _ := cmd.Flags().GetString(flagTo)
			value, _ := cmd.Flags().GetUint(flagValue)

			fromAcc := database.NewAccount(from)
			toAcc := database.NewAccount(to)

			tx := database.NewTx(fromAcc, toAcc, value, "")

			state, err := database.NewStateFromDisk()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			err = state.Add(tx)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}

			err = state.Persist()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)

			}

			fmt.Println("TX is successfully added to the ledger.")
		},
	} 
}



