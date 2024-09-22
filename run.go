package tart

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"
)

// DirMount represents a directory mount with its options
type DirMount struct {
	Name     string `json:"name"`
	Path     string `json:"path"`
	ReadOnly bool   `json:"readOnly"`
	Sync     string `json:"sync"`
	Tag      string `json:"tag"`
}

// RunOptions represents the options for running a VM.
type RunOptions struct {
	NoGraphics        bool       `json:"noGraphics"`
	Serial            bool       `json:"serial"`
	SerialPath        string     `json:"serialPath"`
	NoAudio           bool       `json:"noAudio"`
	NoClipboard       bool       `json:"noClipboard"`
	Recovery          bool       `json:"recovery"`
	VNC               bool       `json:"vnc"`
	VNCExperimental   bool       `json:"vncExperimental"`
	Disk              []string   `json:"disk"`
	Rosetta           string     `json:"rosetta"`
	Dir               []DirMount `json:"dir"`
	NetBridged        string     `json:"netBridged"`
	NetSoftnet        bool       `json:"netSoftnet"`
	NetSoftnetAllow   string     `json:"netSoftnetAllow"`
	NetHost           bool       `json:"netHost"`
	RootDiskOpts      string     `json:"rootDiskOpts"`
	Suspendable       bool       `json:"suspendable"`
	CaptureSystemKeys bool       `json:"captureSystemKeys"`
}

// Run runs a VM with the specified options.
// It returns an error if the VM is already running, doesn't exist, or if the run process fails.
func (t *Tart) Run(name string, options RunOptions) error {
	s, err := t.State(name)
	if err != nil {
		return fmt.Errorf("failed to get VM state: %w", err)
	}
	if s.State == "running" {
		return fmt.Errorf("VM is already running")
	}
	if s.Name != name {
		return fmt.Errorf("VM with name %s does not exist", name)
	}
	args := []string{"run"}
	if options.NoGraphics {
		args = append(args, "--no-graphics")
	}
	if options.Serial {
		args = append(args, "--serial")
	}
	if options.SerialPath != "" {
		args = append(args, "--serial-path", options.SerialPath)
	}
	if options.NoAudio {
		args = append(args, "--no-audio")
	}
	if options.NoClipboard {
		args = append(args, "--no-clipboard")
	}
	if options.Recovery {
		args = append(args, "--recovery")
	}
	if options.VNC {
		args = append(args, "--vnc")
	}
	if options.VNCExperimental {
		args = append(args, "--vnc-experimental")
	}
	for _, disk := range options.Disk {
		args = append(args, "--disk", disk)
	}
	if options.Rosetta != "" {
		args = append(args, "--rosetta", options.Rosetta)
	}
	for _, dir := range options.Dir {
		dirArg := ""
		if dir.Name != "" {
			dirArg += dir.Name + ":"
		}
		dirArg += dir.Path
		if dir.ReadOnly || dir.Tag != "" {
			dirArg += ":"
			if dir.ReadOnly {
				dirArg += "ro"
				if dir.Tag != "" {
					dirArg += ","
				}
			}
			if dir.Tag != "" {
				dirArg += "tag=" + dir.Tag
			}
			if dir.Sync != "" {
				dirArg += "sync=" + dir.Sync
			}
		}
		args = append(args, "--dir", dirArg)
	}
	if options.NetBridged != "" {
		args = append(args, "--net-bridged", options.NetBridged)
	}
	if options.NetSoftnet {
		args = append(args, "--net-softnet")
	}
	if options.NetSoftnetAllow != "" {
		args = append(args, "--net-softnet-allow", options.NetSoftnetAllow)
	}
	if options.NetHost {
		args = append(args, "--net-host")
	}
	if options.RootDiskOpts != "" {
		args = append(args, "--root-disk-opts", options.RootDiskOpts)
	}
	if options.Suspendable {
		args = append(args, "--suspendable")
	}
	if options.CaptureSystemKeys {
		args = append(args, "--capture-system-keys")
	}
	args = append(args, name)

	cmd := exec.Command("tart", args...)
	t.setTartHome(cmd)

	serialOut, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start VM: %w", err)
	}

	reader := bufio.NewReader(serialOut)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return fmt.Errorf("failed to read serial output: %w", err)
		}
		if strings.Contains(line, "VM is up") {
			fmt.Println("VM is up and running")
			return nil
		}
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("VM process exited with error: %w", err)
	}

	return nil
}
