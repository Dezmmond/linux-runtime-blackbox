package procfs

import "testing"

func TestParseStatus(t *testing.T) {
	status := ParseStatus(`Name:	python3
State:	S (sleeping)
Uid:	1000	1000	1000	1000
Gid:	1000	1000	1000	1000
VmSize:	  123456 kB
VmRSS:	   23456 kB
Threads:	4
`)
	if status.UID != 1000 || status.GID != 1000 {
		t.Fatalf("uid/gid = %d/%d, want 1000/1000", status.UID, status.GID)
	}
	if status.State != "S" {
		t.Fatalf("state = %q, want S", status.State)
	}
	if status.Threads != 4 {
		t.Fatalf("threads = %d, want 4", status.Threads)
	}
	if status.VmSizeKB != 123456 || status.VmRSSKB != 23456 {
		t.Fatalf("vm = %d/%d, want 123456/23456", status.VmSizeKB, status.VmRSSKB)
	}
}
