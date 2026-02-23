package cmd

import (
	"fmt"
	"strings"

	"github.com/paulleonhardhellweg/bring-tui/internal/bring"
	"github.com/paulleonhardhellweg/bring-tui/internal/config"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(useCmd)
}

var useCmd = &cobra.Command{
	Use:     "use <list-name>",
	Short:   "Set the default shopping list",
	Args:    cobra.ExactArgs(1),
	Example: "  bring use Wochenmarkt",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, stored, err := bring.Authenticate()
		if err != nil {
			return fmt.Errorf("not logged in. Run 'bring login' first: %w", err)
		}

		lists, err := client.GetLists()
		if err != nil {
			return fmt.Errorf("failed to get lists: %w", err)
		}

		name := args[0]
		for _, l := range lists {
			if strings.EqualFold(l.Name, name) {
				stored.DefaultListUUID = l.ListUUID
				stored.DefaultListName = l.Name
				if err := config.Save(stored); err != nil {
					return fmt.Errorf("failed to save config: %w", err)
				}
				fmt.Printf("Default list set to: %s\n", l.Name)
				return nil
			}
		}

		fmt.Printf("List '%s' not found. Available lists:\n", name)
		for _, l := range lists {
			fmt.Printf("  - %s\n", l.Name)
		}
		return nil
	},
}
