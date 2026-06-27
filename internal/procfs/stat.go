package procfs

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type StatInfo struct {
	State          string
	PPID           int
	StartTimeTicks uint64
}

func ReadStat(pid int) (StatInfo, error) {
	data, err := os.ReadFile(procPIDPath(pid, "stat"))
	if err != nil {
		return StatInfo{}, err
	}
	return ParseStat(string(data))
}

func ParseStat(content string) (StatInfo, error) {
	content = strings.TrimSpace(content)
	open := strings.IndexByte(content, '(')
	close := strings.LastIndexByte(content, ')')
	if open < 0 || close < open {
		return StatInfo{}, fmt.Errorf("invalid proc stat: missing comm parentheses")
	}

	fields := strings.Fields(strings.TrimSpace(content[close+1:]))
	if len(fields) < 20 {
		return StatInfo{}, fmt.Errorf("invalid proc stat: got %d fields after comm", len(fields))
	}

	ppid, err := strconv.Atoi(fields[1])
	if err != nil {
		return StatInfo{}, fmt.Errorf("invalid ppid: %w", err)
	}
	start, err := strconv.ParseUint(fields[19], 10, 64)
	if err != nil {
		return StatInfo{}, fmt.Errorf("invalid starttime: %w", err)
	}

	return StatInfo{
		State:          fields[0],
		PPID:           ppid,
		StartTimeTicks: start,
	}, nil
}
