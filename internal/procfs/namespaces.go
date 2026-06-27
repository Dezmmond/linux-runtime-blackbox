package procfs

import (
	"os"
	"path/filepath"
	"sort"
)

func ReadNamespaces(pid int) (map[string]string, error) {
	dir := procPIDPath(pid, "ns")
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		names = append(names, entry.Name())
	}
	sort.Strings(names)

	out := make(map[string]string, len(names))
	for _, name := range names {
		target, err := os.Readlink(filepath.Join(dir, name))
		if err != nil {
			continue
		}
		out[name] = target
	}
	return out, nil
}
