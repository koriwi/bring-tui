package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/paulleonhardhellweg/bring-tui/internal/bring"
	"github.com/paulleonhardhellweg/bring-tui/internal/tui"
	"github.com/spf13/cobra"
)

var refreshFlag bool

var rootCmd = &cobra.Command{
	Use:   "bring",
	Short: "Bring! shopping list in your terminal",
	Long:  "A TUI for the Bring! shopping list app. Add items from your terminal, because code is life.",
	RunE: func(cmd *cobra.Command, args []string) error {
		if refreshFlag {
			if err := bring.RefreshStoredToken(); err != nil {
				return err
			}
			fmt.Println("Token refreshed.")
			return nil
		}
		p := tea.NewProgram(tui.NewApp(), tea.WithAltScreen())
		if _, err := p.Run(); err != nil {
			return fmt.Errorf("TUI error: %w", err)
		}
		return nil
	},
}

func init() {
	rootCmd.Flags().BoolVar(&refreshFlag, "refresh", false, "Refresh the stored auth token and exit")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
