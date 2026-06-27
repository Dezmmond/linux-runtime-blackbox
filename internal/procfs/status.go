package procfs

import (
	"os"
	"strconv"
	"strings"

	"github.com/Dezmmond/linux-runtime-blackbox/internal/model"
)

type StatusInfo struct {
	UID      int
	GID      int
	State    string
	Threads  int
	VmSizeKB uint64
	VmRSSKB  uint64
}

func ReadStatus(pid int) (StatusInfo, error) {
	data, err := os.ReadFile(procPIDPath(pid, "status"))
	if err != nil {
		return StatusInfo{}, err
	}
	return ParseStatus(string(data)), nil
}

func ParseStatus(content string) StatusInfo {
	var info StatusInfo
	for _, line := range strings.Split(content, "\n") {
		key, value, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		fields := strings.Fields(value)
		if len(fields) == 0 {
			continue
		}
		switch key {
		case "State":
			info.State = fields[0]
		case "Uid":
			info.UID = atoiDefault(fields[0])
		case "Gid":
			info.GID = atoiDefault(fields[0])
		case "Threads":
			info.Threads = atoiDefault(fields[0])
		case "VmSize":
			info.VmSizeKB = atou64Default(fields[0])
		case "VmRSS":
			info.VmRSSKB = atou64Default(fields[0])
		}
	}
	return info
}

func ReadIO(pid int) (model.Resources, error) {
	data, err := os.ReadFile(procPIDPath(pid, "io"))
	if err != nil {
		return model.Resources{}, err
	}
	var resources model.Resources
	for _, line := range strings.Split(string(data), "\n") {
		key, value, ok := strings.Cut(line, ":")
		if !ok {
			continue
		}
		fields := strings.Fields(value)
		if len(fields) == 0 {
			continue
		}
		switch key {
		case "read_bytes":
			resources.ReadBytes = atou64Default(fields[0])
		case "write_bytes":
			resources.WriteBytes = atou64Default(fields[0])
		}
	}
	return resources, nil
}

func atoiDefault(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}

func atou64Default(s string) uint64 {
	v, _ := strconv.ParseUint(s, 10, 64)
	return v
}
