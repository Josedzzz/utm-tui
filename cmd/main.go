package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// VM represents a single virtual machine
type VM struct {
	Name   string
	Status string
}

// model holds the TUI state
type model struct {
	vms      []VM
	err      error
	quitting bool
}

// fetchVMsCmd is a tea.Cmd that runs utmctl list and returns the result
type vmsLoadedMsg struct {
	vms []VM
	err error
}

func fetchVMsCmd() tea.Msg {
	cmd := exec.Command("utmctl", "list")
	out, err := cmd.Output()
	if err != nil {
		return vmsLoadedMsg{err: fmt.Errorf("failed to run utmctl: %w", err)}
	}

	// Parse the table output
	lines := strings.Split(string(out), "\n")
	var vms []VM
	for i, line := range lines {
		// Skip header and empty lines
		if i == 0 || strings.TrimSpace(line) == "" {
			continue
		}
		fields := strings.Fields(line)
		// Expected format: "Name   Status   [UUID]"
		// Usually first two fields are Name and Status; if more, the first is name, second status.
		if len(fields) >= 2 {
			vms = append(vms, VM{
				Name:   fields[0],
				Status: fields[1],
			})
		}
	}
	return vmsLoadedMsg{vms: vms, err: nil}
}

// Init initializes the model and returns the first command
func (m model) Init() tea.Cmd {
	return fetchVMsCmd
}

// Update handles messages and user input
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quitting = true
			return m, tea.Quit
		}

	case vmsLoadedMsg:
		if msg.err != nil {
			m.err = msg.err
			return m, tea.Quit
		}
		m.vms = msg.vms
		return m, nil
	}
	return m, nil
}

// View renders the screen
func (m model) View() string {
	if m.quitting {
		return "Goodbye!\n"
	}

	if m.err != nil {
		return fmt.Sprintf("Error: %v\nPress any key to exit.", m.err)
	}

	if len(m.vms) == 0 {
		return "No VMs found.\nPress q or Ctrl+C to quit."
	}

	var b strings.Builder
	title := lipgloss.NewStyle().Bold(true).Underline(true).Foreground(lipgloss.Color("99"))
	header := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))
	statusStyle := map[string]lipgloss.Style{
		"running":   lipgloss.NewStyle().Foreground(lipgloss.Color("10")), // green
		"stopped":   lipgloss.NewStyle().Foreground(lipgloss.Color("9")),  // red
		"suspended": lipgloss.NewStyle().Foreground(lipgloss.Color("11")), // yellow
	}

	fmt.Fprintf(&b, "%s\n\n", title.Render("UTM Virtual Machines"))
	fmt.Fprintf(&b, "%-30s %s\n", header.Render("NAME"), header.Render("STATUS"))
	fmt.Fprintf(&b, "%s\n", strings.Repeat("-", 45))

	for _, vm := range m.vms {
		status := strings.ToLower(vm.Status)
		style, ok := statusStyle[status]
		if !ok {
			style = lipgloss.NewStyle()
		}
		fmt.Fprintf(&b, "%-30s %s\n", vm.Name, style.Render(vm.Status))
	}
	fmt.Fprintf(&b, "\nPress q or Ctrl+C to quit.")
	return b.String()
}

func main() {
	p := tea.NewProgram(model{})
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
