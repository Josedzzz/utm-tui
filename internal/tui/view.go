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

// actionCompleteMsg indicates an action (delete/clone) completed
type actionCompleteMsg struct {
	success bool
	message string
}

// statusUpdateMsg indicates a status update for a VM
type statusUpdateMsg struct {
	status  string
	message string
	err     error
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
		case deleteConfirmView:
			return m.handleDeleteConfirmUpdate(msg)
		case cloneInputView:
			return m.handleCloneInputUpdate(msg)
		case vmDetailsView:
			return m.handleVMDetailsUpdate(msg)
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

	case actionCompleteMsg:
		m.loading = false
		m.message = msg.message
		m.isSuccess = msg.success
		if msg.success {
			// After successful action, refresh the VM list
			m.state = listVMsView
			m.loading = true
			return m, fetchVMsCmd
		}
		return m, nil

	case statusUpdateMsg:
		m.loading = false
		if msg.err != nil {
			m.message = fmt.Sprintf("Status check failed: %v", msg.err)
			m.isSuccess = false
		} else {
			m.message = fmt.Sprintf("Status: %s", msg.status)
			m.isSuccess = true
			if m.selectedVM != nil {
				m.selectedVM.Status = msg.status
			}
		}
		return m, nil
	}
	return m, nil
}

// handleDeleteConfirmUpdate processes keys in the delete confirmation dialog
func (m Model) handleDeleteConfirmUpdate(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "enter":
		// Confirm delete
		if m.selectedVM != nil {
			m.loading = true
			return m, deleteVMCmd(m.selectedVM.UUID)
		}
	case "n", "esc":
		// Cancel delete
		m.state = listVMsView
		m.selectedVM = nil
		return m, nil
	}
	return m, nil
}

// handleCloneInputUpdate processes keys in the clone name input
func (m Model) handleCloneInputUpdate(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter":
		// Confirm clone with the entered name
		if m.selectedVM != nil && m.inputText != "" {
			m.loading = true
			return m, cloneVMCmd(m.selectedVM.UUID, m.inputText)
		}
	case "esc":
		// Cancel clone
		m.state = listVMsView
		m.selectedVM = nil
		m.inputText = ""
		return m, nil
	case "backspace":
		if m.inputCursor > 0 {
			m.inputText = m.inputText[:m.inputCursor-1] + m.inputText[m.inputCursor:]
			m.inputCursor--
		}
	default:
		// Add character to input
		if len(msg.String()) == 1 {
			m.inputText = m.inputText[:m.inputCursor] + msg.String() + m.inputText[m.inputCursor:]
			m.inputCursor++
		}
	}
	return m, nil
}

