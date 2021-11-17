// wallet commands

package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/kalam360/obboy-blockchain-golang/wallet"
	"github.com/spf13/cobra"
)

func walletCmd() *cobra.Command {
	var walletCmd = &cobra.Command{
		Use:   "wallet",
		Short: "Manages blockchain accounts and keys",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			return incorrectUsageErr()
		},
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	walletCmd.AddCommand(walletNewAccountCmd())
	walletCmd.AddCommand(walletPrintPrivateKeyCmd())

	return walletCmd

}

func walletNewAccountCmd() *cobra.Command {
	var newAccountCmd = &cobra.Command{
		Use:   "new-account",
		Short: "Creates a new account with a new set of a elliptic-curve Private and Public Keys.",
		Run: func(cmd *cobra.Command, args []string) {
			password := getPassPhrase("Please enter a password to encrypt the new wallet:", true)
			dataDir := getDataDirFromCmd(cmd)

			acc, err := wallet.NewKeyStoreAccount(dataDir, password)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			fmt.Printf("New account created: %s \n", acc.Hex())
			fmt.Printf("Saved in: %s\n", wallet.GetKeystoreDirPath(dataDir))

		},
	}

	addDefaultRequiredFlags(newAccountCmd)
	return newAccountCmd
}

func walletPrintPrivateKeyCmd() *cobra.Command {
	var pkPrintCmd = &cobra.Command{
		Use:   "pk-print",
		Short: "Unlocks keystore file and prints the private and public keys.",
		Run: func(cmd *cobra.Command, args []string) {
			ksFile, _ := cmd.Flags().GetString(flagKeystoreFile)
			password := getPassPhrase("Please enter a password to decrypt the wallet:", false)

			keyJson, err := ioutil.ReadFile(ksFile)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			key, err := keystore.DecryptKey(keyJson, password)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			spew.Dump(key)

		},
	}

	addKeystoreFlag(pkPrintCmd)

	return pkPrintCmd
}

func getPassPhrase(prompt string, confirmation bool) string {
	return utils.GetPassPhrase(prompt, confirmation)
}
