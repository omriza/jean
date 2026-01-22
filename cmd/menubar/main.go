package main

import (
	"fmt"
	"os"
	"time"

	"github.com/caseymrm/menuet"
	"github.com/omriza/jean/internal/api"
	"github.com/omriza/jean/internal/auth"
)

var (
	client *api.Client
	usage  *api.UsageResponse
)

func main() {
	// Get credentials from Claude Desktop
	creds, err := auth.GetCredentials()
	if err != nil {
		menuet.App().Label = "com.github.omri.jean"
		menuet.App().SetMenuState(&menuet.MenuState{
			Title: "⚠️ Jean",
		})
		menuet.App().RunApplication()
		return
	}

	client = api.NewClient(creds.SessionKey, creds.OrgID)

	// Configure the app
	menuet.App().Label = "com.github.omri.jean"
	menuet.App().Children = menuItems

	// Set initial title so it's visible while loading
	menuet.App().SetMenuState(&menuet.MenuState{
		Title: "C ░░░░░ --",
	})

	// Initial fetch
	go refreshUsage()

	// Start background refresh
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		for range ticker.C {
			refreshUsage()
		}
	}()

	// Run the app (this blocks)
	menuet.App().RunApplication()
}

func refreshUsage() {
	newUsage, err := client.GetUsage()
	if err != nil {
		menuet.App().SetMenuState(&menuet.MenuState{
			Title: "⚠️",
		})
		return
	}
	usage = newUsage
	updateMenuBar()
}

func updateMenuBar() {
	if usage == nil || usage.FiveHour == nil {
		return
	}

	// Create progress bar for current session with prefix
	// Using shorter bar (5 chars) + percentage to fit near notch
	bar := makeProgressBar(usage.FiveHour.Utilization, 5)
	title := fmt.Sprintf("C %s %.0f%%", bar, usage.FiveHour.Utilization)

	menuet.App().SetMenuState(&menuet.MenuState{
		Title: title,
	})
}

func makeProgressBar(percent float64, width int) string {
	filled := int(percent / 100 * float64(width))
	if filled > width {
		filled = width
	}
	// Show at least 1 filled if there's any usage
	if filled == 0 && percent > 0 {
		filled = 1
	}

	bar := ""
	for i := 0; i < width; i++ {
		if i < filled {
			bar += "▰"
		} else {
			bar += "▱"
		}
	}
	return bar
}

func menuItems() []menuet.MenuItem {
	items := []menuet.MenuItem{}

	if usage == nil {
		items = append(items, menuet.MenuItem{
			Text: "Loading...",
		})
	} else {
		// Current session
		if usage.FiveHour != nil {
			items = append(items, menuet.MenuItem{
				Text: "Current Session",
				FontWeight: menuet.WeightBold,
			})
			items = append(items, menuet.MenuItem{
				Text: fmt.Sprintf("   %s %.0f%%", makeProgressBar(usage.FiveHour.Utilization, 15), usage.FiveHour.Utilization),
			})
			items = append(items, menuet.MenuItem{
				Text: fmt.Sprintf("   Resets %s", formatResetTime(usage.FiveHour.ResetsAt)),
			})
		}

		items = append(items, menuet.MenuItem{Type: menuet.Separator})

		// Weekly limits header
		items = append(items, menuet.MenuItem{
			Text: "Weekly Limits",
			FontWeight: menuet.WeightBold,
		})

		// All models
		if usage.SevenDay != nil {
			items = append(items, menuet.MenuItem{
				Text: fmt.Sprintf("   All Models: %s %.0f%%", makeProgressBar(usage.SevenDay.Utilization, 10), usage.SevenDay.Utilization),
			})
			items = append(items, menuet.MenuItem{
				Text: fmt.Sprintf("   Resets %s", formatResetTime(usage.SevenDay.ResetsAt)),
			})
		}

		// Sonnet only
		if usage.SevenDaySonnet != nil {
			items = append(items, menuet.MenuItem{
				Text: fmt.Sprintf("   Sonnet: %s %.0f%%", makeProgressBar(usage.SevenDaySonnet.Utilization, 10), usage.SevenDaySonnet.Utilization),
			})
			items = append(items, menuet.MenuItem{
				Text: fmt.Sprintf("   Resets %s", formatResetTime(usage.SevenDaySonnet.ResetsAt)),
			})
		}

		items = append(items, menuet.MenuItem{Type: menuet.Separator})

		// Last updated
		items = append(items, menuet.MenuItem{
			Text: fmt.Sprintf("Updated: %s", time.Now().Format("3:04 PM")),
		})
	}

	items = append(items, menuet.MenuItem{Type: menuet.Separator})

	// Refresh button
	items = append(items, menuet.MenuItem{
		Text: "Refresh Now",
		Clicked: func() {
			go refreshUsage()
		},
	})

	// Quit
	items = append(items, menuet.MenuItem{
		Text: "Quit Jean",
		Clicked: func() {
			os.Exit(0)
		},
	})

	return items
}

func formatResetTime(t time.Time) string {
	now := time.Now()
	diff := t.Sub(now)

	if diff < 0 {
		return "now"
	}

	if diff < time.Hour {
		return fmt.Sprintf("in %d min", int(diff.Minutes()))
	}

	if diff < 24*time.Hour {
		hours := int(diff.Hours())
		mins := int(diff.Minutes()) % 60
		if mins > 0 {
			return fmt.Sprintf("in %dh %dm", hours, mins)
		}
		return fmt.Sprintf("in %dh", hours)
	}

	return t.Format("Mon 3:04 PM")
}
