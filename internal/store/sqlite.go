package store

import (
	"Custos/internal/core"
	"bufio"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SQLiteStore persists logs to a database
type SQLiteStore struct {
	db          *gorm.DB
	subscribers []func(core.LogEntry)
	mu          sync.RWMutex

	// Rule Cache
	cachedRules []core.Rule
	rulesLoaded bool
	cacheMu     sync.RWMutex
}

// NewSQLiteStore creates a new persistent store
func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	// Auto-migrate schema
	if err := db.AutoMigrate(&core.LogEntry{}, &core.TrafficStatsModel{}, &core.Rule{}, &core.AppSetting{}); err != nil {
		return nil, err
	}

	// Initialize Stats if not present or resync
	var stats core.TrafficStatsModel
	if err := db.First(&stats, "id = ?", "global").Error; err != nil {
		// Not found, create it with sum from logs (Migration logic)
		var result struct {
			Upload   int64
			Download int64
		}
		db.Model(&core.LogEntry{}).Select("ifnull(sum(bytes_sent),0) as upload, ifnull(sum(bytes_recv),0) as download").Scan(&result)

		stats = core.TrafficStatsModel{
			ID:            "global",
			TotalUpload:   result.Upload,
			TotalDownload: result.Download,
		}
		db.Create(&stats)
		log.Printf("Initialized global stats from logs: Up=%d, Down=%d", stats.TotalUpload, stats.TotalDownload)
	}

	// Seed Default Rules if empty
	s := &SQLiteStore{
		db: db,
	}

	go func() {
		s.seedDefaultRules()
	}()

	return s, nil
}

// Subscribe adds a listener
func (s *SQLiteStore) Subscribe(callback func(core.LogEntry)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.subscribers = append(s.subscribers, callback)
}

// AddLog adds a log entry
func (s *SQLiteStore) AddLog(entry core.LogEntry) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Ensure ID is set
	if entry.ID == "" {
		entry.ID = time.Now().String() // Simple fallback
	}

	// Use Transaction to ensure Log and Stats are in sync
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 1. Persist Log
		if err := tx.Create(&entry).Error; err != nil {
			return err
		}

		// 2. Update Stats Atomic
		if entry.BytesSent > 0 || entry.BytesRecv > 0 {
			if err := tx.Model(&core.TrafficStatsModel{}).Where("id = ?", "global").
				Updates(map[string]interface{}{
					"total_upload":   gorm.Expr("total_upload + ?", entry.BytesSent),
					"total_download": gorm.Expr("total_download + ?", entry.BytesRecv),
				}).Error; err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		log.Printf("Failed to save log and stats: %v", err)
	}

	// Notify subscribers
	for _, sub := range s.subscribers {
		go sub(entry)
	}
}

// UpdateLog updates an existing log entry
func (s *SQLiteStore) UpdateLog(entry core.LogEntry) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Update only the non-zero fields (GORM Updates behavior) to avoid overwriting existing data with empty values
	if err := s.db.Model(&core.LogEntry{}).Where("id = ?", entry.ID).Updates(entry).Error; err != nil {
		log.Printf("Failed to update log: %v", err)
	}

	// Notify subscribers
	// Note: entry here might be partial. Subscribers might need full entry.
	// For now, we broadcast the partial entry. Frontend might merge or replace.
	// If frontend replaces, we lose data in UI until refresh.
	// Ideally we should fetch updated entry or merge before broadcast.
	// But let's stick to simple for now.
	for _, sub := range s.subscribers {
		go sub(entry)
	}
}

// AddTraffic increments persistent traffic stats explicitly
func (s *SQLiteStore) AddTraffic(upload, download int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if upload == 0 && download == 0 {
		return
	}

	if err := s.db.Model(&core.TrafficStatsModel{}).Where("id = ?", "global").
		Updates(map[string]interface{}{
			"total_upload":   gorm.Expr("total_upload + ?", upload),
			"total_download": gorm.Expr("total_download + ?", download),
		}).Error; err != nil {
		log.Printf("Failed to update traffic stats: %v", err)
	}
}

