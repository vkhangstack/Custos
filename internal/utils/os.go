package utils

import (
	"runtime"
)

const (
	Windows = "windows"
	Darwin  = "darwin"
	Linux   = "linux"
	Unknown = "unknown"
)

func GetOS() string {
	switch os := runtime.GOOS; os {
	case "windows":
		return Windows
	case "darwin":
		return Darwin
	case "linux":
		return Linux
	default:
		return Unknown
	}
}
