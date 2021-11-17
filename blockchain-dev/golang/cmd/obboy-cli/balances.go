// list and check the balances of accounts

package main

import (
	"fmt"
	"os"

	"github.com/kalam360/obboy-blockchain-golang/database"
	"github.com/kalam360/obboy-blockchain-golang/node"
	"github.com/spf13/cobra"
)

func balancesCmd() *cobra.Command {
	var balancesCmd = &cobra.Command{
		Use:   "balances",
		Short: "Interacts with the balances (list ...).",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return incorrectUsageErr()
		},
		Run: func(cmd *cobra.Command, args []string) {
		},
	}

	balancesCmd.AddCommand(balancesListCmd())

	return balancesCmd
}

func balancesListCmd() *cobra.Command {
	var balancesListCmd = &cobra.Command{
		Use: "list",
		Short: "Lists all balances.",
		Run: func(cmd *cobra.Command, args []string) {
			state, err := database.NewStateFromDisk(getDataDirFromCmd(cmd), node.DefaultMiningDifficulty)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			defer state.Close()

			fmt.Printf("Accounts Balances at %x: \n", state.LatestBlockHash())
			fmt.Println("-------------------")
			fmt.Println("")
			for account, balance := range state.Balances {
				fmt.Println(fmt.Sprintf("%s: %d", account.String(), balance))
			}

			fmt.Println("")
			fmt.Printf("Accounts Nonces:")
			fmt.Println("") 
			fmt.Println("----------------")
			fmt.Println("")
			for account, nonce := range state.Account2Nonce {
				fmt.Println(fmt.Sprintf("%s: %d", account.String(), nonce))
			}
			
		},
	}

	addDefaultRequiredFlags(balancesListCmd)

	return balancesListCmd
}
