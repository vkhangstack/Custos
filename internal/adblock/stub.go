//go:build !cgo

package adblock

import "log"

type Engine struct {
	ptr unsafePointer // Not used but keeps struct size non-zero if needed
}

type unsafePointer uintptr

func NewEngine(rules string) *Engine {
	log.Println("WARNING: Adblock engine is disabled because CGO is not enabled.")
	return &Engine{}
}

func (e *Engine) Check(url, sourceURL, resourceType string) bool {
	return false
}

func (e *Engine) Close() {
}
