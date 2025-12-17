package proxy

import (
	"Custos/internal/core"
	"Custos/internal/store"
	"fmt"
	"net"
	"sync/atomic"
	"time" // Added for time.Now().UnixMilli()
)

// CountingConn wraps a net.Conn to track bytes sent and received
type CountingConn struct {
	net.Conn
	logID      string
	entry      core.LogEntry
	store      store.Store
	bytesSent  int64
	bytesRecv  int64
	lastUpdate int64 // UnixMilli timestamp of last db update
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

// tryUpdate updates the DB if enough time has passed
func (c *CountingConn) tryUpdate() {
	now := time.Now().UnixMilli()
	last := atomic.LoadInt64(&c.lastUpdate)
	fmt.Printf("[DEBUG] TryUpdate conn %s. Sent: %d, Recv: %d\n", c.entry.ID, c.bytesSent, c.bytesRecv)

	if now-last > updateIntervalMilli {
		// Try to swap to prevent concurrent updates (thundering herd is unlikely on single conn but good practice)
		if atomic.CompareAndSwapInt64(&c.lastUpdate, last, now) {
			go c.doUpdate()
		}
	}
}

// doUpdate performs the actual store update
func (c *CountingConn) doUpdate() {
	sent := atomic.LoadInt64(&c.bytesSent)
	recv := atomic.LoadInt64(&c.bytesRecv)

	// Create a copy for update
	updateEntry := c.entry
	updateEntry.BytesSent = sent
	updateEntry.BytesRecv = recv
	fmt.Printf("[DEBUG] Update conn %s. Sent: %d, Recv: %d\n", c.entry.ID, sent, recv)

	// We only update the log entry table, NOT the global traffic stats (AddedTraffic)
	// Global stats should arguably only update on Close or separate periodic sync to avoid write hotspots.
	// For "Live Logs" visuals, updating the log entry is enough.
	c.store.UpdateLog(updateEntry)
}

// Close wraps Close to log final stats
func (c *CountingConn) Close() error {
	// Update final Log
	sent := atomic.LoadInt64(&c.bytesSent)
	recv := atomic.LoadInt64(&c.bytesRecv)

	// Update Log Entry
	c.entry.BytesSent = sent
	c.entry.BytesRecv = recv

	// Debug Log
	fmt.Printf("[DEBUG] Close conn %s. Sent: %d, Recv: %d\n", c.entry.ID, sent, recv)

	// The core.LogEntry has Latency field. Let's calculate duration.
	// But we don't track start time here. We could.
	// Let's assume latency was initial connection time.

	c.store.UpdateLog(c.entry)

	// Add to global stats ONLY on close to ensure we don't double count or hammer DB
	c.store.AddTraffic(sent, recv)

	return c.Conn.Close()
}
