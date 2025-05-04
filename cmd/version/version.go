package version

import (
	"fmt"
)

// These variables are set during build time via -ldflags
var (
	version   = "n/a"
	gitCommit = "n/a"
)

func VersionTemplate() string {
	return fmt.Sprintf("Version: %s\nGit commit: %s\n", version, gitCommit)
}
