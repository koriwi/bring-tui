package cmd

import (
	"fmt"

	"github.com/paulleonhardhellweg/bring-tui/internal/bring"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listsCmd)
}

var listsCmd = &cobra.Command{
	Use:   "lists",
	Short: "Show all your Bring! lists",
	RunE: func(cmd *cobra.Command, args []string) error {
		client, stored, err := bring.Authenticate()
		if err != nil {
			return fmt.Errorf("not logged in. Run 'bring login' first: %w", err)
		}

		lists, err := client.GetLists()
		if err != nil {
			return fmt.Errorf("failed to get lists: %w", err)
		}

		fmt.Println("Your lists:")
		for _, l := range lists {
			marker := "  "
			if l.ListUUID == stored.DefaultListUUID {
				marker = "▸ "
			}
			fmt.Printf("%s%s\n", marker, l.Name)
		}
		return nil
	},
}
