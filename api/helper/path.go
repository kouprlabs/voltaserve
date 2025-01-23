package helper

import "strings"

func PathFromFilename(name string) []string {
	var components []string
	for _, component := range strings.Split(name, "/") {
		if component != "" {
			components = append(components, component)
		}
	}
	return components
}

func FilenameFromPath(path []string) string {
	return path[len(path)-1]
}
