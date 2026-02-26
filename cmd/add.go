package cmd

import (
	"fmt"
	"strings"

	"github.com/paulleonhardhellweg/bring-tui/internal/bring"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use:   "add <item>[:<description>]",
	Short: "Add an item to your shopping list",
	Long:  `Add an item to your default Bring! shopping list. Optionally include a description after a colon.`,
	Example: `  bring add Milch
  bring add "Milch:1.5%, 2 Liter"
  bring add Milch:1.5%`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, stored, err := bring.Authenticate()
		if err != nil {
			return fmt.Errorf("not logged in. Run 'bring login' first: %w", err)
		}

		item, spec, _ := strings.Cut(strings.Join(args, " "), ":")
		item = strings.TrimSpace(item)
		spec = strings.TrimSpace(spec)

		listUUID := stored.DefaultListUUID
		if listUUID == "" {
			return fmt.Errorf("no default list set. Run 'bring lists' and 'bring use <name>'")
		}

		if err := client.AddItem(listUUID, item, spec); err != nil {
			return fmt.Errorf("failed to add item: %w", err)
		}

		if spec != "" {
			fmt.Printf("Added: %s — %s\n", item, spec)
		} else {
			fmt.Printf("Added: %s\n", item)
		}
		return nil
	},
}
