package utils

import (
	"path/filepath"
	"strings"
)

var supportedAudioExtensions = map[string]struct{}{
	".wav": {},
	".mp3": {},
}

var supportedVideoExtensions = map[string]struct{}{
	".mp4": {},
	".mov": {},
}

func IsAudioFile(file string) bool {
	ext := filepath.Ext(file)
	_, ok := supportedAudioExtensions[strings.ToLower(ext)]
	return ok
}

func IsVideoFile(file string) bool {
	ext := filepath.Ext(file)
	_, ok := supportedVideoExtensions[strings.ToLower(ext)]
	return ok
}
