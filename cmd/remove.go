package cmd

import (
	"fmt"

	"github.com/paulleonhardhellweg/bring-tui/internal/bring"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(removeCmd)
}

var removeCmd = &cobra.Command{
	Use:     "remove <item>",
	Short:   "Remove an item from the list",
	Args:    cobra.ExactArgs(1),
	Example: "  bring remove Milch",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, stored, err := bring.Authenticate()
		if err != nil {
			return fmt.Errorf("not logged in. Run 'bring login' first: %w", err)
		}

		listUUID := stored.DefaultListUUID
		if listUUID == "" {
			return fmt.Errorf("no default list set")
		}

		if err := client.RemoveItem(listUUID, args[0]); err != nil {
			return fmt.Errorf("failed to remove item: %w", err)
		}

		fmt.Printf("Removed: %s\n", args[0])
		return nil
	},
}