// handleVMDetailsUpdate processes keys in the VM details view
func (m Model) handleVMDetailsUpdate(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		// Go back to VM list
		m.state = listVMsView
		m.selectedVM = nil
		return m, nil
	case "s":
		// Start VM
		if m.selectedVM != nil {
			m.loading = true
			return m, startVMCmd(m.selectedVM.UUID)
		}
	case "x":
		// Stop VM
		if m.selectedVM != nil {
			m.loading = true
			return m, stopVMCmd(m.selectedVM.UUID)
		}
	case "p":
		// Suspend (pause) VM
		if m.selectedVM != nil {
			m.loading = true
			return m, suspendVMCmd(m.selectedVM.UUID)
		}
	case "r":
		// Refresh status
		if m.selectedVM != nil {
			m.loading = true
			return m, getStatusCmd(m.selectedVM.UUID)
		}
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
	case "d":
		// Delete selected VM
		if len(m.vms) > 0 {
			m.selectedVM = &m.vms[m.vmsCursor]
			m.state = deleteConfirmView
			m.message = ""
		}
	case "c":
		// Clone selected VM
		if len(m.vms) > 0 {
			m.selectedVM = &m.vms[m.vmsCursor]
			m.state = cloneInputView
			m.inputText = ""
			m.inputCursor = 0
			m.message = ""
		}
	case "enter", "e":
		// Enter VM details view
		if len(m.vms) > 0 {
			m.selectedVM = &m.vms[m.vmsCursor]
			m.state = vmDetailsView
			m.message = ""
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

// deleteVMCmd executes the delete VM command
func deleteVMCmd(uuid string) tea.Cmd {
	return func() tea.Msg {
		err := utm.DeleteVM(uuid)
		if err != nil {
			return actionCompleteMsg{success: false, message: fmt.Sprintf("Delete failed: %v", err)}
		}
		return actionCompleteMsg{success: true, message: "VM deleted successfully"}
	}
}

// cloneVMCmd executes the clone VM command
func cloneVMCmd(uuid, newName string) tea.Cmd {
	return func() tea.Msg {
		err := utm.CloneVM(uuid, newName)
		if err != nil {
			return actionCompleteMsg{success: false, message: fmt.Sprintf("Clone failed: %v", err)}
		}
		return actionCompleteMsg{success: true, message: fmt.Sprintf("VM cloned as '%s'", newName)}
	}
}

// startVMCmd executes the start VM command
func startVMCmd(uuid string) tea.Cmd {
	return func() tea.Msg {
		err := utm.StartVM(uuid)
		if err != nil {
			return actionCompleteMsg{success: false, message: fmt.Sprintf("Start failed: %v", err)}
		}
		return actionCompleteMsg{success: true, message: "VM started successfully"}
	}
}

// stopVMCmd executes the stop VM command
func stopVMCmd(uuid string) tea.Cmd {
	return func() tea.Msg {
		err := utm.StopVM(uuid)
		if err != nil {
			return actionCompleteMsg{success: false, message: fmt.Sprintf("Stop failed: %v", err)}
		}
		return actionCompleteMsg{success: true, message: "VM stopped successfully"}
	}
}

// suspendVMCmd executes the suspend VM command
func suspendVMCmd(uuid string) tea.Cmd {
	return func() tea.Msg {
		err := utm.SuspendVM(uuid)
		if err != nil {
			return actionCompleteMsg{success: false, message: fmt.Sprintf("Suspend failed: %v", err)}
		}
		return actionCompleteMsg{success: true, message: "VM suspended successfully"}
	}
}

// getStatusCmd retrieves the status of a VM
func getStatusCmd(uuid string) tea.Cmd {
	return func() tea.Msg {
		status, err := utm.GetVMStatus(uuid)
		if err != nil {
			return statusUpdateMsg{err: err}
		}
		return statusUpdateMsg{status: status, message: "Status retrieved"}
	}
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

			body.WriteString(footerStyle.Render("\n\n↑/↓: scroll • enter: details • c: clone • d: delete • esc: back to menu"))
		case deleteConfirmView:
			body.WriteString(titleStyle.Render(" DELETE VM ") + "\n\n")
			if m.selectedVM != nil {
				fmt.Fprintf(&body, "Are you sure you want to delete:\n\n")
				fmt.Fprintf(&body, "%s\n\n", selectedItemStyle.Render(m.selectedVM.Name))
				fmt.Fprintf(&body, "Press 'y' to confirm or 'n' to cancel\n")
			}
		case cloneInputView:
			body.WriteString(titleStyle.Render(" CLONE VM ") + "\n\n")
			if m.selectedVM != nil {
				fmt.Fprintf(&body, "Clone from: %s\n\n", selectedItemStyle.Render(m.selectedVM.Name))
				fmt.Fprintf(&body, "New VM name:\n")
				fmt.Fprintf(&body, "%s\n\n", inputStyle.Render(m.inputText+"_"))
				fmt.Fprintf(&body, "Press Enter to clone or Esc to cancel\n")
			}
		case vmDetailsView:
			if m.selectedVM != nil {
				body.WriteString(titleStyle.Render(" VM DETAILS ") + "\n\n")
				fmt.Fprintf(&body, "Name:   %s\n", selectedItemStyle.Render(m.selectedVM.Name))
				fmt.Fprintf(&body, "UUID:   %s\n", m.selectedVM.UUID)
				fmt.Fprintf(&body, "Status: %s\n\n", m.selectedVM.Status)

				if m.loading {
					body.WriteString(processingStyle.Render("Processing action... Please wait.") + "\n")
				} else {
					body.WriteString("Actions:\n")
					body.WriteString("  s - Start VM\n")
					body.WriteString("  x - Stop VM\n")
					body.WriteString("  p - Suspend VM\n")
					body.WriteString("  r - Refresh Status\n\n")
				}

				if m.message != "" {
					style := errorStyle
					if m.isSuccess {
						style = successStyle
					}
					fmt.Fprintf(&body, "%s\n", style.Render(m.message))
				}
			}
			body.WriteString(footerStyle.Render("\n\nEsc: back to list"))
		}
	}

	// Center the content
	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		containerStyle.Render(body.String()),
	)
}
