package procfs

import "testing"

func TestParseStatWithSpacesInComm(t *testing.T) {
	stat, err := ParseStat("1234 (worker thread) S 42 1 1 0 -1 4194304 1 2 3 4 5 6 7 8 20 0 1 0 987654 0 0")
	if err != nil {
		t.Fatalf("ParseStat returned error: %v", err)
	}
	if stat.State != "S" {
		t.Fatalf("state = %q, want S", stat.State)
	}
	if stat.PPID != 42 {
		t.Fatalf("ppid = %d, want 42", stat.PPID)
	}
	if stat.StartTimeTicks != 987654 {
		t.Fatalf("starttime = %d, want 987654", stat.StartTimeTicks)
	}
}
