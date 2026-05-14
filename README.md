# UTM TUI Dashboard

A terminal user interface (TUI) dashboard for managing and controlling UTM virtual machines directly from your command line.

## Overview

**utm-tui** is a lightweight, keyboard-driven dashboard built with Go that provides an intuitive way to manage your UTM (Universal Turing Machine / UTM hypervisor) virtual machines. Instead of navigating through the GUI, you can now control your VMs efficiently from the terminal.

## Features

### VM Management

- **List VMs** - View all your registered virtual machines with their current status
- **VM Details** - Access detailed information about individual VMs
- **Status Monitoring** - Real-time status checking (running, stopped, suspended)

### VM Control Actions

- **Start VM** - Boot up a stopped or suspended virtual machine
- **Stop VM** - Gracefully shut down a running VM
- **Suspend VM** - Pause a running VM to memory
- **Refresh Status** - Check current VM status on demand

### VM Operations

- **Clone VM** - Create a copy of an existing VM with a new name
- **Delete VM** - Remove a VM from your system (with confirmation)

## Installation

### Requirements

- Go 1.26.3 or higher
- UTM hypervisor installed
- `utmctl` CLI tool in your PATH

### Build

```bash
git clone https://github.com/Josedzzz/utm-tui.git
cd utm-tui
go build -o utm-tui ./cmd/main.go
```

### Run

```bash
./utm-tui
```

## Usage

### Main Menu

```
┌─ MAIN MENU ──────────┐
│ > List VMs           │
│   Exit               │
└──────────────────────┘
```

- `↑/↓` or `j/k` - Navigate menu
- `Enter` - Select option
- `Ctrl+C` or `q` - Quit

### VM List View

```
NAME                           STATUS
--------------------------------------------------
Ubuntu Desktop                 running
Debian Server                  stopped
Windows 11                     suspended
```

**Keyboard Shortcuts:**

- `↑/↓` or `j/k` - Scroll through VMs
- `Enter` - View VM details and controls
- `c` - Clone selected VM
- `d` - Delete selected VM
- `Esc` - Return to main menu
- `Ctrl+C` - Quit

### VM Details View

```
┌─ VM DETAILS ─────────┐
│ Name:   Ubuntu 22.04 │
│ UUID:   ABC123...    │
│ Status: running      │
│                      │
│ Actions:             │
│   s - Start VM       │
│   x - Stop VM        │
│   p - Suspend VM     │
│   r - Refresh Status │
└──────────────────────┘
```

**Keyboard Shortcuts:**

- `s` - Start VM
- `x` - Stop VM
- `p` - Suspend (pause) VM
- `r` - Refresh current status
- `Esc` - Return to VM list

### Clone VM

```
┌─ CLONE VM ──────────────┐
│ Clone from: Ubuntu 22.04│
│                         │
│ New VM name:            │
│ ┌────────────────────┐  │
│ │ Ubuntu-Copy-01_    │  │
│ └────────────────────┘  │
│                         │
│ Press Enter to clone or │
│ Esc to cancel           │
└─────────────────────────┘
```

- Type the name for the cloned VM
- `Enter` - Confirm clone
- `Esc` - Cancel

### Delete VM

```
┌─ DELETE VM ──────┐
│ Are you sure you │
│ want to delete:   │
│                   │
│ Old VM Name       │
│                   │
│ Press 'y' to      │
│ confirm or 'n' to │
│ cancel            │
└──────────────────┘
```

- `y` or `Enter` - Confirm deletion
- `n` or `Esc` - Cancel

## Project Structure

```
utm-tui/
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── tui/                    # Terminal UI logic
│   │   ├── model.go            # State management
│   │   ├── view.go             # UI rendering and handlers
│   │   └── styles.go           # Terminal styling
│   └── utm/                    # UTM client integration
│       └── client.go           # UTM command wrappers
├── go.mod                      # Go module definition
├── go.sum                      # Dependency checksums
└── README.md                   # This file
```

## Architecture

The application follows the **Model-View-Update (MVU)** pattern using the [Bubbletea](https://github.com/charmbracelet/bubbletea) TUI framework:

- **Model** (`model.go`) - Holds application state (current screen, VM list, input fields)
- **Update** (`view.go`) - Processes messages and updates state based on user input
- **View** (`view.go`) - Renders the current state to the terminal
- **Styling** (`styles.go`) - Defines colors and text formatting

## Dependencies

- **[Bubbletea](https://github.com/charmbracelet/bubbletea)** - TUI framework with event loop and state management
- **[Lipgloss](https://github.com/charmbracelet/lipgloss)** - Terminal styling and layout library

## Keyboard Quick Reference

| Context    | Key         | Action         |
| ---------- | ----------- | -------------- |
| Main Menu  | `↑/↓` `j/k` | Navigate       |
| Main Menu  | `Enter`     | Select         |
| VM List    | `↑/↓` `j/k` | Scroll         |
| VM List    | `Enter`     | View details   |
| VM List    | `c`         | Clone VM       |
| VM List    | `d`         | Delete VM      |
| VM Details | `s`         | Start VM       |
| VM Details | `x`         | Stop VM        |
| VM Details | `p`         | Suspend VM     |
| VM Details | `r`         | Refresh status |
| Anywhere   | `Esc`       | Go back        |
| Anywhere   | `Ctrl+C`    | Quit           |

## Features Coming Soon

- VM resource monitoring (CPU, RAM, disk usage)
- VM settings modification (network, resources)
- Batch operations on multiple VMs
- VM snapshots management
- SSH connection to VMs
- Configuration file support

## Requirements

- macOS with UTM installed
- `utmctl` command-line tool available in PATH
- Terminal with support for 256 colors (recommended)

## Development

### Build

```bash
go build -o utm-tui ./cmd/main.go
```

### Run with Debug Output

```bash
go run ./cmd/main.go
```

### Code Quality

```bash
go vet ./...
go fmt ./...
```

## License

MIT License - See LICENSE file for details

## Contributing

Contributions are welcome! Feel free to open issues or pull requests to improve the project.

## Support

For issues or feature requests, please visit: https://github.com/Josedzzz/utm-tui

## Acknowledgments

- Built with [Bubbletea](https://github.com/charmbracelet/bubbletea) and [Lipgloss](https://github.com/charmbracelet/lipgloss)
- Inspired by tools like `lazydocker` and `lazygit`

---

**Made with <3 for efficient VM management**
