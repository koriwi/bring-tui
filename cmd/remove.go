package cmd

import (
	"fmt"
	"strings"

	"github.com/paulleonhardhellweg/bring-tui/internal/bring"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(removeCmd)
}

var removeCmd = &cobra.Command{
	Use:     "remove <item>",
	Short:   "Remove an item from the list",
	Args:    cobra.MinimumNArgs(1),
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

		item := strings.TrimSpace(strings.Join(args, " "))
		if err := client.RemoveItem(listUUID, item); err != nil {
			return fmt.Errorf("failed to remove item: %w", err)
		}

		fmt.Printf("Removed: %s\n", item)
		return nil
	},
}
