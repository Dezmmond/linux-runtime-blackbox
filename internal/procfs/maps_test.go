package procfs

import "testing"

func TestParseMapLine(t *testing.T) {
	entry, ok := ParseMapLine("7f000000-7f001000 r-xp 00000000 08:01 123 /usr/lib/libc.so.6")
	if !ok {
		t.Fatal("ParseMapLine returned !ok")
	}
	if entry.StartAddress != "7f000000" || entry.EndAddress != "7f001000" {
		t.Fatalf("address = %s-%s", entry.StartAddress, entry.EndAddress)
	}
	if entry.Perms != "r-xp" {
		t.Fatalf("perms = %q, want r-xp", entry.Perms)
	}
	if entry.Path != "/usr/lib/libc.so.6" {
		t.Fatalf("path = %q", entry.Path)
	}
}
