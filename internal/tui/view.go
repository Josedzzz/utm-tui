// Package tui handles the screens
package tui

import (
	"fmt"
	"strings"

	"github.com/Josedzzz/utm-tui/internal/utm"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const logo = `
████████╗██╗   ██╗██╗    ██╗   ██╗████████╗███╗   ███╗
╚══██╔══╝██║   ██║██║    ██║   ██║╚══██╔══╝████╗ ████║
   ██║   ██║   ██║██║    ██║   ██║   ██║   ██╔████╔██║
   ██║   ██║   ██║██║    ██║   ██║   ██║   ██║╚██╔╝██║
   ██║   ╚██████╔╝██║    ╚██████╔╝   ██║   ██║ ╚═╝ ██║
   ╚═╝    ╚═════╝ ╚═╝     ╚═════╝    ╚═╝   ╚═╝     ╚═╝
`

const miniLogo = " UTMctl "

// vmsLoadedMsg is an internal message type (defined here)
type vmsLoadedMsg struct {
	vms []utm.VM
	err error
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles all messages and user input
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}
		switch m.state {
		case menuView:
			return m.handleMenuUpdate(msg)
		case listVMsView:
			return m.handleListVMsUpdate(msg)
		}

	case vmsLoadedMsg:
		m.loading = false
		if msg.err != nil {
			m.message = fmt.Sprintf("Error loading VMs: %v", msg.err)
			m.isSuccess = false
		} else {
			m.vms = msg.vms
			m.message = fmt.Sprintf("Loaded %d VMs", len(m.vms))
			m.isSuccess = true
		}
		return m, nil
	}
	return m, nil
}

// handleMenuUpdate processes keys in the main menu
func (m Model) handleMenuUpdate(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.choices)-1 {
			m.cursor++
		}
	case "enter":
		switch m.choices[m.cursor] {
		case "List VMs":
			m.state = listVMsView
			m.loading = true
			m.message = ""
			return m, fetchVMsCmd
		case "Exit":
			m.quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

// handleListVMsUpdate processes keys in the VM list view
func (m Model) handleListVMsUpdate(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "backspace":
		m.state = menuView
		m.vms = []utm.VM{}
		m.message = ""
		return m, nil
	case "up", "k":
		if m.vmsCursor > 0 {
			m.vmsCursor--
			if m.vmsCursor < m.vmsOffset {
				m.vmsOffset = m.vmsCursor
			}
		}
	case "down", "j":
		if m.vmsCursor < len(m.vms)-1 {
			m.vmsCursor++
			if m.vmsCursor >= m.vmsOffset+10 {
				m.vmsOffset = m.vmsCursor - 9
			}
		}
	}
	return m, nil
}

// fetchVMsCmd runs utmctl list and returns a vmsLoadedMsg
func fetchVMsCmd() tea.Msg {
	vms, err := utm.ListVMs()
	if err != nil {
		return vmsLoadedMsg{err: err}
	}
	return vmsLoadedMsg{vms: vms, err: nil}
}

// View renders the current screen
func (m Model) View() string {
	if m.quitting {
		return ""
	}

	var body strings.Builder

	if m.state == menuView {
		body.WriteString(asciiStyle.Render(logo) + "\n")
		body.WriteString(titleStyle.Render(" MAIN MENU ") + "\n\n")

		for i, choice := range m.choices {
			if m.cursor == i {
				fmt.Fprintf(&body, "%s\n", selectedItemStyle.Render(choice))
			} else {
				fmt.Fprintf(&body, "%s\n", itemStyle.Render(choice))
			}
		}
		body.WriteString(footerStyle.Render("\n↑/↓: navigate • enter: select • q/Ctrl+C: quit"))
	} else {
		// Submenu: mini logo
		body.WriteString(miniLogoStyle.Render(miniLogo) + "\n")

		switch m.state {
		case listVMsView:
			body.WriteString(titleStyle.Render(" VIRTUAL MACHINES ") + "\n\n")

			if m.loading {
				body.WriteString(processingStyle.Render("Loading VMs... Please wait.") + "\n")
			} else if len(m.vms) == 0 {
				body.WriteString("No VMs found.\n\n")
			} else {
				// Table header with proper formatting
				fmt.Fprintf(&body, "%-30s  STATUS\n", "NAME")
				fmt.Fprintf(&body, "%s\n", strings.Repeat("-", 50))

				// Show visible items (window size ~10)
				start := m.vmsOffset
				end := min(start+10, len(m.vms))
				for i := start; i < end; i++ {
					vm := m.vms[i]
					statusLower := strings.ToLower(vm.Status)
					statusStyle, ok := vmStatusStyle[statusLower]
					if !ok {
						statusStyle = lipgloss.NewStyle()
					}
					// Format with proper spacing: name (30 chars) + 2 spaces + status
					name := fmt.Sprintf("%-30s", vm.Name)
					status := statusStyle.Render(vm.Status)
					line := fmt.Sprintf("%s  %s", name, status)
					if m.vmsCursor == i {
						line = selectedItemStyle.Render(line)
					}
					fmt.Fprintf(&body, "%s\n", line)
				}

				// Scroll indicator
				if len(m.vms) > 10 {
					scrollPercent := int((float64(m.vmsOffset+10) / float64(len(m.vms))) * 100)
					fmt.Fprintf(&body, "\n%s", statusStyle.Render(fmt.Sprintf("Showing %d-%d of %d (%d%%)", start+1, end, len(m.vms), scrollPercent)))
				}
			}

			if m.message != "" {
				style := errorStyle
				if m.isSuccess {
					style = successStyle
				}
				fmt.Fprintf(&body, "\n%s", style.Render(m.message))
			}

			body.WriteString(footerStyle.Render("\n\n↑/↓: scroll • esc: back to menu"))
		}
	}

	// Center the content
	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		containerStyle.Render(body.String()),
	)
}
