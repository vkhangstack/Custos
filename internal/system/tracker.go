package system

import (
	"fmt"
	"log"
	"sync"

	"github.com/shirou/gopsutil/v4/net"
	"github.com/shirou/gopsutil/v4/process"
)

// ConnectionInfo holds details about an active connection
type ConnectionInfo struct {
	PID         int32  `json:"pid"`
	ProcessName string `json:"process_name"`
	LocalAddr   string `json:"local_addr"`
	RemoteAddr  string `json:"remote_addr"`
	Status      string `json:"status"`
	Protocol    string `json:"protocol"` // "tcp" or "udp"
}

// Tracker manages system process and network monitoring
type Tracker struct {
	mu           sync.RWMutex
	processCache map[int32]string
}

// NewTracker creates a new system tracker
func NewTracker() *Tracker {
	return &Tracker{
		processCache: make(map[int32]string),
	}
}

// GetProcessName returns the name of a process by PID, caching the result
func (t *Tracker) GetProcessName(pid int32) string {
	if pid == 0 {
		return "kernel"
	}

	t.mu.RLock()
	// Cache invalidation logic could be added here (e.g. check if PID is reused)
	// For now, let's just rely on simple map
	if name, ok := t.processCache[pid]; ok {
		t.mu.RUnlock()
		return name
	}
	t.mu.RUnlock()

	proc, err := process.NewProcess(pid)
	if err != nil {
		return "unknown"
	}
	name, err := proc.Name()
	if err != nil {
		return "unknown"
	}

	t.mu.Lock()
	t.processCache[pid] = name
	t.mu.Unlock()

	return name
}

// GetProcessFromPort attempts to identify the process owning a local port
// This is used to identify the source of a connection to the proxy
func (t *Tracker) GetProcessFromPort(port int) (string, int32) {
	// 1. Get all TCP connections
	conns, err := net.Connections("tcp")
	if err != nil {
		log.Printf("[Tracker] Error getting connections: %v", err)
		return "unknown", 0
	}

	log.Printf("[Tracker] Searching for process on port %d. Found %d active TCP connections.", port, len(conns))

	// 2. Find the connection where Laddr.Port matches our target port
	for _, conn := range conns {
		if int(conn.Laddr.Port) == port {
			// Found it!
			pid := conn.Pid
			name := t.GetProcessName(pid)
			log.Printf("[Tracker] Found match! Port %d -> PID %d (%s)", port, pid, name)
			return name, pid
		}
	}

	log.Printf("[Tracker] No process found for port %d", port)
	return "unknown", 0
}

// GetActiveConnections returns a list of current network connections
func (t *Tracker) GetActiveConnections() ([]ConnectionInfo, error) {
	conns, err := net.Connections("all")
	if err != nil {
		log.Printf("Error getting active connections: %v", err)
		return nil, err
	}

	var results []ConnectionInfo
	for _, conn := range conns {
		name := t.GetProcessName(conn.Pid)

		info := ConnectionInfo{
			PID:         conn.Pid,
			ProcessName: name,
			LocalAddr:   fmt.Sprintf("%s:%d", conn.Laddr.IP, conn.Laddr.Port),
			RemoteAddr:  fmt.Sprintf("%s:%d", conn.Raddr.IP, conn.Raddr.Port),
			Status:      conn.Status,
		}

		if conn.Type == 1 { // TCP
			info.Protocol = "tcp"
		} else if conn.Type == 2 { // UDP
			info.Protocol = "udp"
		} else {
			info.Protocol = "unknown"
		}

		results = append(results, info)
	}

	return results, nil
}
