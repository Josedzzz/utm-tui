// Package tui handles the screens
package tui

import "github.com/Josedzzz/utm-tui/internal/utm"

// sessionState represents the different screens of the TUI
type sessionState int

const (
	menuView sessionState = iota
	listVMsView
	deleteConfirmView
	cloneInputView
	vmDetailsView
)

// Model holds the entire TUI state
type Model struct {
	state     sessionState
	choices   []string // menu options
	cursor    int      // menu cursor
	quitting  bool
	width     int
	height    int
	message   string
	isSuccess bool

	// VM list view
	vms       []utm.VM
	vmsCursor int
	vmsOffset int
	loading   bool

	// Action-specific fields
	selectedVM  *utm.VM // VM selected for delete or clone
	inputText   string  // Text input for clone name
	inputCursor int     // Cursor position in input
}

// NewModel creates and returns a new model with default values
func NewModel() Model {
	return Model{
		choices:   []string{"List VMs", "Exit"},
		state:     menuView,
		cursor:    0,
		vms:       []utm.VM{},
		loading:   false,
		message:   "",
		isSuccess: true,
	}
}
