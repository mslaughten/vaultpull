package sync

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// HookEvent represents the lifecycle point at which a hook runs.
type HookEvent string

const (
	HookPreSync  HookEvent = "pre-sync"
	HookPostSync HookEvent = "post-sync"
)

// Hook defines a shell command to run at a specific lifecycle event.
type Hook struct {
	Event   HookEvent
	Command string
}

// HookRunner executes lifecycle hooks.
type HookRunner struct {
	hooks  []Hook
	stdout *os.File
	stderr *os.File
}

// NewHookRunner creates a HookRunner that writes output to the provided writers.
func NewHookRunner(hooks []Hook) *HookRunner {
	return &HookRunner{
		hooks:  hooks,
		stdout: os.Stdout,
		stderr: os.Stderr,
	}
}

// Run executes all hooks registered for the given event.
// All hooks are attempted; errors are collected and returned together.
func (r *HookRunner) Run(event HookEvent) error {
	var errs []string
	for _, h := range r.hooks {
		if h.Event != event {
			continue
		}
		if err := r.exec(h.Command); err != nil {
			errs = append(errs, fmt.Sprintf("hook %q failed: %v", h.Command, err))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("%s hooks had errors:\n  %s", event, strings.Join(errs, "\n  "))
	}
	return nil
}

func (r *HookRunner) exec(command string) error {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return fmt.Errorf("empty command")
	}
	cmd := exec.Command(parts[0], parts[1:]...) //nolint:gosec
	cmd.Stdout = r.stdout
	cmd.Stderr = r.stderr
	return cmd.Run()
}
