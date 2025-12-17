package system

// SetStartup configures the application to run on system startup
func SetStartup(enabled bool) error {
	return setStartup(enabled)
}

// IsStartupEnabled checks if the application is configured to run on startup
func IsStartupEnabled() (bool, error) {
	return isStartupEnabled()
}
