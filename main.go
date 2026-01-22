package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/omri/jean/internal/api"
	"github.com/omri/jean/internal/auth"
	"github.com/omri/jean/internal/tui"
)

func main() {
	// Get credentials from Claude Desktop (single keychain prompt)
	creds, err := auth.GetCredentials()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get credentials: %v\n", err)
		fmt.Fprintf(os.Stderr, "\nMake sure Claude Desktop is installed and you're logged in.\n")
		fmt.Fprintf(os.Stderr, "Tip: Click 'Always Allow' when prompted for keychain access.\n")
		os.Exit(1)
	}

	// Create API client
	client := api.NewClient(creds.SessionKey, creds.OrgID)

	// Start TUI
	p := tea.NewProgram(tui.NewModel(client), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running TUI: %v\n", err)
		os.Exit(1)
	}
}
