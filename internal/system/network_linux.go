package system

import (
	"log"
	"os/exec"
	"strconv"
)

func setProxy(enabled bool, port int) error {
	// Check for gsettings
	if _, err := exec.LookPath("gsettings"); err != nil {
		log.Println("gsettings not found, skipping system proxy configuration")
		return nil
	}
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
	// Set ignore hosts
	// Prevent loopback traffic from going through proxy
	if err := exec.Command("gsettings", "set", "org.gnome.system.proxy", "ignore-hosts", "['localhost', '127.0.0.0/8', '::1']").Run(); err != nil {
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
