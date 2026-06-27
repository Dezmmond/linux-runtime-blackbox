package model

type Warning struct {
	ID       string            `json:"id"`
	Severity string            `json:"severity"`
	Message  string            `json:"message"`
	Evidence map[string]string `json:"evidence,omitempty"`
}
