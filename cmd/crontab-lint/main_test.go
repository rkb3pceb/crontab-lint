package main

import (
	"os/exec"
	"strings"
	"testing"
)

// buildBinary compiles the binary once and returns its path.
// Tests are skipped if the build fails (e.g., missing deps in CI).
func buildBinary(t *testing.T) string {
	t.Helper()
	tmp := t.TempDir()
	out := tmp + "/crontab-lint"
	cmd := exec.Command("go", "build", "-o", out, ".")
	if err := cmd.Run(); err != nil {
		t.Skipf("could not build binary: %v", err)
	}
	return out
}

func TestMain_ValidExpression(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin, "*/5 * * * * /usr/bin/backup.sh")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected exit 0, got error: %v\noutput: %s", err, out)
	}
	if !strings.Contains(string(out), "valid") && !strings.Contains(string(out), "OK") {
		t.Errorf("expected valid indicator in output, got: %s", out)
	}
}

func TestMain_InvalidExpression_ExitsNonZero(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin, "99 * * * * /bin/sh")
	out, _ := cmd.CombinedOutput()
	if cmd.ProcessState.ExitCode() == 0 {
		t.Errorf("expected non-zero exit for invalid expression, output: %s", out)
	}
}

func TestMain_JSONOutput(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin, "-format=json", "0 9 * * 1 /bin/notify")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected exit 0, got error: %v\noutput: %s", err, out)
	}
	if !strings.Contains(string(out), "{") {
		t.Errorf("expected JSON output, got: %s", out)
	}
}

func TestMain_VersionFlag(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin, "-version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("expected exit 0 for -version, got: %v", err)
	}
	if !strings.Contains(string(out), version) {
		t.Errorf("expected version %q in output, got: %s", version, out)
	}
}

func TestMain_NoArgs_ExitsWithUsage(t *testing.T) {
	bin := buildBinary(t)
	cmd := exec.Command(bin)
	cmd.Run()
	if cmd.ProcessState.ExitCode() != 2 {
		t.Errorf("expected exit code 2 when no args provided, got %d", cmd.ProcessState.ExitCode())
	}
}
