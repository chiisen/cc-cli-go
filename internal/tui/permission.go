package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/liao-eli/cc-cli-go/internal/permission"
)

type PermissionDialog struct {
	toolName string
	input    map[string]interface{}
	decision *permission.Decision

	selected int
	options  []string

	approved bool
	finished bool
}

func NewPermissionDialog(toolName string, input map[string]interface{}, decision *permission.Decision) PermissionDialog {
	return PermissionDialog{
		toolName: toolName,
		input:    input,
		decision: decision,
		selected: 0,
		options:  []string{"Allow", "Deny", "Always Allow", "Always Deny"},
		approved: false,
		finished: false,
	}
}

func (d PermissionDialog) Update(msg tea.Msg) (PermissionDialog, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyLeft:
			if d.selected > 0 {
				d.selected--
			}

		case tea.KeyRight:
			if d.selected < len(d.options)-1 {
				d.selected++
			}

		case tea.KeyEnter:
			d.approved = d.selected == 0 || d.selected == 2
			d.finished = true
			return d, tea.Quit

		case tea.KeyEscape:
			d.approved = false
			d.finished = true
			return d, tea.Quit
		}
	}

	return d, nil
}

func (d PermissionDialog) View() string {
	var b strings.Builder

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("12")).
		Padding(1, 2).
		Margin(1, 2)

	titleStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("12")).
		Bold(true)

	warningStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("9")).
		Bold(true)

	optionStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("15"))

	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("10")).
		Bold(true)

	b.WriteString(titleStyle.Render("⚠️  Permission Request"))
	b.WriteString("\n\n")

	b.WriteString(fmt.Sprintf("Tool: %s\n", d.toolName))
	b.WriteString("\n")

	b.WriteString("Input:\n")
	for key, value := range d.input {
		if strValue, ok := value.(string); ok {
			if len(strValue) > 100 {
				strValue = strValue[:100] + "..."
			}
			b.WriteString(fmt.Sprintf("  %s: %s\n", key, strValue))
		}
	}

	if d.decision.Reason != "" {
		b.WriteString("\n")
		if strings.Contains(d.decision.Reason, "dangerous") {
			b.WriteString(warningStyle.Render("⚠️  Warning: " + d.decision.Reason))
		} else {
			b.WriteString(fmt.Sprintf("Reason: %s", d.decision.Reason))
		}
	}

	b.WriteString("\n\n")

	for i, option := range d.options {
		if i == d.selected {
			b.WriteString(selectedStyle.Render(fmt.Sprintf("▶ %s", option)))
		} else {
			b.WriteString(optionStyle.Render(fmt.Sprintf("  %s", option)))
		}
		if i < len(d.options)-1 {
			b.WriteString("  ")
		}
	}

	b.WriteString("\n\n")
	b.WriteString("← → to select, Enter to confirm, Esc to deny")

	return borderStyle.Render(b.String())
}

func (d PermissionDialog) GetDecision() (bool, string) {
	if !d.finished {
		return false, ""
	}

	action := d.options[d.selected]
	return d.approved, action
}
