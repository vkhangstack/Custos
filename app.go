package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/vkhangstack/Custos/internal/system"

	"github.com/vkhangstack/Custos/internal/store"

	"github.com/vkhangstack/Custos/internal/proxy"

	"github.com/vkhangstack/Custos/internal/dns"

	"github.com/vkhangstack/Custos/internal/core"
	"github.com/vkhangstack/Custos/internal/utils"

	_ "embed"
)

//go:embed app.json
var appJson []byte

// App struct
type App struct {
	ctx           context.Context
	store         store.Store
	proxyServer   *proxy.Server
	dnsServer     *dns.Server
	systemTracker *system.Tracker
	blocklist     *core.BlocklistManager
	refreshMu     sync.Mutex
}

// NewApp creates a new App application struct
func NewApp() *App {
	// Resolve log path
	homeDir, errPath := os.UserHomeDir()
	if errPath != nil {
		println("Error getting home directory:", errPath.Error())
		return nil
	}
	dataPath := filepath.Join(homeDir, ".custos", "data")
	if err := os.MkdirAll(dataPath, 0755); err != nil {
		println("Error creating data directory:", err.Error())
		return nil
	}

	// Try to initialize SQLite store
	var s store.Store
	var err error
	s, err = store.NewSQLiteStore(filepath.Join(dataPath, "custos.db"))
	if err != nil {
		fmt.Printf("Failed to initialize SQLite store: %v. Falling back to MemoryStore.\n", err)
		s = store.NewMemoryStore()
	}

	bm := core.NewBlocklistManager()

	// Load blocklist in background
	go bm.Load()

	systemTracker := system.NewTracker()

	// Load configured port
	port := 1080
	if val, err := s.GetSetting("proxy_port"); err == nil && val != "" {
		if p, err := strconv.Atoi(val); err == nil {
			port = p
		}
	}

	return &App{
		store:         s,
		proxyServer:   proxy.NewServer(s, bm, systemTracker, port),
		dnsServer:     dns.NewServer(s, bm, 5353),
		systemTracker: systemTracker,
		blocklist:     bm,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.proxyServer.Start()
	// a.dnsServer.Start()

	// Auto-enable system proxy
	// if err := a.SetSystemProxy(true); err != nil {
	// 	fmt.Printf("Failed to set system proxy on startup: %v\n", err)
	// }

	if enabled := a.GetProtectionStatus(); enabled {
		a.proxyServer.SetProtection(true)
		a.SetSystemProxy(true)
	} else {
		a.proxyServer.SetProtection(false)
		a.SetSystemProxy(false)
	}

	// Restore Adblock State
	if enabled := a.GetAdblockStatus(); enabled {
		a.proxyServer.SetAdblockEnabled(true)
	} else {
		a.proxyServer.SetAdblockEnabled(false)
	}

	// Seed and Refresh Filters
	go func() {
		a.seedFilters()
		a.RefreshAdblockFilters()
	}()

	// Start a ticker to emit logs to frontend
	go a.broadcastLogs()
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {
	// Auto-disable system proxy
	if err := a.SetSystemProxy(false); err != nil {
		fmt.Printf("Failed to disable system proxy on shutdown: %v\n", err)
	}
	a.proxyServer.Stop()
	ctx.Done()
}

// broadcastLogs sends new logs to frontend events
func (a *App) broadcastLogs() {
	a.store.Subscribe(func(entry core.LogEntry) {
		// Emit event to Wails runtime
		// runtime.EventsEmit(a.ctx, "log:new", entry)
		// Since we don't have runtime imported here broadly, we'd typically use it
		// For now let's rely on polling or add runtime via import
	})
}

// GetLogs returns recent logs for the frontend
func (a *App) GetLogs() []core.LogEntry {
	return a.store.GetRecentLogs(50)
}

// GetStats returns current stats
func (a *App) GetStats() core.Stats {
	return a.store.GetStats()
}

// GetSystemConnections returns active system connections
func (a *App) GetSystemConnections() []system.ConnectionInfo {
	conns, _ := a.systemTracker.GetActiveConnections()
	return conns
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

// SetSystemProxy enables or disables the system proxy
func (a *App) SetSystemProxy(enabled bool) error {
	return system.SetSystemProxy(enabled, a.proxyServer.GetPort())
}

// EnableProtection toggles HTTP blocking
func (a *App) EnableProtection(enabled bool) {
	a.SetSystemProxy(enabled)
	// Persist
	val := "false"
	if enabled {
		val = "true"
	}
	a.store.SetSetting("protection_enabled", val)
}

// GetProtectionStatus returns the current status
func (a *App) GetProtectionStatus() bool {
	val, err := a.store.GetSetting("protection_enabled")
	if err != nil {
		return false
	}
	return val == "true"
}

// EnableAdblock toggles the adblock engine
func (a *App) EnableAdblock(enabled bool) {
	a.proxyServer.SetAdblockEnabled(enabled)
	// Persist
	val := "false"
	if enabled {
		val = "true"
	}
	a.store.SetSetting("adblock_enabled", val)
}

// GetAdblockStatus returns the current status
func (a *App) GetAdblockStatus() bool {
	val, err := a.store.GetSetting("adblock_enabled")
	if err != nil || val == "" {
		// Default to enabled if not set
		return true
	}
	return val == "true"
}

// GetChartData returns historical traffic data for the chart
func (a *App) GetChartData(durationStr string) []core.TrafficDataPoint {
	// Parse duration
	var duration time.Duration
	switch durationStr {
	case "1h":
		duration = 1 * time.Hour
	case "3h":
		duration = 3 * time.Hour
	case "24h":
		duration = 24 * time.Hour
	default:
		// Default to 1h if invalid or empty
		duration = 1 * time.Hour
	}
	return a.store.GetTrafficHistory(duration)
}

// Rule Management

// AddRule adds a new rule
func (a *App) AddRule(pattern string, ruleType string) error {
	rule := core.Rule{
		ID:      utils.GenerateIDString(),
		Pattern: pattern,
		Type:    core.RuleType(ruleType),
		Enabled: true,
		Source:  core.RuleSourceCustom,
	}
	return a.store.AddRule(rule)
}

// GetRules returns all rules (legacy/internal use)
func (a *App) GetRules() []core.Rule {
	return a.store.GetRules()
}

// GetRulesPaginated returns rules with pagination
func (a *App) GetRulesPaginated(page, pageSize int, search string) core.PaginatedRulesResponse {
	rules, total, err := a.store.GetRulesPaginated(page, pageSize, search)
	if err != nil {
		return core.PaginatedRulesResponse{Rules: []core.Rule{}, Total: 0}
	}
	return core.PaginatedRulesResponse{Rules: rules, Total: total}
}

// DeleteRule deletes a rule by ID
func (a *App) DeleteRule(id string) error {
	return a.store.DeleteRule(id)
}

// ToggleRule toggles a rule's enabled status
// For simplicity, we just look it up or blindly update?
// Let's implement full UpdateRule in App if needed, or just specific toggle.
// But frontend might just call Delete/Add. Or we need explicit Toggle.
// Let's stick to Add/Get/Delete for MVP as per request "add rules and magement".
func (a *App) ToggleRule(id string, enabled bool) error {
	// We need to fetch it first? Or just blindly update?
	// SQLiteStore UpdateRule updates fields present in struct.
	// We just pass ID and Enabled.
	return a.store.UpdateRule(core.Rule{ID: id, Enabled: enabled})
}

// Startup Management

// SetRunOnStartup toggles launch on startup
func (a *App) SetRunOnStartup(enabled bool) error {
	return system.SetStartup(enabled)
}

// GetStartupStatus checks if launch on startup is enabled
func (a *App) GetStartupStatus() bool {
	enabled, err := system.IsStartupEnabled()
	if err != nil {
		fmt.Printf("Failed to check startup status: %v\n", err)
		return false
	}
	return enabled
}

// AppInfo holds application information
type AppInfo struct {
	Name        string  `json:"name"`
	Version     string  `json:"version"`
	Description string  `json:"description"`
	Author      string  `json:"author"`
	Contact     *string `json:"contact,omitempty"`
}

// GetAppInfo returns application information
func (a *App) GetAppInfo() *AppInfo {
	var appInfo AppInfo
	err := json.Unmarshal(appJson, &appInfo)
	if err != nil {
		fmt.Printf("Error decoding embedded app.json: %v\n", err)
		return &AppInfo{
			Name:    "Custos",
			Version: "Unknown",
		}
	}
	return &appInfo
}

// AppSettings defines configurable settings
type AppSettings struct {
	Port           int  `json:"port"`
	Notifications  bool `json:"notifications"`
	AutoStart      bool `json:"auto_start"`
	AdblockEnabled bool `json:"adblock_enabled"`
}

// GetAppSettings returns current settings
func (a *App) GetAppSettings() AppSettings {
	// Port
	port := a.proxyServer.GetPort()

	// Notifications (TODO: Implement actual notification logic storage if specific)
	// For now assume stored in "notifications_enabled"
	notifications := false
	if val, err := a.store.GetSetting("notifications_enabled"); err == nil && val == "true" {
		notifications = true
	}

	// AutoStart
	autoStart := a.GetStartupStatus()

	// Adblock
	adblockEnabled := a.GetAdblockStatus()

	return AppSettings{
		Port:           port,
		Notifications:  notifications,
		AutoStart:      autoStart,
		AdblockEnabled: adblockEnabled,
	}
}

// SaveAppSettings saves settings and applies changes
func (a *App) SaveAppSettings(settings AppSettings) error {
	// AutoStart
	if err := a.SetRunOnStartup(settings.AutoStart); err != nil {
		log.Printf("Failed to set startup: %v", err)
	}

	// Notifications
	notifVal := "false"
	if settings.Notifications {
		notifVal = "true"
	}
	a.store.SetSetting("notifications_enabled", notifVal)

	// Port
	if settings.Port != a.proxyServer.GetPort() {
		// Port changed
		// 1. Save to store
		a.store.SetSetting("proxy_port", fmt.Sprintf("%d", settings.Port))

		// 2. Restart Proxy
		if err := a.proxyServer.Restart(settings.Port); err != nil {
			return fmt.Errorf("failed to restart proxy: %w", err)
		}

		// 3. Re-apply system proxy if enabled
		if a.GetProtectionStatus() {
			a.SetSystemProxy(true)
		}
	}

	// Adblock
	a.EnableAdblock(settings.AdblockEnabled)

	return nil
}

// Adblock Filter Management

func (a *App) GetAdblockFilters() []core.AdblockFilter {
	return a.store.GetAdblockFilters()
}

func (a *App) AddAdblockFilter(name, url string) error {
	filter := core.AdblockFilter{
		ID:      utils.GenerateIDString(),
		Name:    name,
		URL:     url,
		Enabled: true,
	}
	err := a.store.AddAdblockFilter(filter)
	if err == nil {
		go a.RefreshAdblockFilters()
	}
	return err
}

func (a *App) DeleteAdblockFilter(id string) error {
	err := a.store.DeleteAdblockFilter(id)
	if err == nil {
		go a.RefreshAdblockFilters()
	}
	return err
}

func (a *App) ToggleAdblockFilter(id string, enabled bool) error {
	filters := a.store.GetAdblockFilters()
	for _, f := range filters {
		if f.ID == id {
			f.Enabled = enabled
			err := a.store.UpdateAdblockFilter(f)
			if err == nil {
				go a.RefreshAdblockFilters()
			}
			return err
		}
	}
	return fmt.Errorf("filter not found")
}

func (a *App) RefreshAdblockFilters() error {
	a.refreshMu.Lock()
	defer a.refreshMu.Unlock()

	filters := a.store.GetAdblockFilters()
	var allRules strings.Builder
	var blocklistSources []string
	// Always include the default hosts list as a base for the blocklist
	blocklistSources = append(blocklistSources, "https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts")

	// Default hardcoded rules for adblock engine
	allRules.WriteString(`||ads.google.com^
||doubleclick.net^
||adnxs.com^
||googleadservices.com^
||pagead2.googlesyndication.com^
||analytics.google.com^
||facebook.com/tr/^
`)

	for _, f := range filters {
		if !f.Enabled {
			continue
		}

		// Add to blocklist sources
		homeDir, _ := os.UserHomeDir()
		filterDir := filepath.Join(homeDir, ".custos", "filters")
		filePath := filepath.Join(filterDir, f.ID+".txt")

		if _, err := os.Stat(filePath); err == nil {
			blocklistSources = append(blocklistSources, filePath)
		} else if f.URL != "" {
			blocklistSources = append(blocklistSources, f.URL)
		}

		content, err := a.getFilterContent(f)
		if err == nil {
			allRules.WriteString("\n")
			allRules.WriteString(content)
		}
	}

	// Update and reload adblock engine
	a.proxyServer.ReloadAdblockEngine(allRules.String())

	// Update and reload blocklist
	if a.blocklist != nil {
		a.blocklist.SetSources(blocklistSources)
		a.blocklist.Load()
	}

	return nil
}

func (a *App) getFilterContent(f core.AdblockFilter) (string, error) {
	homeDir, _ := os.UserHomeDir()
	filterDir := filepath.Join(homeDir, ".custos", "filters")
	os.MkdirAll(filterDir, 0755)
	filePath := filepath.Join(filterDir, f.ID+".txt")

	if f.URL == "" {
		return "", nil
	}

	// Check if file exists and is less than 24h old
	info, err := os.Stat(filePath)
	if err == nil && time.Since(info.ModTime()) < 24*time.Hour {
		content, err := os.ReadFile(filePath)
		if err == nil {
			return string(content), nil
		}
	}

	// Download
	log.Printf("Downloading adblock filter: %s from %s", f.Name, f.URL)
	resp, err := http.Get(f.URL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	os.WriteFile(filePath, content, 0644)

	f.LastUpdated = time.Now()
	a.store.UpdateAdblockFilter(f)

	return string(content), nil
}

func (a *App) seedFilters() {
	// Truncate before seeding as requested
	a.store.ClearAdblockFilters()

	defaults := []struct {
		Name string
		URL  string
	}{
		{"AdGuard DNS", "https://justdomains.github.io/blocklists/lists/adguarddns-justdomains.txt"},
		{"Easy List", "https://justdomains.github.io/blocklists/lists/easylist-justdomains.txt"},
		{"Easy Privacy", "https://justdomains.github.io/blocklists/lists/easyprivacy-justdomains.txt"},
		{"NoCoin", "https://justdomains.github.io/blocklists/lists/nocoin-justdomains.txt"},
		{"Youtube Clear View", "https://github.com/yokoffing/filterlists/blob/main/youtube_clear_view.txt"},
		{"Privacy Essentials", "https://github.com/yokoffing/filterlists/blob/main/privacy_essentials.txt"},
		{"Abblock Pro Mini", "https://github.com/hagezi/dns-blocklists/blob/main/adblock/pro.mini.txt"},
		{"Pi-hole", "https://raw.githubusercontent.com/xxcriticxx/.pl-host-file/master/hosts.txt"},
		{"Ramnit", "https://1275.ru/DGA/ramnit.txt"},
		{"SharkBot", "https://1275.ru/DGA/sharkbot.txt"},
		{"QSnatch", "https://1275.ru/DGA/qsnatch.txt"},
		{"CryptoLocker", "https://1275.ru/DGA/cryptolocker.txt"},
		{"1024 Hosts", "https://raw.githubusercontent.com/Goooler/1024_hosts/master/hosts"},
	}

	existingFilters := a.store.GetAdblockFilters()
	existingURLs := make(map[string]bool)
	for _, f := range existingFilters {
		existingURLs[f.URL] = true
	}

	added := false
	for _, d := range defaults {
		if !existingURLs[d.URL] {
			filter := core.AdblockFilter{
				ID:      utils.GenerateIDString(),
				Name:    d.Name,
				URL:     d.URL,
				Enabled: true,
			}
			if err := a.store.AddAdblockFilter(filter); err == nil {
				added = true
			}
		}
	}

	if added {
		go a.RefreshAdblockFilters()
	}
}
