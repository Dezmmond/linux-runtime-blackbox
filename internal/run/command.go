package run

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/Dezmmond/linux-runtime-blackbox/internal/collector"
	"github.com/Dezmmond/linux-runtime-blackbox/internal/model"
)

type Options struct {
	OutputPath string
}

func Command(args []string) (model.Report, error) {
	if len(args) == 0 {
		return model.Report{}, fmt.Errorf("missing command")
	}

	started := time.Now().UTC()
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return model.Report{}, fmt.Errorf("start command: %w", err)
	}

	pid := cmd.Process.Pid
	r := model.Report{
		SchemaVersion: collector.SchemaVersion,
		CollectedAt:   started.Format(time.RFC3339Nano),
		Run: &model.RunInfo{
			Command:   append([]string(nil), args...),
			PID:       pid,
			StartedAt: started.Format(time.RFC3339Nano),
			ExitCode:  -1,
		},
		Target: model.Target{
			PID:     pid,
			Cmdline: append([]string(nil), args...),
		},
		Process: model.Process{
			PID: pid,
		},
		Namespaces: map[string]string{},
		Warnings:   []model.Warning{},
	}

	time.Sleep(25 * time.Millisecond)
	initial, err := collector.CollectPID(pid)
	if err != nil {
		r.Warnings = append(r.Warnings, collectWarning("initial_snapshot_failed", err))
	} else {
		initial.Run = nil
		initial.InitialSnapshot = nil
		initial.FinalSnapshot = nil
		r.InitialSnapshot = &initial
		copySnapshotFields(&r, initial)
	}

	waitErr := cmd.Wait()
	ended := time.Now().UTC()
	r.CollectedAt = ended.Format(time.RFC3339Nano)
	r.Run.EndedAt = ended.Format(time.RFC3339Nano)
	r.Run.DurationMillis = ended.Sub(started).Milliseconds()
	r.Run.ExitCode, r.Run.Signal = exitStatus(waitErr)

	final, err := collector.CollectPID(pid)
	if err != nil {
		if errors.Is(err, collector.ErrProcessNotFound) {
			r.Warnings = append(r.Warnings, model.Warning{
				ID:       "final_snapshot_unavailable",
				Severity: "info",
				Message:  "Process exited before a final procfs snapshot could be collected.",
				Evidence: map[string]string{
					"pid": fmt.Sprintf("%d", pid),
				},
			})
		} else {
			r.Warnings = append(r.Warnings, collectWarning("final_snapshot_failed", err))
		}
	} else {
		final.Run = nil
		final.InitialSnapshot = nil
		final.FinalSnapshot = nil
		r.FinalSnapshot = &final
		copySnapshotFields(&r, final)
	}

	return r, nil
}

func copySnapshotFields(dst *model.Report, src model.Report) {
	dst.Target = src.Target
	dst.Process = src.Process
	dst.Resources = src.Resources
	dst.FDs = src.FDs
	dst.Maps = src.Maps
	dst.Namespaces = src.Namespaces
	dst.Cgroups = src.Cgroups
	dst.Warnings = append(dst.Warnings, src.Warnings...)
}

func collectWarning(id string, err error) model.Warning {
	return model.Warning{
		ID:       id,
		Severity: "low",
		Message:  "Could not collect procfs snapshot.",
		Evidence: map[string]string{
			"error": err.Error(),
		},
	}
}

func exitStatus(err error) (int, string) {
	if err == nil {
		return 0, ""
	}
	var exitErr *exec.ExitError
	if !errors.As(err, &exitErr) {
		return 1, ""
	}
	status, ok := exitErr.Sys().(syscall.WaitStatus)
	if !ok {
		return exitErr.ExitCode(), ""
	}
	if status.Signaled() {
		return 128 + int(status.Signal()), status.Signal().String()
	}
	return status.ExitStatus(), ""
}
