package utils

import (
	"path/filepath"
	"strings"
)

var SupportedAudioExtensions = map[string]struct{}{
	".wav": {},
	".mp3": {},
}

var SupportedVideoExtensions = map[string]struct{}{
	".mp4": {},
	".mkv": {},
	".avi": {},
	".mov": {},
	".wmv": {},
}

func IsAudioFile(file string) bool {
	ext := filepath.Ext(file)
	_, ok := SupportedAudioExtensions[strings.ToLower(ext)]
	return ok
}

func IsVideoFile(file string) bool {
	ext := filepath.Ext(file)
	_, ok := SupportedVideoExtensions[strings.ToLower(ext)]
	return ok
}
