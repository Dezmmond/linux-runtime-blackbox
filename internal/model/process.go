package model

type Target struct {
	PID     int      `json:"pid"`
	Comm    string   `json:"comm,omitempty"`
	Exe     string   `json:"exe,omitempty"`
	Cmdline []string `json:"cmdline,omitempty"`
	CWD     string   `json:"cwd,omitempty"`
}

type Process struct {
	PID            int    `json:"pid"`
	PPID           int    `json:"ppid,omitempty"`
	UID            int    `json:"uid,omitempty"`
	GID            int    `json:"gid,omitempty"`
	State          string `json:"state,omitempty"`
	Threads        int    `json:"threads,omitempty"`
	StartTimeTicks uint64 `json:"start_time_ticks,omitempty"`
}

type Resources struct {
	VmSizeKB   uint64 `json:"vmsize_kb,omitempty"`
	VmRSSKB    uint64 `json:"vmrss_kb,omitempty"`
	ReadBytes  uint64 `json:"read_bytes,omitempty"`
	WriteBytes uint64 `json:"write_bytes,omitempty"`
}
