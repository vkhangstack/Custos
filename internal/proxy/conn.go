package proxy

import (
	"Custos/internal/core"
	"Custos/internal/store"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

// CountingConn wraps a net.Conn to track bytes sent and received
type CountingConn struct {
	net.Conn
	logID      string
	entry      core.LogEntry
	store      store.Store
	bytesSent  int64
	bytesRecv  int64
	lastUpdate int64 // UnixMilli timestamp of last checking

	mu           sync.Mutex
	reportedSent int64
	reportedRecv int64
}

const updateIntervalMilli = 1000 // Update DB at most once per second

// Read wraps Read to count bytes
func (c *CountingConn) Read(b []byte) (n int, err error) {
	n, err = c.Conn.Read(b)
	if n > 0 {
		atomic.AddInt64(&c.bytesRecv, int64(n))
		c.tryUpdate()
	}
	return
}

// Write wraps Write to count bytes
func (c *CountingConn) Write(b []byte) (n int, err error) {
	n, err = c.Conn.Write(b)
	if n > 0 {
		atomic.AddInt64(&c.bytesSent, int64(n))
		c.tryUpdate()
	}
	return
}

// tryUpdate checks if a report is needed
func (c *CountingConn) tryUpdate() {
	now := time.Now().UnixMilli()
	last := atomic.LoadInt64(&c.lastUpdate)

	if now-last > updateIntervalMilli {
		// Try to swap to prevent concurrent updates
		if atomic.CompareAndSwapInt64(&c.lastUpdate, last, now) {
			go c.report()
		}
	}
}

// report performs the actual DB update calculating deltas
func (c *CountingConn) report() {
	c.mu.Lock()
	defer c.mu.Unlock()

	currSent := atomic.LoadInt64(&c.bytesSent)
	currRecv := atomic.LoadInt64(&c.bytesRecv)

	deltaSent := currSent - c.reportedSent
	deltaRecv := currRecv - c.reportedRecv

	if deltaSent == 0 && deltaRecv == 0 {
		return
	}

	// 1. Update Global Stats (Incremental)
	c.store.AddTraffic(deltaSent, deltaRecv)

	// 2. Update reported cursors
	c.reportedSent = currSent
	c.reportedRecv = currRecv

	// 3. Update Log Entry (Total)
	c.entry.BytesSent = currSent
	c.entry.BytesRecv = currRecv

	// Debug log only on significant updates or errors?
	// fmt.Printf("[DEBUG] Report conn %s. Delta: %d/%d Total: %d/%d\n", c.entry.ID, deltaSent, deltaRecv, currSent, currRecv)

	c.store.UpdateLog(c.entry)
}

// Close wraps Close to log final stats
func (c *CountingConn) Close() error {
	// Force a final report
	c.report()
	fmt.Printf("[DEBUG] Closed conn %s. Final sent/recv: %d/%d\n", c.entry.ID, c.reportedSent, c.reportedRecv)
	return c.Conn.Close()
}
