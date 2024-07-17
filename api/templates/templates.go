// Package templates contains the virtual filesystem containing the email templates.
package templates

import "embed"

//go:embed *
var FS embed.FS
