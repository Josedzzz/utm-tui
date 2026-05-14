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

// CloneVM clones an existing virtual machine with a new name
func CloneVM(uuid, newName string) error {
	cmd := exec.Command("utmctl", "clone", uuid, "--name", newName)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to clone VM: %w", err)
	}
	return nil
}

// DeleteVM deletes a virtual machine by UUID
func DeleteVM(uuid string) error {
	cmd := exec.Command("utmctl", "delete", uuid)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to delete VM: %w", err)
	}
	return nil
}
