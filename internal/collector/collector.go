package collector

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Dezmmond/linux-runtime-blackbox/internal/model"
	"github.com/Dezmmond/linux-runtime-blackbox/internal/procfs"
)

const SchemaVersion = "0.1"

var ErrProcessNotFound = errors.New("target process not found")

type Options struct {
	ManyFDThreshold int
	LargeRSSKB      uint64
}

func DefaultOptions() Options {
	return Options{
		ManyFDThreshold: 256,
		LargeRSSKB:      1024 * 1024,
	}
}

func CollectPID(pid int) (model.Report, error) {
	return CollectPIDWithOptions(pid, DefaultOptions())
}

func CollectPIDWithOptions(pid int, opts Options) (model.Report, error) {
	exists, err := procfs.ProcessExists(pid)
	if err != nil {
		return model.Report{}, err
	}
	if !exists {
		return model.Report{}, ErrProcessNotFound
	}

	r := model.Report{
		SchemaVersion: SchemaVersion,
		CollectedAt:   time.Now().UTC().Format(time.RFC3339Nano),
		Target: model.Target{
			PID: pid,
		},
		Process: model.Process{
			PID: pid,
		},
		Namespaces: map[string]string{},
		Warnings:   []model.Warning{},
	}

	var partial bool
	addPartial := func(source string, err error) {
		if err == nil {
			return
		}
		partial = true
		r.Warnings = append(r.Warnings, model.Warning{
			ID:       "partial_procfs_read",
			Severity: "low",
			Message:  "Could not read one procfs source; report may be incomplete.",
			Evidence: map[string]string{
				"source": source,
				"error":  err.Error(),
			},
		})
	}

	comm, err := procfs.ReadComm(pid)
	if err != nil {
		addPartial("comm", err)
	} else {
		r.Target.Comm = comm
	}

	cmdline, err := procfs.ReadCmdline(pid)
	if err != nil {
		addPartial("cmdline", err)
	} else {
		r.Target.Cmdline = cmdline
	}

	exe, err := procfs.ReadExe(pid)
	if err != nil {
		addPartial("exe", err)
		r.Warnings = append(r.Warnings, model.Warning{
			ID:       "unknown_exe",
			Severity: "low",
			Message:  "Executable path could not be resolved from procfs.",
			Evidence: map[string]string{
				"error": err.Error(),
			},
		})
	} else {
		r.Target.Exe = exe
	}

	cwd, err := procfs.ReadCWD(pid)
	if err != nil {
		addPartial("cwd", err)
	} else {
		r.Target.CWD = cwd
	}

	status, err := procfs.ReadStatus(pid)
	if err != nil {
		addPartial("status", err)
	} else {
		r.Process.UID = status.UID
		r.Process.GID = status.GID
		r.Process.State = status.State
		r.Process.Threads = status.Threads
		r.Resources.VmSizeKB = status.VmSizeKB
		r.Resources.VmRSSKB = status.VmRSSKB
	}

	stat, err := procfs.ReadStat(pid)
	if err != nil {
		addPartial("stat", err)
	} else {
		r.Process.PPID = stat.PPID
		if r.Process.State == "" {
			r.Process.State = stat.State
		}
		r.Process.StartTimeTicks = stat.StartTimeTicks
	}

	io, err := procfs.ReadIO(pid)
	if err != nil {
		addPartial("io", err)
	} else {
		r.Resources.ReadBytes = io.ReadBytes
		r.Resources.WriteBytes = io.WriteBytes
	}

	fds, err := procfs.ReadFDs(pid)
	if err != nil {
		addPartial("fd", err)
	} else {
		r.FDs = fds
	}

	maps, err := procfs.ReadMaps(pid)
	if err != nil {
		addPartial("maps", err)
	} else {
		r.Maps = maps
	}

	namespaces, err := procfs.ReadNamespaces(pid)
	if err != nil {
		addPartial("ns", err)
	} else {
		r.Namespaces = namespaces
	}

	cgroups, err := procfs.ReadCgroups(pid)
	if err != nil {
		addPartial("cgroup", err)
	} else {
		r.Cgroups = cgroups
	}

	r.Warnings = append(r.Warnings, heuristicWarnings(r, opts)...)
	if partial {
		r.Warnings = append(r.Warnings, model.Warning{
			ID:       "permission_denied_partial_report",
			Severity: "info",
			Message:  "Some procfs sources were unavailable. This can be caused by permissions or a process changing state during collection.",
		})
	}
	return r, nil
}

func heuristicWarnings(r model.Report, opts Options) []model.Warning {
	var warnings []model.Warning
	if opts.ManyFDThreshold <= 0 {
		opts.ManyFDThreshold = DefaultOptions().ManyFDThreshold
	}
	if opts.LargeRSSKB == 0 {
		opts.LargeRSSKB = DefaultOptions().LargeRSSKB
	}

	if len(r.FDs) > opts.ManyFDThreshold {
		warnings = append(warnings, model.Warning{
			ID:       "many_open_fds",
			Severity: "medium",
			Message:  "Process has many open file descriptors.",
			Evidence: map[string]string{
				"count":     fmt.Sprintf("%d", len(r.FDs)),
				"threshold": fmt.Sprintf("%d", opts.ManyFDThreshold),
			},
		})
	}

	hasSocket := false
	for _, fd := range r.FDs {
		if strings.HasSuffix(fd.Target, " (deleted)") || fd.Kind == "deleted_file" {
			warnings = append(warnings, model.Warning{
				ID:       "deleted_file_open",
				Severity: "low",
				Message:  "A file descriptor points to a deleted file.",
				Evidence: map[string]string{
					"fd":     fmt.Sprintf("%d", fd.FD),
					"target": fd.Target,
				},
			})
		}
		if fd.Kind == "socket" {
			hasSocket = true
		}
	}
	if hasSocket {
		warnings = append(warnings, model.Warning{
			ID:       "network_socket_open",
			Severity: "info",
			Message:  "Process has socket file descriptors.",
		})
	}

	if r.Target.Exe != "" && isTempExecutable(r.Target.Exe) {
		warnings = append(warnings, model.Warning{
			ID:       "unexpected_tmp_executable",
			Severity: "medium",
			Message:  "Executable path is under a temporary or shared-memory directory.",
			Evidence: map[string]string{
				"exe": r.Target.Exe,
			},
		})
	}

	if r.Resources.VmRSSKB > opts.LargeRSSKB {
		warnings = append(warnings, model.Warning{
			ID:       "large_rss",
			Severity: "info",
			Message:  "Process RSS is above the configured threshold.",
			Evidence: map[string]string{
				"rss_kb":    fmt.Sprintf("%d", r.Resources.VmRSSKB),
				"threshold": fmt.Sprintf("%d", opts.LargeRSSKB),
			},
		})
	}

	return warnings
}

func isTempExecutable(path string) bool {
	clean := strings.TrimSuffix(path, " (deleted)")
	for _, prefix := range []string{"/tmp/", "/var/tmp/", "/dev/shm/"} {
		if strings.HasPrefix(clean, prefix) {
			return true
		}
	}
	home := os.Getenv("HOME")
	if home != "" && strings.HasPrefix(clean, home+"/tmp/") {
		return true
	}
	return false
}
