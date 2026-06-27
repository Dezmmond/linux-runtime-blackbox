package procfs

import (
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/yourname/linux-runtime-blackbox/internal/model"
)

func ReadFDs(pid int) ([]model.FD, error) {
	dir := procPIDPath(pid, "fd")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	fds := make([]model.FD, 0, len(entries))
	for _, entry := range entries {
		fdNum, err := strconv.Atoi(entry.Name())
		if err != nil {
			continue
		}
		target, err := os.Readlink(filepath.Join(dir, entry.Name()))
		if err != nil {
			target = ""
		}
		fds = append(fds, model.FD{
			FD:     fdNum,
			Target: target,
			Kind:   ClassifyFDTarget(target),
		})
	}

	sort.Slice(fds, func(i, j int) bool {
		return fds[i].FD < fds[j].FD
	})
	return fds, nil
}

func ClassifyFDTarget(target string) string {
	switch {
	case target == "":
		return "unknown"
	case strings.HasPrefix(target, "socket:["):
		return "socket"
	case strings.HasPrefix(target, "pipe:["):
		return "pipe"
	case strings.HasPrefix(target, "anon_inode:"):
		return "anon_inode"
	case strings.HasPrefix(target, "memfd:"):
		return "memfd"
	case strings.Contains(target, " (deleted)"):
		return "deleted_file"
	default:
		return "file"
	}
}
