package system

import (
	"os/exec"
	"strconv"
)

// SetSystemProxy configures the system SOCKS5 proxy
// Currently supports GNOME via gsettings
func SetSystemProxy(enabled bool, port int) error {
	if enabled {
		return enableGnomeProxy(port)
	}
	return disableGnomeProxy()
}

func enableGnomeProxy(port int) error {
	// Set SOCKS host
	if err := exec.Command("gsettings", "set", "org.gnome.system.proxy.socks", "host", "127.0.0.1").Run(); err != nil {
		return err
	}
	// Set SOCKS port
	if err := exec.Command("gsettings", "set", "org.gnome.system.proxy.socks", "port", strconv.Itoa(port)).Run(); err != nil {
		return err
	}
	// Set mode to manual
	if err := exec.Command("gsettings", "set", "org.gnome.system.proxy", "mode", "manual").Run(); err != nil {
		return err
	}
	return nil
}

func disableGnomeProxy() error {
	// Set mode to none
	return exec.Command("gsettings", "set", "org.gnome.system.proxy", "mode", "none").Run()
}
