package procfs

import "testing"

func TestClassifyFDTarget(t *testing.T) {
	tests := map[string]string{
		"socket:[123]":           "socket",
		"pipe:[123]":             "pipe",
		"anon_inode:[eventpoll]": "anon_inode",
		"memfd:tmp":              "memfd",
		"/tmp/example (deleted)": "deleted_file",
		"/home/user/example.txt": "file",
		"":                       "unknown",
	}
	for target, want := range tests {
		if got := ClassifyFDTarget(target); got != want {
			t.Fatalf("ClassifyFDTarget(%q) = %q, want %q", target, got, want)
		}
	}
}
