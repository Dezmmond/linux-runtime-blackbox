package report

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/yourname/linux-runtime-blackbox/internal/model"
)

func WritePretty(w io.Writer, r model.Report) error {
	write := func(format string, args ...any) error {
		_, err := fmt.Fprintf(w, format, args...)
		return err
	}

	if err := write("Linux Runtime Blackbox report\n"); err != nil {
		return err
	}
	if err := write("schema: %s\ncollected_at: %s\n\n", r.SchemaVersion, r.CollectedAt); err != nil {
		return err
	}
	if err := write("Target\n"); err != nil {
		return err
	}
	if err := write("  pid: %d\n  comm: %s\n  exe: %s\n  cwd: %s\n", r.Target.PID, dash(r.Target.Comm), dash(r.Target.Exe), dash(r.Target.CWD)); err != nil {
		return err
	}
	if len(r.Target.Cmdline) > 0 {
		if err := write("  cmdline: %s\n", strings.Join(r.Target.Cmdline, " ")); err != nil {
			return err
		}
	}

	if err := write("\nProcess\n  ppid: %d\n  uid: %d\n  gid: %d\n  state: %s\n  threads: %d\n  start_time_ticks: %d\n",
		r.Process.PPID, r.Process.UID, r.Process.GID, dash(r.Process.State), r.Process.Threads, r.Process.StartTimeTicks); err != nil {
		return err
	}
	if err := write("\nResources\n  vmsize_kb: %d\n  vmrss_kb: %d\n  read_bytes: %d\n  write_bytes: %d\n",
		r.Resources.VmSizeKB, r.Resources.VmRSSKB, r.Resources.ReadBytes, r.Resources.WriteBytes); err != nil {
		return err
	}

	if err := write("\nFile descriptors (%d)\n", len(r.FDs)); err != nil {
		return err
	}
	for _, fd := range r.FDs {
		if err := write("  %d %-12s %s\n", fd.FD, fd.Kind, fd.Target); err != nil {
			return err
		}
	}

	if err := write("\nMappings (%d)\n", len(r.Maps)); err != nil {
		return err
	}
	limit := len(r.Maps)
	if limit > 25 {
		limit = 25
	}
	for _, m := range r.Maps[:limit] {
		if err := write("  %s-%s %-4s %s\n", m.StartAddress, m.EndAddress, m.Perms, m.Path); err != nil {
			return err
		}
	}
	if len(r.Maps) > limit {
		if err := write("  ... %d more mappings\n", len(r.Maps)-limit); err != nil {
			return err
		}
	}

	if err := write("\nNamespaces\n"); err != nil {
		return err
	}
	nsKeys := make([]string, 0, len(r.Namespaces))
	for k := range r.Namespaces {
		nsKeys = append(nsKeys, k)
	}
	sort.Strings(nsKeys)
	for _, k := range nsKeys {
		if err := write("  %s: %s\n", k, r.Namespaces[k]); err != nil {
			return err
		}
	}

	if err := write("\nCgroups\n"); err != nil {
		return err
	}
	for _, cg := range r.Cgroups {
		if err := write("  %s:%s:%s\n", cg.Hierarchy, cg.Controllers, cg.Path); err != nil {
			return err
		}
	}

	if err := write("\nWarnings (%d)\n", len(r.Warnings)); err != nil {
		return err
	}
	for _, warning := range r.Warnings {
		if err := write("  [%s] %s: %s\n", warning.Severity, warning.ID, warning.Message); err != nil {
			return err
		}
	}
	return nil
}

func dash(value string) string {
	if value == "" {
		return "-"
	}
	return value
}
