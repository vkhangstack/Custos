package system

// SetSystemProxy configures the system SOCKS5 proxy
// Implementation is platform-specific in network_*.go files
func SetSystemProxy(enabled bool, port int) error {
	return setProxy(enabled, port)
}
