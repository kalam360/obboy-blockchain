// Copyright 2021 The Obboylabs Ltd

package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/kalam360/obboy-blockchain-golang/fs"
)

const flagDataDir = "datadir"
const flagKeystoreFile = "keystore"

func main() {
	var obboyCmd = &cobra.Command{
		Use:   "obboy-cli",
		Short: "Obboy Blockchain Cli",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("CLI for Obboy Blockchain")
		},
	}

	obboyCmd.AddCommand(versionCmd)
	obboyCmd.AddCommand(balancesCmd())
	obboyCmd.AddCommand(walletCmd())

	err := obboyCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func addDefaultRequiredFlags(cmd *cobra.Command) {
	cmd.Flags().String(flagDataDir, "", "Absolute path to your node's data dir where the DB will be/is stored")
	cmd.MarkFlagRequired(flagDataDir)
}

func addKeystoreFlag(cmd *cobra.Command) {
	cmd.Flags().String(flagKeystoreFile, "", "Absolute path to the encrypted keystore file")
	cmd.MarkFlagRequired(flagKeystoreFile)
}

func getDataDirFromCmd(cmd *cobra.Command) string {
	dataDir, _ := cmd.Flags().GetString(flagDataDir)

	return fs.ExpandPath(dataDir)
}

func incorrectUsageErr() error {
	return fmt.Errorf("incorrect usage")
}
