package tart

import "fmt"

// VMConfig represents the parameters of a VM.
type VMConfig struct {
	Version       int    `json:"version"`
	OS            string `json:"os"`
	Arch          string `json:"arch"`
	CPUCountMin   int    `json:"cpuCountMin"`
	CPUCount      int    `json:"cpuCount"`
	MemorySizeMin uint64 `json:"memorySizeMin"`
	MemorySize    uint64 `json:"memorySize"`
	MACAddress    string `json:"macAddress"`
	Display       struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"display"`
}

// SetConfig modifies a VM's configuration.
// It returns an error if the configuration update process fails.
func (t *Tart) SetConfig(name string, config VMConfig) error {
	args := []string{"set", name}
	if config.CPUCount > 0 {
		args = append(args, "--cpu", fmt.Sprintf("%d", config.CPUCount))
	}
	if config.MemorySize > 0 {
		args = append(args, "--memory", fmt.Sprintf("%d", config.MemorySize))
	}
	if config.Display.Width > 0 && config.Display.Height > 0 {
		args = append(args, "--display", fmt.Sprintf("%dx%d", config.Display.Width, config.Display.Height))
	}
	if config.MACAddress == "random" {
		args = append(args, "--random-mac")
	}
	output, err := t.run(args...)
	if err != nil {
		return fmt.Errorf("failed to set VM configuration: %w, output: %s", err, string(output))
	}
	return nil
}

// GetConfig retrieves a VM's configuration.
// It returns the configuration as a string and an error if the retrieval process fails.
func (t *Tart) GetConfig(name string, format string) (string, error) {
	args := []string{"get", name}
	if format != "" {
		args = append(args, "--format", format)
	}
	output, err := t.run(args...)
	if err != nil {
		return "", fmt.Errorf("failed to get VM configuration: %w, output: %s", err, string(output))
	}
	return string(output), nil
}

// Rename renames a local VM.
// It returns an error if the rename process fails.
func (t *Tart) Rename(oldName string, newName string) error {
	output, err := t.run("rename", oldName, newName)
	if err != nil {
		return fmt.Errorf("failed to rename VM: %w, output: %s", err, string(output))
	}
	return nil
}

// CreateOptions represents the configuration for creating a new VM.
type CreateOptions struct {
	FromIPSW string `json:"fromIPSW"`
	Linux    bool   `json:"linux"`
	DiskSize int    `json:"diskSize"`
}

// Create creates a new VM and returns it.
// It returns an error if a VM with the same name already exists or if the creation process fails.
func (t *Tart) Create(name string, options CreateOptions) error {
	// Check if the VM name is already taken
	localVMs, err := t.List(ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list local VMs: %w", err)
	}
	for _, existingVM := range localVMs {
		if existingVM.Name == name {
			return fmt.Errorf("VM with name %s already exists", name)
		}
	}
	args := []string{"create", name}
	if options.FromIPSW != "" {
		args = append(args, "--from-ipsw", options.FromIPSW)
	}
	if options.Linux {
		args = append(args, "--linux")
	}
	if options.DiskSize > 0 {
		args = append(args, "--disk-size", fmt.Sprintf("%d", options.DiskSize))
	}
	output, err := t.run(args...)
	if err != nil {
		return fmt.Errorf("failed to create VM: %w, output: %s", err, string(output))
	}
	return nil
}

// CloneOptions represents the configuration for cloning a VM.
type CloneOptions struct {
	NewName     string `json:"newName"`
	Insecure    bool   `json:"insecure"`
	Concurrency int    `json:"concurrency"`
}

// Clone clones an existing VM.
// It returns an error if a VM with the new name already exists or if the cloning process fails.
func (t *Tart) Clone(sourceName string, newName string, options CloneOptions) error {
	// Check if the new VM name is already taken
	localVMs, err := t.List(ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list local VMs: %w", err)
	}
	for _, existingVM := range localVMs {
		if existingVM.Name == newName {
			return fmt.Errorf("VM with name %s already exists", newName)
		}
	}
	args := []string{"clone", sourceName, newName}
	if options.Insecure {
		args = append(args, "--insecure")
	}
	if options.Concurrency > 0 {
		args = append(args, "--concurrency", fmt.Sprintf("%d", options.Concurrency))
	}
	output, err := t.run(args...)
	if err != nil {
		return fmt.Errorf("failed to clone VM: %w, output: %s", err, string(output))
	}
	return nil
}

// ImportOptions represents the configuration for importing an IPSW.
type ImportOptions struct {
	Path string `json:"path"`
}

// Import imports a VM from a compressed .tvm file.
// It returns an error if a VM with the same name already exists or if the import process fails.
func (t *Tart) Import(path string, name string) error {
	// Check if the VM name is already taken
	localVMs, err := t.List(ListOptions{})
	if err != nil {
		return fmt.Errorf("failed to list local VMs: %w", err)
	}
	for _, existingVM := range localVMs {
		if existingVM.Name == name {
			return fmt.Errorf("VM with name %s already exists", name)
		}
	}
	output, err := t.run("import", path, name)
	if err != nil {
		return fmt.Errorf("failed to import VM: %w, output: %s", err, string(output))
	}
	return nil
}

// Export exports a VM to a compressed .tvm file.
// It returns an error if the export process fails.
func (t *Tart) Export(name string, path string) error {
	args := []string{"export", name}
	if path != "" {
		args = append(args, path)
	}
	output, err := t.run(args...)
	if err != nil {
		return fmt.Errorf("failed to export VM: %w, output: %s", err, string(output))
	}
	return nil
}

// Suspend suspends a VM.
// It returns an error if the suspension process fails.
func (t *Tart) Suspend(name string) error {
	output, err := t.run("suspend", name)
	if err != nil {
		return fmt.Errorf("failed to suspend VM: %w, output: %s", err, string(output))
	}
	return nil
}

// Stop stops a VM.
// It returns an error if the stop process fails.
func (t *Tart) Stop(name string, timeout int) error {
	args := []string{"stop", name}
	if timeout > 0 {
		args = append(args, "--timeout", fmt.Sprintf("%d", timeout))
	}
	output, err := t.run(args...)
	if err != nil {
		return fmt.Errorf("failed to stop VM: %w, output: %s", err, string(output))
	}
	return nil
}

// Delete deletes a VM.
// It returns an error if the deletion process fails.
func (t *Tart) Delete(name string) error {
	output, err := t.run("delete", name)
	if err != nil {
		return fmt.Errorf("failed to delete VM: %w, output: %s", err, string(output))
	}
	return nil
}

// PruneOptions represents the options for pruning.
type PruneOptions struct {
	Entries     string `json:"entries"`
	OlderThan   int    `json:"olderThan"`
	SpaceBudget int    `json:"spaceBudget"`
}

// Prune prunes OCI and IPSW caches or local VMs.
// It returns an error if the pruning process fails.
func (t *Tart) Prune(options PruneOptions) error {
	args := []string{"prune"}
	if options.Entries != "" {
		args = append(args, "--entries", options.Entries)
	}
	if options.OlderThan > 0 {
		args = append(args, "--older-than", fmt.Sprintf("%d", options.OlderThan))
	}
	if options.SpaceBudget > 0 {
		args = append(args, "--space-budget", fmt.Sprintf("%d", options.SpaceBudget))
	}
	output, err := t.run(args...)
	if err != nil {
		return fmt.Errorf("failed to prune: %w, output: %s", err, string(output))
	}
	return nil
}
