package procfs

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var ProcRoot = "/proc"

func procPIDPath(pid int, parts ...string) string {
	all := append([]string{ProcRoot, strconvItoa(pid)}, parts...)
	return filepath.Join(all...)
}

func strconvItoa(v int) string {
	return strconv.Itoa(v)
}

func ProcessExists(pid int) (bool, error) {
	_, err := os.Stat(procPIDPath(pid))
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

func ReadComm(pid int) (string, error) {
	data, err := os.ReadFile(procPIDPath(pid, "comm"))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

func ReadCmdline(pid int) ([]string, error) {
	data, err := os.ReadFile(procPIDPath(pid, "cmdline"))
	if err != nil {
		return nil, err
	}
	data = bytes.TrimRight(data, "\x00")
	if len(data) == 0 {
		return nil, nil
	}
	parts := bytes.Split(data, []byte{0})
	out := make([]string, 0, len(parts))
	for _, part := range parts {
		out = append(out, string(part))
	}
	return out, nil
}

func ReadExe(pid int) (string, error) {
	return os.Readlink(procPIDPath(pid, "exe"))
}

func ReadCWD(pid int) (string, error) {
	return os.Readlink(procPIDPath(pid, "cwd"))
}
