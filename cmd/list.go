package cmd

import (
	"fmt"

	"github.com/paulleonhardhellweg/bring-tui/internal/bring"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Show items on your shopping list",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, stored, err := bring.Authenticate()
		if err != nil {
			return fmt.Errorf("not logged in. Run 'bring login' first: %w", err)
		}

		listUUID := stored.DefaultListUUID
		if listUUID == "" {
			return fmt.Errorf("no default list set. Run 'bring lists' and 'bring use <name>'")
		}

		items, err := client.GetItems(listUUID)
		if err != nil {
			return fmt.Errorf("failed to get items: %w", err)
		}

		listName := stored.DefaultListName
		if listName == "" {
			listName = "Shopping List"
		}

		if len(items.Purchase) == 0 {
			fmt.Printf("  %s is empty.\n", listName)
			return nil
		}

		fmt.Printf("  %s\n", listName)
		fmt.Println("  " + "────────────────────────────────")
		for _, item := range items.Purchase {
			if item.Spec != "" {
				fmt.Printf("  ● %s — %s\n", item.ItemID, item.Spec)
			} else {
				fmt.Printf("  ● %s\n", item.ItemID)
			}
		}

		if len(items.Recently) > 0 {
			fmt.Println()
			fmt.Println("  Recently bought:")
			for _, item := range items.Recently {
				fmt.Printf("  ✓ %s\n", item.ItemID)
			}
		}

		return nil
	},
}
