package core

import (
	"bufio"
	"net/http"
	"strings"
	"sync"
	"time"
)

// BlocklistManager handles the loading and checking of blocked domains
type BlocklistManager struct {
	mu      sync.RWMutex
	domains map[string]bool
	sources []string
}

// NewBlocklistManager creates a new manager
func NewBlocklistManager() *BlocklistManager {
	return &BlocklistManager{
		domains: make(map[string]bool),
		sources: []string{
			"https://raw.githubusercontent.com/StevenBlack/hosts/master/hosts",
		},
	}
}

// Load loads all configured sources
func (m *BlocklistManager) Load() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Clear existing
	m.domains = make(map[string]bool)

	for _, source := range m.sources {
		if err := m.loadSource(source); err != nil {
			// Log error but continue
			continue
		}
	}
	return nil
}

func (m *BlocklistManager) loadSource(url string) error {
	// Create a client that explicitly bypasses system proxy
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: nil, // Bypass system proxy
		},
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	// Force connection close to avoid "Unsolicited response received on idle HTTP channel"
	req.Close = true

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		// Hosts file format: 0.0.0.0 domain.com
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			domain := parts[1]
			// Simple validation
			if !strings.Contains(domain, ".") {
				continue
			}
			m.domains[domain] = true
		}
	}
	return scanner.Err()
}

// IsBlocked checks if a domain is blocked
func (m *BlocklistManager) IsBlocked(domain string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// Remove trailing dot if present (DNS validity)
	domain = strings.TrimSuffix(domain, ".")

	return m.domains[domain]
}

// Count returns the number of blocked domains
func (m *BlocklistManager) Count() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.domains)
}
