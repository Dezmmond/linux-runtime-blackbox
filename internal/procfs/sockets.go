package procfs

import "strings"

func SocketInode(target string) string {
	if !strings.HasPrefix(target, "socket:[") || !strings.HasSuffix(target, "]") {
		return ""
	}
	return strings.TrimSuffix(strings.TrimPrefix(target, "socket:["), "]")
}
