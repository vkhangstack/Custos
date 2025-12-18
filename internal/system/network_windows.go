package system

import (
	"fmt"
	"log"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

var (
	modwininet            = windows.NewLazySystemDLL("wininet.dll")
	procInternetSetOption = modwininet.NewProc("InternetSetOptionW")
)

const (
	INTERNET_OPTION_SETTINGS_CHANGED = 39
	INTERNET_OPTION_REFRESH          = 37
)

func setProxy(enabled bool, port int) error {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Internet Settings`, registry.SET_VALUE)
	if err != nil {
		return fmt.Errorf("failed to open registry key: %w", err)
	}
	defer k.Close()

	if enabled {
		proxyServer := fmt.Sprintf("socks=socks5://127.0.0.1:%d", port)
		if err = k.SetStringValue("ProxyServer", proxyServer); err != nil {
			return fmt.Errorf("failed to set ProxyServer: %w", err)
		}
		if err = k.SetStringValue("ProxyOverride", "localhost;127.*;<local>"); err != nil {
			return fmt.Errorf("failed to set ProxyOverride: %w", err)
		}
		if err = k.SetDWordValue("ProxyEnable", 1); err != nil {
			return fmt.Errorf("failed to enable ProxyEnable: %w", err)
		}
		log.Printf("Enabled system proxy at %s", proxyServer)
	} else {
		if err = k.SetDWordValue("ProxyEnable", 0); err != nil {
			return fmt.Errorf("failed to disable ProxyEnable: %w", err)
		}
		log.Println("Disabled system proxy")
	}

	// Notify system of settings change
	// Call InternetSetOption to refresh settings immediately
	// 0 means NULL handle
	ret, _, err := procInternetSetOption.Call(0, uintptr(INTERNET_OPTION_SETTINGS_CHANGED), 0, 0)
	if ret == 0 {
		log.Printf("InternetSetOption(SETTINGS_CHANGED) failed: %v", err)
	}

	ret, _, err = procInternetSetOption.Call(0, uintptr(INTERNET_OPTION_REFRESH), 0, 0)
	if ret == 0 {
		log.Printf("InternetSetOption(REFRESH) failed: %v", err)
	}

	return nil
}
