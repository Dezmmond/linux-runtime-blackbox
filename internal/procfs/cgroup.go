package procfs

import (
	"os"
	"strings"

	"github.com/Dezmmond/linux-runtime-blackbox/internal/model"
)

func ReadCgroups(pid int) ([]model.Cgroup, error) {
	data, err := os.ReadFile(procPIDPath(pid, "cgroup"))
	if err != nil {
		return nil, err
	}
	var cgroups []model.Cgroup
	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ":", 3)
		if len(parts) != 3 {
			continue
		}
		cgroups = append(cgroups, model.Cgroup{
			Hierarchy:   parts[0],
			Controllers: parts[1],
			Path:        parts[2],
		})
	}
	return cgroups, nil
}
