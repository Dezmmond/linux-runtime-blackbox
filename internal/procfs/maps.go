package procfs

import (
	"bufio"
	"os"
	"strings"

	"github.com/Dezmmond/linux-runtime-blackbox/internal/model"
)

func ReadMaps(pid int) ([]model.MapEntry, error) {
	f, err := os.Open(procPIDPath(pid, "maps"))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var entries []model.MapEntry
	seen := make(map[string]struct{})
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		entry, ok := ParseMapLine(scanner.Text())
		if !ok {
			continue
		}
		key := entry.StartAddress + "|" + entry.EndAddress + "|" + entry.Perms + "|" + entry.Path
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		entries = append(entries, entry)
	}
	if err := scanner.Err(); err != nil {
		return entries, err
	}
	return entries, nil
}

func ParseMapLine(line string) (model.MapEntry, bool) {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return model.MapEntry{}, false
	}
	start, end, ok := strings.Cut(fields[0], "-")
	if !ok {
		return model.MapEntry{}, false
	}
	entry := model.MapEntry{
		StartAddress: start,
		EndAddress:   end,
		Perms:        fields[1],
	}
	if len(fields) >= 6 {
		entry.Path = strings.Join(fields[5:], " ")
	}
	return entry, true
}
