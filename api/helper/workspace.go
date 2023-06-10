package helper

import (
	"fmt"
	"strings"

	"github.com/gosimple/slug"
)

func WorkspaceToSlug(id string, name string) string {
	return fmt.Sprintf("%s-%s", slug.Make(name), id)
}

func SlugToWorkspaceId(slug string) string {
	parts := strings.Split(slug, "-")
	return parts[len(parts)-1]
}
