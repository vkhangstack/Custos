//go:build linux

package system

import (
	"fmt"
	"os"
	"path/filepath"
)

func setStartup(enabled bool) error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}
	autostartDir := filepath.Join(configDir, "autostart")
	desktopFile := filepath.Join(autostartDir, "custos.desktop")

	if enabled {
		if err := os.MkdirAll(autostartDir, 0755); err != nil {
			return err
		}

		exePath, err := os.Executable()
		if err != nil {
			return err
		}

		content := fmt.Sprintf(`[Desktop Entry]
Type=Application
Name=Custos
Comment=Custos Network Monitor
Exec=%s
Hidden=false
NoDisplay=false
X-GNOME-Autostart-enabled=true
`, exePath)

		return os.WriteFile(desktopFile, []byte(content), 0644)
	}

	// Remove file
	if err := os.Remove(desktopFile); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func isStartupEnabled() (bool, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return false, err
	}
	desktopFile := filepath.Join(configDir, "autostart", "custos.desktop")

	if _, err := os.Stat(desktopFile); err == nil {
		return true, nil
	}
	return false, nil
}
