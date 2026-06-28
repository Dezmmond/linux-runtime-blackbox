package run

import (
	"testing"
)

func TestCommandCapturesExitCode(t *testing.T) {
	report, err := Command([]string{"sh", "-c", "exit 7"})
	if err != nil {
		t.Fatalf("Command returned error: %v", err)
	}
	if report.Run == nil {
		t.Fatal("Run is nil")
	}
	if report.Run.ExitCode != 7 {
		t.Fatalf("exit code = %d, want 7", report.Run.ExitCode)
	}
	if report.Run.PID <= 0 {
		t.Fatalf("pid = %d, want positive", report.Run.PID)
	}
	if len(report.Run.Command) != 3 {
		t.Fatalf("command len = %d, want 3", len(report.Run.Command))
	}
	if report.Run.StartedAt == "" || report.Run.EndedAt == "" {
		t.Fatalf("started_at/ended_at should be set: %#v", report.Run)
	}
}

func TestCommandCollectsInitialSnapshot(t *testing.T) {
	report, err := Command([]string{"sh", "-c", "sleep 0.2"})
	if err != nil {
		t.Fatalf("Command returned error: %v", err)
	}
	if report.InitialSnapshot == nil {
		t.Fatal("InitialSnapshot is nil")
	}
	if report.Target.PID != report.Run.PID {
		t.Fatalf("target pid = %d, run pid = %d", report.Target.PID, report.Run.PID)
	}
}
