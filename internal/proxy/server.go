package proxy

import (
	"Custos/internal/core"
	"Custos/internal/store"
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/armon/go-socks5"
)

// Server manages the proxy listeners
// Server manages the proxy listeners
// Server manages the proxy listeners
type Server struct {
	store             store.Store
	blocklist         *core.BlocklistManager
	socksServer       *socks5.Server
	port              int
	running           bool
	listener          net.Listener
	protectionEnabled bool
}

// NewServer creates a new proxy server
func NewServer(store store.Store, blocklist *core.BlocklistManager, port int) *Server {
	return &Server{
		store:     store,
		blocklist: blocklist,
		port:      port,
	}
}

// SetProtection enables or disables the HTTP protection
func (s *Server) SetProtection(enabled bool) {
	s.protectionEnabled = enabled
	log.Printf("Protection mode set to: %v", enabled)
}

const logIDKey = "logID"

// Start starts the SOCKS5 proxy
func (s *Server) Start() error {
	conf := &socks5.Config{
		Logger: log.New(log.Writer(), "[SOCKS5] ", log.LstdFlags),
		Rules:  &LoggingRuleSet{store: s.store, blocklist: s.blocklist, server: s},
		Dial: func(ctx context.Context, network, addr string) (net.Conn, error) {
			// Extract ID early to update status on failure
			logID, hasLogID := ctx.Value(logIDKey).(string)

			// Dial upstream
			conn, err := net.Dial(network, addr)
			if err != nil {
				if hasLogID {
					// Update log status to error
					s.store.UpdateLog(core.LogEntry{
						ID:     logID,
						Status: "connection_failed",
					})
				}
				return nil, err
			}

			// Wrap if we have a logID
			if hasLogID {
				fmt.Printf("[DEBUG] Dialing for logID: %s\n", logID)
				return &CountingConn{
					Conn:  conn,
					logID: logID,
					entry: core.LogEntry{
						ID: logID,
					},
					store: s.store,
				}, nil
			} else {
				fmt.Printf("[DEBUG] No logID in Dial context!\n")
			}
			return conn, nil
		},
	}
	// ...

	server, err := socks5.New(conf)
	if err != nil {
		return err
	}
	s.socksServer = server

	addr := fmt.Sprintf(":%d", s.port)
	s.listener, err = net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	s.running = true
	log.Printf("SOCKS5 Proxy started on %s", addr)

	go func() {
		if err := s.socksServer.Serve(s.listener); err != nil && s.running {
			log.Printf("SOCKS5 server error: %v", err)
		}
	}()

	return nil
}

// Stop stops the proxy
func (s *Server) Stop() {
	s.running = false
	if s.listener != nil {
		s.listener.Close()
	}
}

// LoggingRuleSet intercepts requests for logging
type LoggingRuleSet struct {
	store     store.Store
	blocklist *core.BlocklistManager
	server    *Server
}

func (r *LoggingRuleSet) Allow(ctx context.Context, req *socks5.Request) (context.Context, bool) {
	domain := req.DestAddr.FQDN

	// Whitelist Localhost/Loopback
	// Always allow local traffic to bypass protection and blocks
	if domain == "localhost" || req.DestAddr.IP.IsLoopback() {
		return ctx, true
	}

	// Check Protection Mode (Block HTTP Port 80)
	if r.server.protectionEnabled {
		if req.DestAddr.Port == 80 {
			r.logBlock(req, domain, "protection_http_blocked")
			log.Printf("Blocked HTTP request to %s (Port 80) due to active protection", domain)
			return ctx, false
		}
	}

	// Check Blocklist
	if r.blocklist.IsBlocked(domain) {
		r.logBlock(req, domain, "blocklist")
		return ctx, false
	}

	// Check Custom Rules
	// Optimized: Could cache this or use a more efficient matcher
	rules := r.store.GetRules()
	for _, rule := range rules {
		if !rule.Enabled {
			continue
		}
		// Simple wildcard matching or substring?
		// User said "pattern". Let's assume glob or exact.
		// For MVP, let's support exact match or domain suffix.
		// If rule.Pattern == domain ...
		// Go's filepath.Match is good for globs.
		if matched, _ := matchDomain(rule.Pattern, domain); matched {
			if rule.Type == core.RuleBlock {
				r.logBlock(req, domain, "custom_rule")
				return ctx, false
			}
			// If ALLOW, we stop checking other block rules?
			// Typically whitelist overrides blacklist.
			// But here we just proceed.
			break
		}
	}

	// ... continue to allow

	// Log the connection attempt
	// Define variables
	destIP := req.DestAddr.IP.String()
	procName := "custos" // Placeholder, would need SystemTracker integration
	procID := int32(0)

	entry := core.LogEntry{
		ID:          fmt.Sprintf("%d", time.Now().UnixNano()), // Simple ID
		Timestamp:   time.Now(),
		Type:        "proxy",
		Domain:      domain, // Could be empty if IP
		DstIP:       destIP,
		DstPort:     req.DestAddr.Port,
		Protocol:    "tcp", // SOCKS5 is usually TCP
		ProcessName: procName,
		ProcessID:   procID,
		Status:      "allowed",
	}

	r.store.AddLog(entry)

	// Inject logID into context for Dial to pick up
	return context.WithValue(ctx, logIDKey, entry.ID), true
}

// matchDomain checks if domain matches pattern
// Pattern support: *.google.com, google.com (exact)
func matchDomain(pattern, domain string) (bool, error) {
	// Simple implementation
	// If pattern starts with *., match suffix
	if len(pattern) > 2 && pattern[:2] == "*." {
		suffix := pattern[2:]
		if len(domain) >= len(suffix) && domain[len(domain)-len(suffix):] == suffix {
			return true, nil
		}
	}
	return pattern == domain, nil
}

func (r *LoggingRuleSet) logBlock(req *socks5.Request, domain string, reason string) {
	entry := core.LogEntry{
		ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
		Timestamp: time.Now(),
		Type:      "proxy",
		DstIP:     req.DestAddr.IP.String(),
		DstPort:   req.DestAddr.Port,
		Domain:    domain,
		Protocol:  "tcp",
		Status:    "blocked", // or "blocked:" + reason
		BytesSent: 0,
		BytesRecv: 0,
	}
	r.store.AddLog(entry)
}
