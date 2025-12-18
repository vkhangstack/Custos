package store

import (
	"fmt"
	"sync"
	"time"

	"github.com/vkhangstack/Custos/internal/core"
	"github.com/vkhangstack/Custos/internal/utils"
)

// MemoryStore holds logs and stats in memory
type MemoryStore struct {
	mu          sync.RWMutex
	logs        []core.LogEntry
	maxLogs     int
	stats       core.Stats
	subscribers []func(core.LogEntry)
}

// NewMemoryStore creates a new store
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		logs:    make([]core.LogEntry, 0, 1000),
		maxLogs: 1000,
		stats: core.Stats{
			TopDomains: make(map[string]int64),
		},
	}
}

// AddLog adds a log entry and updates stats
func (s *MemoryStore) AddLog(entry core.LogEntry) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Ensure ID is set
	if entry.ID == "" {
		entry.ID = utils.GenerateIDString()
	}

	// Add to logs (circular buffer logic simplified)
	if len(s.logs) >= s.maxLogs {
		// Remove oldest
		s.logs = s.logs[1:]
	}
	s.logs = append(s.logs, entry)

	// Update stats
	s.stats.TotalUpload += entry.BytesSent
	s.stats.TotalDownload += entry.BytesRecv
	if entry.Domain != "" {
		s.stats.TopDomains[entry.Domain] += (entry.BytesSent + entry.BytesRecv)
	}

	// Notify subscribers (simple broadcast)
	for _, sub := range s.subscribers {
		// Run in goroutine to avoid blocking
		go sub(entry)
	}
}

// UpdateLog updates an existing log entry (naive implementation: linear search)
// For MemoryStore, this might be slow, but it's for dev/fallback.
func (s *MemoryStore) UpdateLog(entry core.LogEntry) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for i, log := range s.logs {
		if log.ID == entry.ID {
			// Merge fields from entry into s.logs[i]
			if !entry.Timestamp.IsZero() {
				s.logs[i].Timestamp = entry.Timestamp
			}
			if entry.Type != "" {
				s.logs[i].Type = entry.Type
			}
			if entry.Domain != "" {
				s.logs[i].Domain = entry.Domain
			}
			if entry.SrcIP != "" {
				s.logs[i].SrcIP = entry.SrcIP
			}
			if entry.DstIP != "" {
				s.logs[i].DstIP = entry.DstIP
			}
			if entry.DstPort != 0 {
				s.logs[i].DstPort = entry.DstPort
			}
			if entry.Protocol != "" {
				s.logs[i].Protocol = entry.Protocol
			}
			if entry.ProcessName != "" {
				s.logs[i].ProcessName = entry.ProcessName
			}
			if entry.ProcessID != 0 {
				s.logs[i].ProcessID = entry.ProcessID
			}
			if entry.BytesSent != 0 {
				s.logs[i].BytesSent = entry.BytesSent
			}
			if entry.BytesRecv != 0 {
				s.logs[i].BytesRecv = entry.BytesRecv
			}
			if entry.Status != "" {
				s.logs[i].Status = entry.Status
			}
			if entry.Latency != 0 {
				s.logs[i].Latency = entry.Latency
			}

			// Notify subscribers of update with full entry
			updatedEntry := s.logs[i]
			for _, sub := range s.subscribers {
				go sub(updatedEntry)
			}
			return
		}
	}
}

// AddTraffic increments traffic stats
func (s *MemoryStore) AddTraffic(upload, download int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stats.TotalUpload += upload
	s.stats.TotalDownload += download
}

// GetRecentLogs returns the last N logs
func (s *MemoryStore) GetRecentLogs(limit int) []core.LogEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if limit > len(s.logs) {
		limit = len(s.logs)
	}

	// Return a copy to be safe
	result := make([]core.LogEntry, limit)
	copy(result, s.logs[len(s.logs)-limit:])

	// Reverse to show newest first? Frontend usually handles sorting, but let's just return raw
	return result
}

// GetStats returns current stats
func (s *MemoryStore) GetStats() core.Stats {
	s.mu.RLock()
	defer s.mu.RUnlock()
	stats := s.stats
	stats.Timestamp = time.Now()
	return stats
}

// GetTrafficHistory returns (empty/dummy) traffic history for memory store
func (s *MemoryStore) GetTrafficHistory(duration time.Duration) []core.TrafficDataPoint {
	// Simple implementation: Aggregate recent logs manually?
	// For now, return empty or mock.
	return []core.TrafficDataPoint{}
}

// Subscribe adds a listener for new logs
func (s *MemoryStore) Subscribe(callback func(core.LogEntry)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.subscribers = append(s.subscribers, callback)
}

// ResetData resets data
func (s *MemoryStore) ResetData() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stats = core.Stats{
		TopDomains: make(map[string]int64),
	}
}

func (s *MemoryStore) AddRule(rule core.Rule) error { return nil }
func (s *MemoryStore) GetRules() []core.Rule        { return nil }
func (s *MemoryStore) GetRulesPaginated(page, pageSize int, search string) ([]core.Rule, int64, error) {
	return []core.Rule{}, 0, nil
}

func (s *MemoryStore) GetSetting(key string) (string, error) {
	return "", fmt.Errorf("not implemented")
}

func (s *MemoryStore) SetSetting(key, value string) error {
	return nil
}
func (s *MemoryStore) DeleteRule(id string) error      { return nil }
func (s *MemoryStore) UpdateRule(rule core.Rule) error { return nil }

func (s *MemoryStore) AddAdblockFilter(filter core.AdblockFilter) error    { return nil }
func (s *MemoryStore) GetAdblockFilters() []core.AdblockFilter             { return nil }
func (s *MemoryStore) DeleteAdblockFilter(id string) error                 { return nil }
func (s *MemoryStore) UpdateAdblockFilter(filter core.AdblockFilter) error { return nil }
func (s *MemoryStore) ClearAdblockFilters() error                          { return nil }

func (s *MemoryStore) IncrementRuleHit(id string, domain string) error {
	// Not fully implemented for rules in MemoryStore yet as MemoryStore
	// doesn't actually store/manage rules in the current implementation.
	return nil
}

func (s *MemoryStore) IncrementAdblockHit(domain string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.stats.AdblockHits++
	return nil
}
