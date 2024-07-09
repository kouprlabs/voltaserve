package helper

import "strings"

func IsMicrosoftOfficeLockFile(name string) bool {
	return strings.HasPrefix(name, "~$")
}

func IsOpenOfficeOfficeLockFile(name string) bool {
	return strings.HasPrefix(name, ".~lock.") && strings.HasSuffix(name, "#")
}
