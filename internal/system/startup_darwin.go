//go:build darwin

package system

import "fmt"

func setStartup(enabled bool) error {
	return fmt.Errorf("not implemented for darwin")
}

func isStartupEnabled() (bool, error) {
	return false, nil
}
