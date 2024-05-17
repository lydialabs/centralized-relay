package cmd

import (
	"github.com/icon-project/centralized-relay/relayer/chains/bitcoin"
	"github.com/spf13/cobra"
)

func btcCmd(a *appState) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "bitcoin",
		Aliases: []string{"btc"},
		Short:   "Command line of bitcoin",
	}

	cmd.AddCommand(
		syncBlockEvent(a),
	)

	return cmd
}

func syncBlockEvent(a *appState) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "sync",
		Aliases: []string{"sync"},
		Short:   "Sync block on-chain of bitcoin",
		Args:    withUsage(cobra.NoArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			config, _ := bitcoin.LoadConfig()
			bitcoin.ListenerBitcoin(config)
			return nil
		},
	}
	return yamlFlag(a.viper, jsonFlag(a.viper, cmd))
}
