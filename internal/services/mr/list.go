package mr

import (
	"fmt"
	"strings"
	"time"

	"gitbroski/internal/services"
	"gitbroski/utils/logger"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// displayItem holds MR data plus fetched status for display.
type displayItem struct {
	MR     MR
	Status Status
}

// listModel is the bubbletea model for the MR list.
type listModel struct {
	items    []displayItem
	cursor   int
	quitting bool
	showHelp bool
}

func (m listModel) Init() tea.Cmd {
	return nil
}

func (m listModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.showHelp {
			m.showHelp = false
			return m, nil
		}

		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.items)-1 {
				m.cursor++
			}
		case "enter":
			if len(m.items) > 0 {
				services.OpenBrowser(m.items[m.cursor].MR.URL)
			}
		case "d":
			if len(m.items) > 0 {
				m.items = append(m.items[:m.cursor], m.items[m.cursor+1:]...)
				if m.cursor >= len(m.items) && m.cursor > 0 {
					m.cursor--
				}
				mrs := make([]MR, len(m.items))
				for i, item := range m.items {
					mrs[i] = item.MR
				}
				SaveStore(&Store{MRs: mrs})
			}
		case "h":
			m.showHelp = true
		}
	}
	return m, nil
}

func (m listModel) View() string {
	if m.quitting {
		return ""
	}

	if m.showHelp {
		return renderAuthHelp()
	}

	if len(m.items) == 0 {
		box := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(BorderColor).
			Padding(1, 2).
			Render("No MRs saved yet.\n\nUse " + CodeStyle.Render("gitbroski mr save <url>") + " to add one.")
		return "\n" + box + "\n"
	}

	var b strings.Builder

	b.WriteString("\n")
	b.WriteString(TitleStyle.Render("  📋 Saved MRs") + "\n\n")

	header := fmt.Sprintf("   %-3s │ %-32s │ %-10s │ %-10s │ %-8s", "#", "Title", "Status", "Service", "Saved")
	b.WriteString(DimStyle.Render(header) + "\n")
	b.WriteString(DimStyle.Render("  ─────┼──────────────────────────────────┼────────────┼────────────┼─────────") + "\n")

	for i, item := range m.items {
		title := item.MR.Title
		if title == "" {
			title = item.Status.Title
		}
		if title == "" {
			title = "(no title)"
		}
		if len(title) > 30 {
			title = title[:27] + "..."
		}

		service := item.MR.Service
		if service == "" {
			service = item.Status.Repo
		}
		if len(service) > 10 {
			service = service[:7] + "..."
		}

		savedAgo := formatTimeAgo(item.MR.SavedAt)
		status := fmt.Sprintf("%s %-6s", item.Status.Icon(), item.Status.Label())

		if i == m.cursor {
			line := fmt.Sprintf(" ▸ %-3d │ %-32s │ %-10s │ %-10s │ %-8s", i+1, title, status, service, savedAgo)
			b.WriteString(SelectedStyle.Render(line) + "\n")
		} else {
			line := fmt.Sprintf("   %-3d │ %-32s │ %-10s │ %-10s │ %-8s", i+1, title, status, service, savedAgo)
			b.WriteString(NormalStyle.Render(line) + "\n")
		}
	}

	b.WriteString("\n")
	keys := []string{"↑↓", "enter", "d", "h", "q"}
	actions := []string{"navigate", "open", "delete", "help", "quit"}
	var helpParts []string
	for i := range keys {
		key := lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Render(keys[i])
		helpParts = append(helpParts, key+" "+actions[i])
	}
	b.WriteString(HelpStyle.Render("  " + strings.Join(helpParts, "  •  ")))
	b.WriteString("\n")

	return b.String()
}

func renderAuthHelp() string {
	var b strings.Builder

	b.WriteString(AuthTitleStyle.Render("🔐 Authentication Setup") + "\n\n")

	b.WriteString(SectionStyle.Render("GitHub (gh CLI)") + "\n")
	b.WriteString("─────────────────\n")
	b.WriteString("1. Install: " + CodeStyle.Render("brew install gh") + "\n")
	b.WriteString("2. Login:   " + CodeStyle.Render("gh auth login") + "\n")
	b.WriteString("3. Verify:  " + CodeStyle.Render("gh auth status") + "\n\n")

	b.WriteString(SectionStyle.Render("GitLab (glab CLI)") + "\n")
	b.WriteString("─────────────────\n")
	b.WriteString("1. Install: " + CodeStyle.Render("brew install glab") + "\n")
	b.WriteString("2. Login:   " + CodeStyle.Render("glab auth login") + "\n")
	b.WriteString("3. Verify:  " + CodeStyle.Render("glab auth status") + "\n\n")

	b.WriteString(HelpStyle.Render("Press any key to close"))

	return b.String()
}

func formatTimeAgo(t time.Time) string {
	diff := time.Since(t)

	switch {
	case diff < time.Minute:
		return "just now"
	case diff < time.Hour:
		mins := int(diff.Minutes())
		if mins == 1 {
			return "1m ago"
		}
		return fmt.Sprintf("%dm ago", mins)
	case diff < 24*time.Hour:
		hours := int(diff.Hours())
		if hours == 1 {
			return "1h ago"
		}
		return fmt.Sprintf("%dh ago", hours)
	default:
		days := int(diff.Hours() / 24)
		if days == 1 {
			return "1d ago"
		}
		return fmt.Sprintf("%dd ago", days)
	}
}

// List displays an interactive list of saved MRs.
func List() {
	store, err := LoadStore()
	if err != nil {
		logger.Error("Failed to load MR store: " + err.Error())
		return
	}

	items := make([]displayItem, len(store.MRs))
	for i, m := range store.MRs {
		items[i] = displayItem{
			MR:     m,
			Status: FetchStatus(m.URL),
		}
	}

	model := listModel{items: items}
	p := tea.NewProgram(model)

	_, err = p.Run()
	if err != nil {
		logger.Error("Error running list: " + err.Error())
		return
	}
}
