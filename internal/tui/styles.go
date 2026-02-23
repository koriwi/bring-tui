package tui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	colorPrimary = lipgloss.Color("#4A7C59") // Sage Green
	colorMuted   = lipgloss.Color("#6B7280")
	colorDim     = lipgloss.Color("#9CA3AF")
	colorText    = lipgloss.Color("#1A1A1A")
	colorError   = lipgloss.Color("#DC2626")
	colorSuccess = lipgloss.Color("#16A34A")

	// Title
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorPrimary).
			PaddingLeft(1)

	// Items
	itemStyle = lipgloss.NewStyle().
			Foreground(colorText)

	selectedItemStyle = lipgloss.NewStyle().
				Foreground(colorPrimary).
				Bold(true)

	specStyle = lipgloss.NewStyle().
			Foreground(colorMuted)

	doneItemStyle = lipgloss.NewStyle().
			Foreground(colorDim).
			Strikethrough(true)

	// Status bar
	statusBarStyle = lipgloss.NewStyle().
			Foreground(colorMuted).
			PaddingLeft(1).
			PaddingTop(1)

	helpKeyStyle = lipgloss.NewStyle().
			Foreground(colorPrimary).
			Bold(true)

	helpDescStyle = lipgloss.NewStyle().
			Foreground(colorMuted)

	// Messages
	errorStyle = lipgloss.NewStyle().
			Foreground(colorError).
			PaddingLeft(1)

	successStyle = lipgloss.NewStyle().
			Foreground(colorSuccess).
			PaddingLeft(1)

	// Divider
	dividerStyle = lipgloss.NewStyle().
			Foreground(colorDim).
			PaddingLeft(1)

	// Section header (e.g., "recently bought")
	sectionStyle = lipgloss.NewStyle().
			Foreground(colorDim).
			Italic(true).
			PaddingLeft(1)
)
