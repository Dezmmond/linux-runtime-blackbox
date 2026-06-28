package model

type Report struct {
	SchemaVersion string            `json:"schema_version"`
	CollectedAt   string            `json:"collected_at"`
	Run           *RunInfo          `json:"run,omitempty"`
	Target        Target            `json:"target"`
	Process       Process           `json:"process"`
	Resources     Resources         `json:"resources"`
	FDs           []FD              `json:"fds"`
	Maps          []MapEntry        `json:"maps"`
	Namespaces    map[string]string `json:"namespaces"`
	Cgroups       []Cgroup          `json:"cgroups"`
	Warnings      []Warning         `json:"warnings"`

	InitialSnapshot *Report `json:"initial_snapshot,omitempty"`
	FinalSnapshot   *Report `json:"final_snapshot,omitempty"`
}

type RunInfo struct {
	Command        []string `json:"command"`
	PID            int      `json:"pid"`
	StartedAt      string   `json:"started_at"`
	EndedAt        string   `json:"ended_at,omitempty"`
	DurationMillis int64    `json:"duration_millis,omitempty"`
	ExitCode       int      `json:"exit_code"`
	Signal         string   `json:"signal,omitempty"`
}
