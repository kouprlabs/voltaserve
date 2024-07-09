package helper

import "strings"

func ExtractWorkspaceIDFromPath(path string) string {
	slashParts := strings.Split(strings.TrimPrefix(path, "/"), "/")
	dashParts := strings.Split(slashParts[0], "-")
	if len(dashParts) > 1 {
		return dashParts[len(dashParts)-1]
	}
	return ""
}