// GetRecentLogs returns the last N logs
func (s *SQLiteStore) GetRecentLogs(limit int) []core.LogEntry {
	var logs []core.LogEntry
	s.db.Order("timestamp desc").Limit(limit).Find(&logs)
	return logs
}

// GetStats calculates stats from DB
func (s *SQLiteStore) GetStats() core.Stats {
	var stats core.Stats

	// Get Totals from TrafficStatsModel (Fast)
	var model core.TrafficStatsModel
	if err := s.db.First(&model, "id = ?", "global").Error; err == nil {
		stats.TotalUpload = model.TotalUpload
		stats.TotalDownload = model.TotalDownload
	} else {
		// Fallback to sum (shouldn't happen if init worked)
		var res struct {
			TotalUpload   int64
			TotalDownload int64
		}
		s.db.Model(&core.LogEntry{}).Select("sum(bytes_sent) as total_upload, sum(bytes_recv) as total_download").Scan(&res)
		stats.TotalUpload = res.TotalUpload
		stats.TotalDownload = res.TotalDownload
	}

	// Top Domains (Last 1000 logs for performance)
	// Fallback to DstIP if Domain is empty
	rows, _ := s.db.Model(&core.LogEntry{}).
		Select("ifnull(nullif(domain, ''), dst_ip) as target, sum(bytes_sent + bytes_recv) as total").
		Group("target").
		Order("total desc").
		Limit(5).
		Rows()
	defer rows.Close()

	stats.TopDomains = make(map[string]int64)
	for rows.Next() {
		var target string
		var total int64
		rows.Scan(&target, &total)
		if target != "" {
			stats.TopDomains[target] = total
		}
	}

	return stats
}

// GetTrafficHistory returns traffic stats grouped by time
func (s *SQLiteStore) GetTrafficHistory(duration time.Duration) []core.TrafficDataPoint {
	var points []core.TrafficDataPoint

	// SQLite specific query to group by minute (or rough interval)
	threshold := time.Now().Add(-duration)

	// Group by %H:%M using Local Time
	rows, err := s.db.Model(&core.LogEntry{}).
		Select("strftime('%H:%M', timestamp, 'localtime') as name, sum(bytes_sent) as upload, sum(bytes_recv) as download").
		Where("timestamp > ?", threshold).
		Group("name").
		Order("name ASC").
		Rows()

	if err != nil {
		log.Printf("Failed to get traffic history: %v", err)
		return points
	}
	defer rows.Close()

	for rows.Next() {
		var p core.TrafficDataPoint
		rows.Scan(&p.Name, &p.Upload, &p.Download)
		points = append(points, p)
	}

	return points
}

// Subscribe adds a listener
func (s *SQLiteStore) ResetStats() {
	// Use Exec for direct deletion to bypass GORM's global delete protection if enabled
	// and to ensure efficient clearing.
	s.db.Exec("DELETE FROM traffic_stats_models WHERE id = ?", "global")
	s.db.Exec("DELETE FROM log_entries")
	s.db.Exec("DELETE FROM traffic_stats_models")
	s.db.Exec("DELETE FROM stats")
}

// Rule Management Implementation

func (s *SQLiteStore) AddRule(rule core.Rule) error {
	if err := s.db.Create(&rule).Error; err != nil {
		return err
	}
	s.invalidateCache()
	return nil
}

func (s *SQLiteStore) GetRules() []core.Rule {
	s.cacheMu.RLock()
	if s.rulesLoaded {
		defer s.cacheMu.RUnlock()
		return s.cachedRules
	}
	s.cacheMu.RUnlock()

	s.cacheMu.Lock()
	defer s.cacheMu.Unlock()

	// Double check
	if s.rulesLoaded {
		return s.cachedRules
	}

	var rules []core.Rule
	s.db.Find(&rules)
	s.cachedRules = rules
	s.rulesLoaded = true
	return rules
}

