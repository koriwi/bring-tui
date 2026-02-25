package cmd

import (
	"fmt"

	"github.com/paulleonhardhellweg/bring-tui/internal/bring"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(doneCmd)
}

var doneCmd = &cobra.Command{
	Use:     "done <item>",
	Short:   "Mark an item as bought",
	Args:    cobra.ExactArgs(1),
	Example: "  bring done Milch",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, stored, err := bring.Authenticate()
		if err != nil {
			return fmt.Errorf("not logged in. Run 'bring login' first: %w", err)
		}

		listUUID := stored.DefaultListUUID
		if listUUID == "" {
			return fmt.Errorf("no default list set")
		}

		if err := client.CompleteItem(listUUID, args[0], ""); err != nil {
			return fmt.Errorf("failed to complete item: %w", err)
		}

		fmt.Printf("Done: %s ✓\n", args[0])
		return nil
	},
}
