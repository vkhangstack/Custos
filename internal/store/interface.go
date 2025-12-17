package store

import (
	"Custos/internal/core"
	"time"
)

// Store defines the interface for data persistence
type Store interface {
	AddLog(entry core.LogEntry)
	UpdateLog(entry core.LogEntry)
	AddTraffic(upload, download int64)
	GetTrafficHistory(duration time.Duration) []core.TrafficDataPoint
	GetRecentLogs(limit int) []core.LogEntry
	GetStats() core.Stats
	Subscribe(callback func(core.LogEntry))
	ResetStats()
	// Rule Management
	AddRule(rule core.Rule) error
	GetRules() []core.Rule
	GetRulesPaginated(page, pageSize int, search string) ([]core.Rule, int64, error)
	DeleteRule(id string) error
	UpdateRule(rule core.Rule) error
	// Settings
	GetSetting(key string) (string, error)
	SetSetting(key, value string) error
}
