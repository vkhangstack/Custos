//go:build windows

package system

import (
	"os"
	"os/exec"
	"strings"
	"syscall"
)

const registryKey = `HKCU\Software\Microsoft\Windows\CurrentVersion\Run`
const appName = "Custos"

func setStartup(enabled bool) error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	if enabled {
		// Add to registry: reg add HKCU\... /v AppName /t REG_SZ /d "path" /f
		cmd := exec.Command("reg", "add", registryKey, "/v", appName, "/t", "REG_SZ", "/d", exePath, "/f")
		// Hide window
		cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true} // Oops, syscall needs import or use generic exec
		// Just run it normally
		return cmd.Run()
	}

	// Remove from registry: reg delete HKCU\... /v AppName /f
	return exec.Command("reg", "delete", registryKey, "/v", appName, "/f").Run()
}

func isStartupEnabled() (bool, error) {
	// Query registry: reg query HKCU\... /v AppName
	out, err := exec.Command("reg", "query", registryKey, "/v", appName).Output()
	if err != nil {
		// If key doesn't exist, it returns error
		return false, nil
	}
	return strings.Contains(string(out), appName), nil
}
