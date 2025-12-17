package core

import "time"

// LogEntry represents a single traffic event
type LogEntry struct {
	ID          string    `json:"id"`
	Timestamp   time.Time `json:"timestamp"`
	Type        string    `json:"type"` // "dns", "proxy", "system"
	Domain      string    `json:"domain"`
	SrcIP       string    `json:"src_ip"`
	DstIP       string    `json:"dst_ip"`
	DstPort     int       `json:"dst_port"`
	Protocol    string    `json:"protocol"` // "tcp", "udp", "http", "https"
	ProcessName string    `json:"process_name"`
	ProcessID   int32     `json:"process_id"`
	BytesSent   int64     `json:"bytes_sent"`
	BytesRecv   int64     `json:"bytes_recv"`
	Status      string    `json:"status"`  // "allowed", "blocked", "error"
	Latency     int64     `json:"latency"` // in ms
}

// Stats represents aggregated statistics
type Stats struct {
	TotalUpload   int64            `json:"total_upload"`
	TotalDownload int64            `json:"total_download"`
	ActiveConns   int              `json:"active_connections"`
	TopDomains    map[string]int64 `json:"top_domains"`
}

// TrafficDataPoint represents a point in the traffic chart
type TrafficDataPoint struct {
	Name     string `json:"name"`     // Time label (e.g., "10:00")
	Upload   int64  `json:"upload"`   // Bytes sent
	Download int64  `json:"download"` // Bytes received
}

// RuleType defines blocking or allowing
type RuleType string

const (
	RuleBlock RuleType = "BLOCK"
	RuleAllow RuleType = "ALLOW"
)

const (
	RuleSourceDefault RuleType = "default"
	RuleSourceCustom  RuleType = "custom"
)

// Rule represents a filtering rule
type Rule struct {
	ID      string   `json:"id"`
	Type    RuleType `json:"type"`
	Pattern string   `json:"pattern"` // e.g., "*.ads.com"
	Enabled bool     `json:"enabled"`
	Source  RuleType `json:"source"` // "custom" or "default"
}

// TrafficStatsModel is the DB model for persistent stats
type TrafficStatsModel struct {
	ID            string `gorm:"primaryKey"`
	TotalUpload   int64
	TotalDownload int64
}

// PaginatedRulesResponse wraps rules and total count
type PaginatedRulesResponse struct {
	Rules []Rule `json:"rules"`
	Total int64  `json:"total"`
}

// AppSetting represents a key-value setting
type AppSetting struct {
	Key   string `gorm:"primaryKey" json:"key"`
	Value string `json:"value"`
}
