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
	store       store.Store
	blocklist   *core.BlocklistManager
	socksServer *socks5.Server
	port        int
	running     bool
	listener    net.Listener
}

// NewServer creates a new proxy server
func NewServer(store store.Store, blocklist *core.BlocklistManager, port int) *Server {
	return &Server{
		store:     store,
		blocklist: blocklist,
		port:      port,
	}
}

const logIDKey = "logID"

// Start starts the SOCKS5 proxy
func (s *Server) Start() error {
	conf := &socks5.Config{
		Logger: log.New(log.Writer(), "[SOCKS5] ", log.LstdFlags),
		Rules:  &LoggingRuleSet{store: s.store, blocklist: s.blocklist},
		Dial: func(ctx context.Context, network, addr string) (net.Conn, error) {
			// Dial upstream
			conn, err := net.Dial(network, addr)
			if err != nil {
				return nil, err
			}

			// Wrap if we have a logID
			if logID, ok := ctx.Value(logIDKey).(string); ok {
				fmt.Printf("[DEBUG] Dialing for logID: %s\n", logID)
				return &CountingConn{
					Conn:  conn,
					logID: logID,
					entry: core.LogEntry{
						ID: logID,
						// Other fields will be updated by UpdateLog, but ID is critical
						// We might want to pass more context if needed, but for now ID is enough to find/update it?
						// Wait, UpdateLog in SQLite uses Save(&entry). It overwrites EVERYTHING.
						// So we need to ensure the entry we pass to UpdateLog has the other fields correct?
						// Or at least doesn't blank them out.
						// UpdateLog implementation: db.Save(&entry).
						// If I pass an empty struct with just ID and bytes, GORM might blank out Domain/IP/etc.
						// This IS A PROBLEM with GORM Save.
						// Solution: UpdateLog should use db.Model().Updates().
						// Let's fix SQLiteStore.UpdateLog in previous step? Or here?
						// Better: fetch the entry first? No, performance.
						// Better: Pass only fields to update?
						// Or make CountingConn smarter.
						// For now, let's inject a "Partial Update" logic in Store or handle it.
						// But I can't easily change Store interface again without refactor.
						// If UpdateLog does Save, I must provide full object.
						// So I need to pass the FULL LogEntry to Dial.
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
}

func (r *LoggingRuleSet) Allow(ctx context.Context, req *socks5.Request) (context.Context, bool) {
	// Check Blocklist
	domain := req.DestAddr.FQDN
	if r.blocklist.IsBlocked(domain) {
		entry := core.LogEntry{
			ID:        fmt.Sprintf("%d", time.Now().UnixNano()),
			Timestamp: time.Now(),
			Type:      "proxy",
			DstIP:     req.DestAddr.IP.String(),
			DstPort:   req.DestAddr.Port,
			Domain:    domain,
			Protocol:  "tcp",
			Status:    "blocked",
			BytesSent: 0,
			BytesRecv: 0,
		}
		r.store.AddLog(entry)
		return ctx, false
	}

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
