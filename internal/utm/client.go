// Package utm provides logic to handle the vms
package utm

import (
	"fmt"
	"os/exec"
	"strings"
)

// VM represents a single virtual machine
type VM struct {
	UUID   string
	Name   string
	Status string
}

// ListVMs fetches the list of virtual machines from utmctl
func ListVMs() ([]VM, error) {
	cmd := exec.Command("utmctl", "list")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to run utmctl: %w", err)
	}

	lines := strings.Split(string(out), "\n")
	var vms []VM
	for i, line := range lines {
		if i == 0 || strings.TrimSpace(line) == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 3 {
			vms = append(vms, VM{
				UUID:   fields[0],
				Status: fields[1],
				Name:   fields[2],
			})
		}
	}
	return vms, nil
}
