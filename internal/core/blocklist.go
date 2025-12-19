package core

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
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
		sources: []string{}, // Start empty, will be seeded/populated by App
	}
}

// SetSources updates the sources list
func (m *BlocklistManager) SetSources(sources []string) {
	m.mu.Lock()
	m.sources = sources
	m.mu.Unlock()
}

// Load loads all configured sources
func (m *BlocklistManager) Load() error {
	m.mu.RLock()
	sources := make([]string, len(m.sources))
	copy(sources, m.sources)
	m.mu.RUnlock()

	newDomains := make(map[string]bool)

	for _, source := range sources {
		fmt.Println("Loading blocklist source:", source)
		if err := m.loadSource(source, newDomains); err != nil {
			// Log error but continue
			continue
		}
	}

	m.mu.Lock()
	m.domains = newDomains
	m.mu.Unlock()

	return nil
}

func (m *BlocklistManager) loadSource(source string, domains map[string]bool) error {
	if source == "" {
		return nil
	}

	var body io.ReadCloser
	if strings.HasPrefix(source, "http://") || strings.HasPrefix(source, "https://") {
		client := &http.Client{
			Transport: &http.Transport{
				Proxy: nil, // Bypass system proxy
			},
			Timeout: 30 * time.Second,
		}

		req, err := http.NewRequest("GET", source, nil)
		if err != nil {
			return err
		}
		req.Close = true

		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		body = resp.Body
	} else {
		// Assume local file
		f, err := os.Open(source)
		if err != nil {
			return fmt.Errorf("failed to open local source %s: %v", source, err)
		}
		body = f
	}
	defer body.Close()

	count := 0
	scanner := bufio.NewScanner(body)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		parts := strings.Fields(line)
		var domain string

		if len(parts) >= 2 {
			// Hosts file format: 0.0.0.0 domain.com
			domain = parts[1]
		} else if len(parts) == 1 {
			// Simple list format: domain.com
			domain = parts[0]
		}

		if domain != "" {
			// Simple validation
			if strings.Contains(domain, ".") {
				domains[domain] = true
				count++
			}
		}
	}
	fmt.Printf("Loaded %d domains from source: %s\n", count, source)
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
