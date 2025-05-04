package utils

import (
	"fmt"
	"strings"

	"github.com/pterm/pterm"
)

var DependencyMap = map[string][]Command{
	"convert": {
		{
			Name:  "ffmpeg",
			Probe: "-version",
		},
		{
			Name:  "ffprobe",
			Probe: "-version",
		},
	},
}

func CheckDependencies(cmd string) error {
	dependencies := DependencyMap[cmd]

	missing := []string{}
	installed := []string{}
	for _, dependency := range dependencies {
		if dependency.Exists() {
			installed = append(installed, dependency.Name)
		} else {
			missing = append(missing, dependency.Name)
		}
	}

	if len(missing) > 0 {
		if len(installed) > 0 {
			pterm.ThemeDefault.SuccessMessageStyle.Println("Installed dependencies:")
			for _, dependency := range installed {
				pterm.ThemeDefault.SuccessMessageStyle.Printf("✓ %s\n", dependency)
			}
		}

		pterm.ThemeDefault.WarningMessageStyle.Println("Missing dependencies:")
		for _, dependency := range missing {
			pterm.ThemeDefault.WarningMessageStyle.Printf("✗ %s\n", dependency)
		}

		return fmt.Errorf("missing dependencies: %s", strings.Join(missing, ", "))
	}

	return nil
}
