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
		s.logBlock(q.Name, w.RemoteAddr().String(), string(core.RuleSourceBlocklist))
		s.replyBlocked(w, r, q.Name)
		return
	}

	// Check Custom Rules
	rules := s.store.GetRules()
	for _, rule := range rules {
		if !rule.Enabled {
			continue
		}
		if matched, _ := matchDomain(rule.Pattern, q.Name); matched {
			s.store.IncrementRuleHit(rule.ID, q.Name)
			if rule.Type == core.RuleBlock {
				s.logBlock(q.Name, w.RemoteAddr().String(), string(core.RuleSourceCustom))
				s.replyBlocked(w, r, q.Name)
				return
			}
			// If ALLOW, skip further block checks
			break
		}
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
		Type:      core.LogSourceDNS,
		Domain:    q.Name,
		SrcIP:     w.RemoteAddr().String(),
		Protocol:  core.ProtocolUDP,
		Status:    core.LogStatusAllowed,
	}
	s.store.AddLog(entry)

	resp.Id = r.Id
	w.WriteMsg(resp)
}

func (s *Server) logBlock(domain, srcIP, source string) {
	entry := core.LogEntry{
		ID:        utils.GenerateIDString(),
		Timestamp: time.Now(),
		Type:      core.LogSourceDNS,
		Domain:    domain,
		SrcIP:     srcIP,
		Protocol:  core.ProtocolUDP,
		Status:    core.LogStatusBlocked,
		BytesSent: 0,
		BytesRecv: 0,
	}
	s.store.AddLog(entry)
}

func (s *Server) replyBlocked(w dns.ResponseWriter, r *dns.Msg, domain string) {
	m := new(dns.Msg)
	m.SetReply(r)
	rr, _ := dns.NewRR(fmt.Sprintf("%s A 0.0.0.0", domain))
	m.Answer = append(m.Answer, rr)
	w.WriteMsg(m)
}

// matchDomain checks if domain matches pattern
func matchDomain(pattern, domain string) (bool, error) {
	domain = fmt.Sprintf("%s", domain)
	// Remove trailing dot if present
	if len(domain) > 0 && domain[len(domain)-1] == '.' {
		domain = domain[:len(domain)-1]
	}
	if len(pattern) > 2 && pattern[:2] == "*." {
		suffix := pattern[2:]
		if len(domain) >= len(suffix) && domain[len(domain)-len(suffix):] == suffix {
			return true, nil
		}
	}
	return pattern == domain, nil
}
