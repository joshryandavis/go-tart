package tart

import "fmt"

// LoginOptions represents options for logging in to a registry.
type LoginOptions struct {
	Username      string `json:"username"`
	PasswordStdin bool   `json:"password_stdin"`
	Insecure      bool   `json:"insecure"`
	NoValidate    bool   `json:"no_validate"`
}

// Login logs in to a registry.
//
// It takes a LoginOptions struct as parameter.
//
// It returns an error if the login process fails.
func (t *Tart) Login(opts LoginOptions) error {
	args := []string{"login", t.Host}
	if opts.Username != "" {
		args = append(args, "--username", opts.Username)
	}
	if opts.PasswordStdin {
		args = append(args, "--password-stdin")
	}
	if opts.Insecure {
		args = append(args, "--insecure")
	}
	if opts.NoValidate {
		args = append(args, "--no-validate")
	}
	output, err := t.run(args...)
	if err != nil {
		return fmt.Errorf("failed to login: %w, output: %s", err, string(output))
	}
	return nil
}

// Logout logs out from a registry.
//
// It returns an error if the logout process fails.
func (t *Tart) Logout() error {
	output, err := t.run("logout", t.Host)
	if err != nil {
		return fmt.Errorf("failed to logout: %w, output: %s", err, string(output))
	}
	return nil
}

// PushOptions represents the options for pushing a VM to a registry.
type PushOptions struct {
	RemoteNames   []string `json:"remoteNames"`
	Insecure      bool     `json:"insecure"`
	Concurrency   int      `json:"concurrency"`
	ChunkSize     int      `json:"chunkSize"`
	PopulateCache bool     `json:"populateCache"`
}

// Push pushes a VM to a registry.
// It returns an error if the push process fails.
func (t *Tart) Push(name string, options PushOptions) error {
	args := []string{"push", name}
	args = append(args, options.RemoteNames...)
	if options.Insecure {
		args = append(args, "--insecure")
	}
	if options.Concurrency > 0 {
		args = append(args, "--concurrency", fmt.Sprintf("%d", options.Concurrency))
	}
	if options.ChunkSize > 0 {
		args = append(args, "--chunk-size", fmt.Sprintf("%d", options.ChunkSize))
	}
	if options.PopulateCache {
		args = append(args, "--populate-cache")
	}
	output, err := t.run(args...)
	if err != nil {
		return fmt.Errorf("failed to push VM: %w, output: %s", err, string(output))
	}
	return nil
}

// Pull pulls a VM from a registry.
// It returns an error if the pull process fails.
func (t *Tart) Pull(name string, insecure bool, concurrency int) error {
	args := []string{"pull", name}
	if insecure {
		args = append(args, "--insecure")
	}
	if concurrency > 0 {
		args = append(args, "--concurrency", fmt.Sprintf("%d", concurrency))
	}
	output, err := t.run(args...)
	if err != nil {
		return fmt.Errorf("failed to pull VM: %w, output: %s", err, string(output))
	}
	return nil
}
