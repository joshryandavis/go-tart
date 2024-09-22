package tart

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Constants representing the options for a VM source.
const (
	SourceLocal  = "local"
	SourceRemote = "remote"
)

// ListOptions represents the options for listing VMs.
type ListOptions struct {
	Source *string `json:"source,omitempty"`
}

// VMState represents the state of a VM.
type VMState struct {
	SizeOnDisk int    `json:"sizeOnDisk"`
	Disk       int    `json:"disk"`
	Name       string `json:"name"`
	Source     string `json:"source"`
	Size       int    `json:"size"`
	State      string `json:"state"`
}

// UnmarshalJSON implements the json.Unmarshaler interface for VMState.
func (v *VMState) UnmarshalJSON(data []byte) error {
	type Alias VMState
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(v),
	}
	return json.Unmarshal(data, &aux)
}

// List lists VMs.
// It returns a slice of VMState and an error if the listing process fails.
func (t *Tart) List(config ListOptions) ([]VMState, error) {
	// source can be empty but if not should be either "local" or "remote"
	if config.Source != nil && *config.Source != SourceLocal && *config.Source != SourceRemote {
		return nil, fmt.Errorf("invalid source: %s", *config.Source)
	}
	args := []string{"list", "--format", "json"}
	if config.Source != nil {
		args = append(args, "--source", *config.Source)
	}
	output, err := t.run(args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list VMs: %w, output: %s", err, string(output))
	}
	var vms []VMState
	err = json.Unmarshal(output, &vms)
	if err != nil {
		return nil, err
	}
	return vms, nil
}

// State gets the state of a VM.
// It returns a VMState struct and an error if the state retrieval process fails.
func (t *Tart) State(name string) (VMState, error) {
	var ret VMState
	vms, err := t.List(ListOptions{})
	if err != nil {
		return ret, fmt.Errorf("failed to get VM state: %w", err)
	}
	for _, s := range vms {
		if s.Name == name {
			ret = s
			break
		}
	}
	return ret, nil
}

// IP retrieves a VM's IP address.
// It returns the IP address as a string and an error if the retrieval process fails.
func (t *Tart) IP(name string, wait int, resolver string) (string, error) {
	args := []string{"ip", name}
	if wait > 0 {
		args = append(args, "--wait", fmt.Sprintf("%d", wait))
	}
	if resolver != "" {
		args = append(args, "--resolver", resolver)
	}
	output, err := t.run(args...)
	if err != nil {
		return "", fmt.Errorf("failed to get VM IP: %w, output: %s", err, string(output))
	}
	return strings.TrimSpace(string(output)), nil
}

// Exists checks if a VM exists
func (t *Tart) Exists(name string) (bool, error) {
	localVMs, err := t.List(ListOptions{})
	if err != nil {
		return false, fmt.Errorf("failed to list local VMs: %w", err)
	}
	for _, existingVM := range localVMs {
		if existingVM.Name == name {
			return true, nil
		}
	}
	return false, nil
}

// Running checks if a VM is running.
func (t *Tart) Running(name string) (bool, error) {
	s, err := t.State(name)
	if err != nil {
		return false, fmt.Errorf("failed to get VM state: %w", err)
	}
	return s.State == "running", nil
}

// Stopped checks if a VM is stopped.
func (t *Tart) Stopped(name string) (bool, error) {
	s, err := t.State(name)
	if err != nil {
		return false, fmt.Errorf("failed to get VM state: %w", err)
	}
	return s.State == "stopped", nil
}

// Suspended checks if a VM is suspended.
func (t *Tart) Suspended(name string) (bool, error) {
	s, err := t.State(name)
	if err != nil {
		return false, fmt.Errorf("failed to get VM state: %w", err)
	}
	return s.State == "suspended", nil
}
