// running node
package main

import (
	"context"
	"fmt"

	"github.com/kalam360/obboy-blockchain-golang/database"
	"github.com/kalam360/obboy-blockchain-golang/node"
	"github.com/spf13/cobra"
)


func runCmd() *cobra.Command {
	var runCmd = &cobra.Command{
		Use: "run",
		Short: "Launches the Obboy node and its HTTP API.",
		Run: func(cmd *cobra.Command, args []string) {
			miner, _ := cmd.Flags().GetString(flagMiner)
			sslEmail, _ := cmd.Flags().GetString(flagSSLEmail)
			isSSLDisabled, _ := cmd.Flags().GetBool(flagDisableSSL)
			ip, _ := cmd.Flags().GetString(flagIP)
			port, _ := cmd.Flags().GetUint64(flagPort)
			bootstrapIp, _ := cmd.Flags().GetString(flagBootstrapIp)
			bootstrapPort, _ := cmd.Flags().GetUint64(flagBootstrapPort)
			bootstrapAcc, _ := cmd.Flags().GetString(flagBootstrapAcc)

			fmt.Println("Launching Obboy node and its HTTP API...")
			

			bootstrap := node.NewPeerNode(
				bootstrapIp,
				bootstrapPort,
				true,
				database.NewAccount(bootstrapAcc),
				false,
				"",
			)

			if !isSSLDisabled {
				port = node.HttpSSLPort
			}

			version := fmt.Sprintf("%s.%s.%s-alpha %s %s", Major, Minor, Fix, shortGitCommit(GitCommit), verbal)
			n := node.New(getDataDirFromCmd(cmd), ip, port, database.NewAccount(miner), bootstrap, version, node.DefaultMiningDifficulty)
			err := n.Run(context.Background(), isSSLDisabled, sslEmail)

		}
	}
}