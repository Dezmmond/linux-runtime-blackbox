package model

type Socket struct {
	FD    int    `json:"fd,omitempty"`
	Inode string `json:"inode,omitempty"`
	Kind  string `json:"kind,omitempty"`
}
