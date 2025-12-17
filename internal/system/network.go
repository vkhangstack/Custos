package system

import (
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strconv"
)

// SetSystemProxy configures the system SOCKS5 proxy
// Currently supports GNOME via gsettings
func SetSystemProxy(enabled bool, port int) error {
	switch runtime.GOOS {
	case "linux":
		// Check for gsettings
		if _, err := exec.LookPath("gsettings"); err != nil {
			log.Println("gsettings not found, skipping system proxy configuration")
			return nil
		}
		if enabled {
			return enableGnomeProxy(port)
		}
		return disableGnomeProxy()
	case "darwin", "windows":
		log.Printf("System proxy configuration not yet implemented for %s", runtime.GOOS)
		return nil
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}
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
