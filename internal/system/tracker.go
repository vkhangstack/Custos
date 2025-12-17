package system

import (
	"fmt"
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

// GetActiveConnections returns a list of current network connections
func (t *Tracker) GetActiveConnections() ([]ConnectionInfo, error) {
	conns, err := net.Connections("all")
	if err != nil {
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
