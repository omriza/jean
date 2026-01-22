package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/omriza/jean/internal/api"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("15"))

	subtitleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8"))

	barEmptyStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8"))

	barFillStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("4"))

	percentStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("7")).
			Width(10).
			Align(lipgloss.Right)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("1"))

	refreshStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("8"))
)

type tickMsg time.Time
type usageMsg *api.UsageResponse
type errMsg error

type Model struct {
	client      *api.Client
	usage       *api.UsageResponse
	lastUpdated time.Time
	err         error
	spinner     spinner.Model
	loading     bool
	width       int
}

func NewModel(client *api.Client) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))

	return Model{
		client:  client,
		spinner: s,
		loading: true,
		width:   60,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		fetchUsage(m.client),
		tickEvery(30*time.Second),
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "r":
			m.loading = true
			return m, fetchUsage(m.client)
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width

	case tickMsg:
		m.loading = true
		return m, tea.Batch(
			fetchUsage(m.client),
			tickEvery(30*time.Second),
		)

	case usageMsg:
		m.usage = msg
		m.lastUpdated = time.Now()
		m.loading = false
		m.err = nil
		return m, nil

	case errMsg:
		m.err = msg
		m.loading = false
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m Model) View() string {
	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(titleStyle.Render("Plan usage limits"))
	b.WriteString("\n\n")

	if m.err != nil {
		b.WriteString(errorStyle.Render(fmt.Sprintf("Error: %v", m.err)))
		b.WriteString("\n\n")
	}

	if m.usage != nil {
		// Current session (five_hour)
		if m.usage.FiveHour != nil {
			b.WriteString(m.renderUsageRow("Current session", m.usage.FiveHour))
			b.WriteString("\n\n")
		}

		// Weekly limits section
		b.WriteString(titleStyle.Render("Weekly limits"))
		b.WriteString("\n")
		b.WriteString(subtitleStyle.Render("Learn more about usage limits"))
		b.WriteString("\n\n")

		// All models (seven_day)
		if m.usage.SevenDay != nil {
			b.WriteString(m.renderUsageRow("All models", m.usage.SevenDay))
			b.WriteString("\n\n")
		}

		// Sonnet only (seven_day_sonnet)
		if m.usage.SevenDaySonnet != nil {
			b.WriteString(m.renderUsageRow("Sonnet only", m.usage.SevenDaySonnet))
			b.WriteString("\n\n")
		}
	} else if m.loading {
		b.WriteString(m.spinner.View())
		b.WriteString(" Loading usage data...")
		b.WriteString("\n\n")
	}

	// Last updated
	if !m.lastUpdated.IsZero() {
		ago := time.Since(m.lastUpdated).Round(time.Second)
		loadingIndicator := ""
		if m.loading {
			loadingIndicator = " " + m.spinner.View()
		}
		b.WriteString(refreshStyle.Render(fmt.Sprintf("Last updated: %s ago%s", ago, loadingIndicator)))
	}

	b.WriteString("\n\n")
	b.WriteString(subtitleStyle.Render("Press 'r' to refresh, 'q' to quit"))
	b.WriteString("\n")

	return b.String()
}

func (m Model) renderUsageRow(label string, limit *api.UsageLimit) string {
	var b strings.Builder

	// Label and reset time
	b.WriteString(titleStyle.Render(label))
	b.WriteString("\n")
	b.WriteString(subtitleStyle.Render(formatResetTime(limit.ResetsAt)))
	b.WriteString("\n")

	// Progress bar
	barWidth := 40
	filled := int(limit.Utilization / 100 * float64(barWidth))
	if filled > barWidth {
		filled = barWidth
	}

	bar := barFillStyle.Render(strings.Repeat("█", filled))
	bar += barEmptyStyle.Render(strings.Repeat("░", barWidth-filled))

	percent := percentStyle.Render(fmt.Sprintf("%.0f%% used", limit.Utilization))

	b.WriteString(bar)
	b.WriteString("  ")
	b.WriteString(percent)

	return b.String()
}

func formatResetTime(t time.Time) string {
	now := time.Now()
	diff := t.Sub(now)

	if diff < 0 {
		return "Resetting..."
	}

	if diff < time.Hour {
		return fmt.Sprintf("Resets in %d min", int(diff.Minutes()))
	}

	if diff < 24*time.Hour {
		hours := int(diff.Hours())
		mins := int(diff.Minutes()) % 60
		if mins > 0 {
			return fmt.Sprintf("Resets in %dh %dm", hours, mins)
		}
		return fmt.Sprintf("Resets in %dh", hours)
	}

	return fmt.Sprintf("Resets %s", t.Format("Mon 3:04 PM"))
}

func fetchUsage(client *api.Client) tea.Cmd {
	return func() tea.Msg {
		usage, err := client.GetUsage()
		if err != nil {
			return errMsg(err)
		}
		return usageMsg(usage)
	}
}

func tickEvery(d time.Duration) tea.Cmd {
	return tea.Tick(d, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
