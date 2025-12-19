package proxy

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"github.com/vkhangstack/Custos/internal/adblock"
	"github.com/vkhangstack/Custos/internal/core"
	"github.com/vkhangstack/Custos/internal/store"
	"github.com/vkhangstack/Custos/internal/system"
	"github.com/vkhangstack/Custos/internal/utils"

	"github.com/armon/go-socks5"
)

// Server manages the proxy listeners
// Server manages the proxy listeners
// Server manages the proxy listeners
type Server struct {
	store             store.Store
	blocklist         *core.BlocklistManager
	systemTracker     *system.Tracker
	socksServer       *socks5.Server
	port              int
	running           bool
	listener          net.Listener
	protectionEnabled bool
	adblockEnabled    bool
	adblockEngine     *adblock.Engine
	mu                sync.RWMutex
}

// NewServer creates a new proxy server
func NewServer(store store.Store, blocklist *core.BlocklistManager, systemTracker *system.Tracker, port int) *Server {
	// Initialize adblock engine with default rules for now
	// In the future, this can be loaded from DB or files
	adblockRules := `||ads.google.com^
||doubleclick.net^
||adnxs.com^
||googleadservices.com^
||pagead2.googlesyndication.com^
||analytics.google.com^
||facebook.com/tr/^
`
	engine := adblock.NewEngine(adblockRules)

	return &Server{
		store:             store,
		blocklist:         blocklist,
		systemTracker:     systemTracker,
		port:              port,
		adblockEngine:     engine,
		protectionEnabled: true, // Default
	}
}

// SetProtection enables or disables the HTTP protection
func (s *Server) SetProtection(enabled bool) {
	s.mu.Lock()
	s.protectionEnabled = enabled
	s.mu.Unlock()
	log.Printf("Protection mode set to: %v", enabled)
}

// SetAdblockEnabled enables or disables the adblock engine
func (s *Server) SetAdblockEnabled(enabled bool) {
	s.mu.Lock()
	s.adblockEnabled = enabled
	s.mu.Unlock()
	log.Printf("Adblock engine enabled: %v", enabled)
}

func (s *Server) ReloadAdblockEngine(rules string) {
	newEngine := adblock.NewEngine(rules)
	log.Printf("Adblock engine parsed with %d bytes of rules", len(rules))

	s.mu.Lock()
	oldEngine := s.adblockEngine
	s.adblockEngine = newEngine
	s.mu.Unlock()

	if oldEngine != nil {
		oldEngine.Close()
	}
	log.Printf("Adblock engine swapped successfully")
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
	if s.adblockEngine != nil {
		s.adblockEngine.Close()
	}
}

// GetPort returns the current port
func (s *Server) GetPort() int {
	return s.port
}

// Restart restarts the proxy with a new port
func (s *Server) Restart(port int) error {
	s.Stop()
	time.Sleep(100 * time.Millisecond) // Give it a moment
	s.port = port
	return s.Start()
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
	if domain == core.ProtocolLocalhost || req.DestAddr.IP.IsLoopback() {
		return ctx, true
	}

	// Resolve Process Name
	procName := "unknown"
	procID := int32(0)
	// req.RemoteAddr is *socks5.AddrSpec in this library version
	if r.server.systemTracker != nil && req.RemoteAddr != nil {
		procName, procID = r.server.systemTracker.GetProcessFromPort(req.RemoteAddr.Port)
	}

	// Check Adblock Engine
	sEnabled := false
	var engine *adblock.Engine

	r.server.mu.RLock()
	sEnabled = r.server.adblockEnabled
	engine = r.server.adblockEngine
	r.server.mu.RUnlock()

	if sEnabled && engine != nil {
		testURL := "http://" + domain
		log.Printf("[DEBUG] Checking adblock for: %s", testURL)
		if engine.Check(testURL, "http://"+domain, "other") {
			r.store.IncrementAdblockHit(domain)
			r.logBlock(req, domain, string(core.RuleSourceAdsblock), &core.Process{
				PID:  procID,
				Name: procName,
			})
			log.Printf("Blocked by adblock engine: %s", domain)
			return ctx, false
		} else {
			log.Printf("[DEBUG] Not blocked by adblock: %s", domain)
		}
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
			r.store.IncrementRuleHit(rule.ID, domain)

			if rule.Type == core.RuleAllow {
				r.logAllow(req, domain, &core.Process{
					PID:  procID,
					Name: procName,
				}, utils.GenerateIDString())
				return ctx, true
			}

			if rule.Type == core.RuleBlock {
				r.store.IncrementAdblockHit(domain)
				r.logBlock(req, domain, string(core.RuleSourceAdsblock), &core.Process{
					PID:  procID,
					Name: procName,
				})
				return ctx, false
			}
		}
	}

	// Check Blocklist
	if r.blocklist.IsBlocked(domain) {
		r.store.IncrementAdblockHit(domain)
		r.logBlock(req, domain, string(core.RuleSourceBlocklist), &core.Process{
			PID:  procID,
			Name: procName,
		})
		return ctx, false
	}

	// Log the connection attempt
	logID := utils.GenerateIDString()

	r.logAllow(req, domain, &core.Process{
		PID:  procID,
		Name: procName,
	}, logID)

	// Inject logID into context for Dial to pick up
	return context.WithValue(ctx, logIDKey, logID), true
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

func (r *LoggingRuleSet) logBlock(req *socks5.Request, domain string, reason string, process *core.Process) {
	entry := core.LogEntry{
		ID:          utils.GenerateIDString(),
		Timestamp:   time.Now(),
		Type:        core.LogSourceProxy,
		DstIP:       req.DestAddr.IP.String(),
		DstPort:     req.DestAddr.Port,
		SrcIP:       req.RemoteAddr.IP.String(),
		Domain:      domain,
		Protocol:    core.ProtocolTCP,
		Status:      core.LogStatusBlocked,
		BytesSent:   0,
		BytesRecv:   0,
		ProcessName: process.Name,
		ProcessID:   process.PID,
		Reason:      &reason,
	}
	r.store.AddLog(entry)
}

func (r *LoggingRuleSet) logAllow(req *socks5.Request, domain string, process *core.Process, id string) {
	entry := core.LogEntry{
		ID:          id,
		Timestamp:   time.Now(),
		Type:        core.LogSourceProxy,
		DstIP:       req.DestAddr.IP.String(),
		DstPort:     req.DestAddr.Port,
		SrcIP:       req.RemoteAddr.IP.String(),
		Domain:      domain,
		Protocol:    core.ProtocolTCP,
		Status:      core.LogStatusAllowed,
		BytesSent:   0,
		BytesRecv:   0,
		ProcessName: process.Name,
		ProcessID:   process.PID,
	}

	r.store.AddLog(entry)
}
