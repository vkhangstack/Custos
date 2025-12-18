package utils

import (
	"fmt"
	"sync"
	"time"
)

// Snowflake holds the basic settings for a snowflake ID generator
// format: 41 bits (ms) | 10 bits (node) | 12 bits (seq)
type Snowflake struct {
	mu       sync.Mutex
	epoch    int64
	nodeID   int64
	lastTime int64
	sequence int64
}

var (
	defaultSnowflake *Snowflake
	once             sync.Once
)

// NewSnowflake creates a new generator
func NewSnowflake(nodeID int64) *Snowflake {
	return &Snowflake{
		epoch:  1734512400000,  // 2024-12-18 00:00:00 UTC in ms
		nodeID: nodeID & 0x3FF, // 10 bits
	}
}

// GetID returns a new snowflake ID
func (s *Snowflake) Generate() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UnixMilli()

	if now == s.lastTime {
		s.sequence = (s.sequence + 1) & 0xFFF // 12 bits
		if s.sequence == 0 {
			// Sequence exhausted, wait for next ms
			for now <= s.lastTime {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		s.sequence = 0
	}

	s.lastTime = now

	id := ((now - s.epoch) << 22) | (s.nodeID << 12) | s.sequence
	return id
}

// GetGlobalID is a helper for a default global instance
func GenerateID() int64 {
	once.Do(func() {
		defaultSnowflake = NewSnowflake(1)
	})
	return defaultSnowflake.Generate()
}

// GenerateIDString returns the ID as a string
func GenerateIDString() string {
	return fmt.Sprintf("%d", GenerateID())
}
