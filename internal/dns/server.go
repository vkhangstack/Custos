package dns

import (
	"fmt"
	"log"
	"time"

	"github.com/vkhangstack/Custos/internal/core"
	"github.com/vkhangstack/Custos/internal/store"
	"github.com/vkhangstack/Custos/internal/utils"

	"github.com/miekg/dns"
)

// Server manages the DNS listener
type Server struct {
	store     store.Store
	blocklist *core.BlocklistManager
	server    *dns.Server
	port      int
	upstream  string
}

// NewServer creates a new DNS server
func NewServer(store store.Store, blocklist *core.BlocklistManager, port int) *Server {
	return &Server{
		store:     store,
		blocklist: blocklist,
		port:      port,
		upstream:  "8.8.8.8:53",
	}
}

// Start starts the DNS server
func (s *Server) Start() {
	addr := fmt.Sprintf(":%d", s.port)
	s.server = &dns.Server{Addr: addr, Net: "udp"}

	dns.HandleFunc(".", s.handleRequest)

	log.Printf("DNS Server started on %s", addr)
	go func() {
		if err := s.server.ListenAndServe(); err != nil {
			log.Printf("DNS server error: %v", err)
		}
	}()
}

// Stop stops the DNS server
func (s *Server) Stop() {
	if s.server != nil {
		s.server.Shutdown()
	}
}

func (s *Server) handleRequest(w dns.ResponseWriter, r *dns.Msg) {
	if len(r.Question) == 0 {
		return
	}
	q := r.Question[0]

	// Check Blocklist
	if s.blocklist.IsBlocked(q.Name) {
		// Log blocked
		entry := core.LogEntry{
			ID:        utils.GenerateIDString(),
			Timestamp: time.Now(),
			Type:      "dns",
			Domain:    q.Name,
			SrcIP:     w.RemoteAddr().String(),
			Protocol:  "udp",
			Status:    "blocked",
			BytesSent: 0,
			BytesRecv: 0,
		}
		s.store.AddLog(entry)

		// Return 0.0.0.0
		m := new(dns.Msg)
		m.SetReply(r)
		rr, _ := dns.NewRR(fmt.Sprintf("%s A 0.0.0.0", q.Name))
		m.Answer = append(m.Answer, rr)
		w.WriteMsg(m)
		return
	}

	// Forward to Upstream
	c := new(dns.Client)
	resp, _, err := c.Exchange(r, s.upstream)
	if err != nil {
		log.Printf("DNS upstream error: %v", err)
		return
	}

	// Log allowed
	entry := core.LogEntry{
		ID:        utils.GenerateIDString(),
		Timestamp: time.Now(),
		Type:      "dns",
		Domain:    q.Name,
		SrcIP:     w.RemoteAddr().String(),
		Protocol:  "udp",
		Status:    "allowed",
	}
	s.store.AddLog(entry)

	resp.Id = r.Id
	w.WriteMsg(resp)
}
