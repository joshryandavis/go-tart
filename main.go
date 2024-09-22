// Package tart provides a Go API for interacting with Tart.
// It allows you to create, clone, run, and manage virtual machines.
package tart

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
)

// Tart represents the Tart hypervisor.
type Tart struct {
	ConfigDir string `json:"configDir"`
	Host      string `json:"host"`
}

// New creates a new Tart instance with a custom config directory.
// It returns an error if the 'tart' command is not found in the system PATt.
func New() (*Tart, error) {
	// Validate that TART is on the path
	_, err := exec.LookPath("tart")
	if err != nil {
		return nil, errors.New("tart command not found in PATH")
	}
	// Configure the config directory
	configDir := getConfigDir()
	return &Tart{
		ConfigDir: configDir,
	}, nil
}

// setTartHome sets the TART_HOME environment variable for the given command.
// It returns an error if the specified config directory does not exist.
func (t *Tart) setTartHome(cmd *exec.Cmd) error {
	if t.ConfigDir != "" {
		if _, err := os.Stat(t.ConfigDir); os.IsNotExist(err) {
			return errors.New("config directory does not exist")
		}
		cmd.Env = append(os.Environ(), fmt.Sprintf("TART_HOME=%s", t.ConfigDir))
	}
	return nil
}

// run is a helper function to execute Tart commands
func (t *Tart) run(args ...string) ([]byte, error) {
	cmd := exec.Command("tart", args...)
	t.setTartHome(cmd)

	// Create pipes for stdout and stderr
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start command: %w", err)
	}

	// Read stdout and stderr
	stdout, err := io.ReadAll(stdoutPipe)
	if err != nil {
		return nil, fmt.Errorf("failed to read stdout: %w", err)
	}
	stderr, err := io.ReadAll(stderrPipe)
	if err != nil {
		return nil, fmt.Errorf("failed to read stderr: %w", err)
	}

	// Wait for the command to finish
	if err := cmd.Wait(); err != nil {
		return nil, fmt.Errorf("command failed: %w, stderr: %s", err, stderr)
	}

	return stdout, nil
}

// Returns the directory where we store our configuration
func getConfigDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	configDir := filepath.Join(homeDir, ".tart")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		os.Mkdir(configDir, 0700)
	}
	return configDir
}
