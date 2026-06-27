package model

type FD struct {
	FD     int    `json:"fd"`
	Target string `json:"target"`
	Kind   string `json:"kind"`
}

type MapEntry struct {
	StartAddress string `json:"start_address,omitempty"`
	EndAddress   string `json:"end_address,omitempty"`
	Perms        string `json:"perms"`
	Path         string `json:"path,omitempty"`
}

type Cgroup struct {
	Hierarchy   string `json:"hierarchy"`
	Controllers string `json:"controllers"`
	Path        string `json:"path"`
}
