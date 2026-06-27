package model

type Report struct {
	SchemaVersion string            `json:"schema_version"`
	CollectedAt   string            `json:"collected_at"`
	Target        Target            `json:"target"`
	Process       Process           `json:"process"`
	Resources     Resources         `json:"resources"`
	FDs           []FD              `json:"fds"`
	Maps          []MapEntry        `json:"maps"`
	Namespaces    map[string]string `json:"namespaces"`
	Cgroups       []Cgroup          `json:"cgroups"`
	Warnings      []Warning         `json:"warnings"`
}
