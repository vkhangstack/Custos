package main

import (
	"Custos/internal/core"
	"Custos/internal/dns"
	"Custos/internal/proxy"
	"Custos/internal/store"
	"Custos/internal/system"
	"context"
	"fmt"
	"time"
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

	return &App{
		store:         s,
		proxyServer:   proxy.NewServer(s, bm, 1080),
		dnsServer:     dns.NewServer(s, bm, 5353),
		systemTracker: system.NewTracker(),
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
	if err := a.SetSystemProxy(true); err != nil {
		fmt.Printf("Failed to set system proxy on startup: %v\n", err)
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
	a.store.ResetStats()
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

// GetChartData returns historical traffic data for the chart
func (a *App) GetChartData() []core.TrafficDataPoint {
	// Get last 30 minutes of history
	return a.store.GetTrafficHistory(30 * time.Minute)
}
