package utils

import (
	"path/filepath"
	"strings"
)

var supportedAudioExtensions = map[string]struct{}{
	".wav": {},
	".mp3": {},
}

func IsAudioFile(file string) bool {
	ext := filepath.Ext(file)
	_, ok := supportedAudioExtensions[strings.ToLower(ext)]
	return ok
}