func (s *SQLiteStore) GetRulesPaginated(page, pageSize int, search string) ([]core.Rule, int64, error) {
	var rules []core.Rule
	var total int64

	offset := (page - 1) * pageSize
	query := s.db.Model(&core.Rule{})

	if search != "" {
		likePattern := "%" + search + "%"
		query = query.Where("pattern LIKE ?", likePattern)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Sort by Source ASC (Custom first), then ID DESC (Newest first)
	if err := query.Order("source asc, id desc").Offset(offset).Limit(pageSize).Find(&rules).Error; err != nil {
		return nil, 0, err
	}

	return rules, total, nil
}

func (s *SQLiteStore) DeleteRule(id string) error {
	if err := s.db.Delete(&core.Rule{}, "id = ?", id).Error; err != nil {
		return err
	}
	s.invalidateCache()
	return nil
}

func (s *SQLiteStore) UpdateRule(rule core.Rule) error {
	// GORM Updates with struct ignores zero values (false, 0, "").
	// We use map to ensure Enabled=false is applied, but we must only include
	// other fields if they are set (to avoid overwriting with empty strings during toggle).
	updates := map[string]interface{}{
		"enabled": rule.Enabled,
	}

	if rule.Pattern != "" {
		updates["pattern"] = rule.Pattern
	}
	if rule.Type != "" {
		updates["type"] = rule.Type
	}
	if rule.Source != "" {
		updates["source"] = rule.Source
	}

	err := s.db.Model(&core.Rule{}).Where("id = ?", rule.ID).Updates(updates).Error

	if err != nil {
		return err
	}
	s.invalidateCache()
	return nil
}

func (s *SQLiteStore) invalidateCache() {
	s.cacheMu.Lock()
	defer s.cacheMu.Unlock()
	s.rulesLoaded = false
}

// Settings
func (s *SQLiteStore) GetSetting(key string) (string, error) {
	var setting core.AppSetting
	if err := s.db.First(&setting, "key = ?", key).Error; err != nil {
		return "", err
	}
	return setting.Value, nil
}

func (s *SQLiteStore) SetSetting(key, value string) error {
	setting := core.AppSetting{Key: key, Value: value}
	// Upsert
	return s.db.Save(&setting).Error
}

// seedDefaultRules fetches and applies a default blocklist
func (s *SQLiteStore) seedDefaultRules() {
	url := "https://adaway.org/hosts.txt"
	log.Println("Seeding default rules from", url)

	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Failed to fetch blocklist: %v", err)
		return
	}
	defer resp.Body.Close()

	var rules []core.Rule
	scanner := bufio.NewScanner(resp.Body)

	// Fetch existing default rules to avoid duplicates
	existingRules := make(map[string]bool)
	var currentDefaults []core.Rule
	// Optimize: Select only pattern for default rules
	s.db.Model(&core.Rule{}).Where("source = ?", core.RuleSourceDefault).Select("pattern").Find(&currentDefaults)
	for _, r := range currentDefaults {
		existingRules[r.Pattern] = true
	}

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			domain := parts[1]
			if domain == "localhost" || domain == "broadcasthost" || domain == "::1" {
				continue
			}

			// Deduplicate against existing DB rules
			if existingRules[domain] {
				continue
			}

			// Also deduplicate within the new batch (hosts file might have dupes)
			existingRules[domain] = true

			rules = append(rules, core.Rule{
				ID:      fmt.Sprintf("seed-%d", time.Now().UnixNano()+int64(len(rules))),
				Type:    core.RuleBlock,
				Pattern: domain,
				Enabled: true,
				Source:  core.RuleSourceDefault,
			})
		}
	}

	if len(rules) > 0 {
		// Batch insert in chunks to avoid SQLite limits
		chunkSize := 500
		for i := 0; i < len(rules); i += chunkSize {
			end := i + chunkSize
			if end > len(rules) {
				end = len(rules)
			}
			log.Println("Seeding chunk", i, "to", end)
			if err := s.db.Create(rules[i:end]).Error; err != nil {
				log.Printf("Failed to seed chunk: %v", err)
			}
		}
		log.Printf("Seeded %d new rules", len(rules))
		s.invalidateCache()
	} else {
		log.Println("No new rules to seed.")
	}

	// Invalidate cache so next GetRules fetches them
	s.invalidateCache()
}

func (s *SQLiteStore) TruncateRules() error {
	return s.db.Exec("DELETE FROM rules").Error
}
