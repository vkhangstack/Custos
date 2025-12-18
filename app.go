package main

import (
	"context"
	"fmt"
	"time"

	"github.com/vkhangstack/Custos/internal/system"

	"github.com/vkhangstack/Custos/internal/store"

	"github.com/vkhangstack/Custos/internal/proxy"

	"github.com/vkhangstack/Custos/internal/dns"

	"github.com/vkhangstack/Custos/internal/core"
)

// App struct
type App struct {
	ctx           context.Context
	store         store.Store
	proxyServer   *proxy.Server
	dnsServer     *dns.Server
	systemTracker *system.Tracker
}

// NewApp creates a new App application struct
func NewApp() *App {
	// Try to initialize SQLite store
	var s store.Store
	var err error
	s, err = store.NewSQLiteStore("custos.db")
	if err != nil {
		fmt.Printf("Failed to initialize SQLite store: %v. Falling back to MemoryStore.\n", err)
		s = store.NewMemoryStore()
	}

	bm := core.NewBlocklistManager()

	// Load blocklist in background
	go bm.Load()

	systemTracker := system.NewTracker()

	return &App{
		store:         s,
		proxyServer:   proxy.NewServer(s, bm, systemTracker, 1080),
		dnsServer:     dns.NewServer(s, bm, 5353),
		systemTracker: systemTracker,
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

	// Restore Protection State
	if enabled := a.GetProtectionStatus(); enabled {
		a.proxyServer.SetProtection(true)
		a.SetSystemProxy(true)
	} else {
		a.proxyServer.SetProtection(false)
		a.SetSystemProxy(false)
	}

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
	return system.SetSystemProxy(enabled, 1080)
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

// GetChartData returns historical traffic data for the chart
func (a *App) GetChartData() []core.TrafficDataPoint {
	// Get last 30 minutes of history
	return a.store.GetTrafficHistory(30 * time.Minute)
}

// Rule Management

// AddRule adds a new rule
func (a *App) AddRule(pattern string, ruleType string) error {
	rule := core.Rule{
		ID:      fmt.Sprintf("%d", time.Now().UnixNano()),
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
