package helper

import (
	"fmt"
	"strings"

	"github.com/gosimple/slug"
)

func SlugFromWorkspace(id string, name string) string {
	return fmt.Sprintf("%s-%s", slug.Make(name), id)
}

func WorkspaceIDFromSlug(slug string) string {
	parts := strings.Split(slug, "-")
	return parts[len(parts)-1]
}
