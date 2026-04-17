package main

importf13/cobra"

	"github.com/vaultwatch/internal/audit"
	"github.com/vaultwatch/internal/snapshot"
	"github.com/vaultwatch/internal/vault"
)

func runSnapshot(cmd *cobra.Command, args []string) error {
	env := args[0]

	client, err := vault.NewClient(vault.Config{
		Address: flagVaultAddr,
		Token:   flagVaultToken,
	})
	if err != nil {
		return fmt.Errorf("vault client: %w", err)
	}

	lister := vault.NewLister(client)
	paths, err := lister.ListPaths(cmd.Context(), "secret/")
	if err != nil {
		return fmt.Errorf("listing paths: %w", err)
	}

	mgr := snapshot.NewManager(flagSnapshotDir)
	if err := mgr.Save(env, paths); err != nil {
		return fmt.Errorf("saving snapshot: %w", err)
	}

	fmt.Fprintf(os.Stdout, "Snapshot saved for environment %q (%d paths)\n", env, len(paths))
	return nil
}

func runDiff(cmd *cobra.Command, args []string) error {
	env1, env2 := args[0], args[1]

	mgr := snapshot.NewManager(flagSnapshotDir)
	auditor := audit.NewAuditor(mgr)

	report, err := auditor.Audit(env1, env2)
	if err != nil {
		return fmt.Errorf("audit: %w", err)
	}

	if flagOutputJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(report)
	}

	audit.Format(os.Stdout, report)
	return nil
}
