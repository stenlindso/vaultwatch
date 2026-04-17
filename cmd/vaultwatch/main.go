package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "vaultwatch",
	Short: "Audit and diff HashiCorp Vault secret paths across environments",
}

var snapshotCmd = &cobra.Command{
	Use:   "snapshot [environment]",
	Short: "Take a snapshot of Vault secret paths for a given environment",
	Args:  cobra.ExactArgs(1),
	RunE:  runSnapshot,
}

var diffCmd = &cobra.Command{
	Use:   "diff [env1] [env2]",
	Short: "Diff secret paths between two environments",
	Args:  cobra.ExactArgs(2),
	RunE:  runDiff,
}

var (
	flagSnapshotDir string
	flagVaultAddr   string
	flagVaultToken  string
	flagOutputJSON  bool
)

func init() {
	rootCmd.PersistentFlags().StringVar(&flagSnapshotDir, "snapshot-dir", ".vaultwatch", "Directory to store snapshots")
	rootCmd.PersistentFlags().StringVar(&flagVaultAddr, "vault-addr", "", "Vault server address (overrides VAULT_ADDR)")
	rootCmd.PersistentFlags().StringVar(&flagVaultToken, "vault-token", "", "Vault token (overrides VAULT_TOKEN)")
	diffCmd.Flags().BoolVar(&flagOutputJSON, "json", false, "Output diff as JSON")
	rootCmd.AddCommand(snapshotCmd)
	rootCmd.AddCommand(diffCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
