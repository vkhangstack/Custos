package store

import (
	"time"

	"github.com/vkhangstack/Custos/internal/core"
)

// Store defines the interface for data persistence
type Store interface {
	AddLog(entry core.LogEntry)
	UpdateLog(entry core.LogEntry)
	AddTraffic(upload, download int64)
	GetTrafficHistory(duration time.Duration) []core.TrafficDataPoint
	GetRecentLogs(limit int) []core.LogEntry
	GetLogsPaginated(cursor string, limit int, search, status, logType string) ([]core.LogEntry, string, bool, int64, error)
	GetStats() core.Stats
	Subscribe(callback func(core.LogEntry))
	ResetData()
	// Rule Management
	AddRule(rule core.Rule) error
	GetRules() []core.Rule
	GetRulesPaginated(page, pageSize int, search string) ([]core.Rule, int64, error)
	DeleteRule(id string) error
	UpdateRule(rule core.Rule) error
	IncrementRuleHit(id string, domain string) error
	IncrementAdblockHit(domain string) error
	// Adblock Filters
	AddAdblockFilter(filter core.AdblockFilter) error
	GetAdblockFilters() []core.AdblockFilter
	DeleteAdblockFilter(id string) error
	UpdateAdblockFilter(filter core.AdblockFilter) error
	ClearAdblockFilters() error
	// Settings
	GetSetting(key string) (string, error)
	SetSetting(key, value string) error
}

type Cache interface {
	Add(key string, value interface{}) error
	Get(key string) (interface{}, error)
	Delete(key string) error
	Clear()
}
