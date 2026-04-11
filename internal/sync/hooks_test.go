package sync

import (
	"runtime"
	"testing"
)

func TestHookRunner_Run_NoMatchingHooks(t *testing.T) {
	r := NewHookRunner([]Hook{
		{Event: HookPreSync, Command: "echo pre"},
	})
	// Running post-sync should not execute the pre-sync hook and return nil.
	if err := r.Run(HookPostSync); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestHookRunner_Run_SuccessfulHook(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping shell command test on Windows")
	}
	r := NewHookRunner([]Hook{
		{Event: HookPostSync, Command: "echo hello"},
	})
	if err := r.Run(HookPostSync); err != nil {
		t.Fatalf("expected nil error, got: %v", err)
	}
}

func TestHookRunner_Run_FailingHook(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping shell command test on Windows")
	}
	r := NewHookRunner([]Hook{
		{Event: HookPreSync, Command: "false"},
	})
	err := r.Run(HookPreSync)
	if err == nil {
		t.Fatal("expected an error from failing hook, got nil")
	}
}

func TestHookRunner_Run_MultipleHooks_CollectsErrors(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skipping shell command test on Windows")
	}
	r := NewHookRunner([]Hook{
		{Event: HookPostSync, Command: "false"},
		{Event: HookPostSync, Command: "false"},
	})
	err := r.Run(HookPostSync)
	if err == nil {
		t.Fatal("expected combined error, got nil")
	}
}

func TestHookRunner_Run_EmptyCommand_ReturnsError(t *testing.T) {
	r := NewHookRunner([]Hook{
		{Event: HookPreSync, Command: ""},
	})
	err := r.Run(HookPreSync)
	if err == nil {
		t.Fatal("expected error for empty command, got nil")
	}
}

func TestHookEvent_Constants(t *testing.T) {
	if HookPreSync != "pre-sync" {
		t.Errorf("unexpected HookPreSync value: %s", HookPreSync)
	}
	if HookPostSync != "post-sync" {
		t.Errorf("unexpected HookPostSync value: %s", HookPostSync)
	}
}
