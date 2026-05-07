package logic

import "strings"

func extractPathParam(path, prefix string) string {
	idx := strings.Index(path, prefix)
	if idx == -1 {
		return ""
	}
	result := path[idx+len(prefix):]
	if i := strings.Index(result, "/"); i != -1 {
		result = result[:i]
	}
	if i := strings.Index(result, "?"); i != -1 {
		result = result[:i]
	}
	return result
}
