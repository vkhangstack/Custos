//go:build cgo

package adblock

/*
#cgo LDFLAGS: -L${SRCDIR}/../../lib/adblock/target/release -ladblock -ldl -lpthread -lm
#include <stdlib.h>
#include <stdbool.h>

typedef void* adblock_engine_t;

adblock_engine_t adblock_engine_create(const char* rules);
bool adblock_engine_check(adblock_engine_t engine, const char* url, const char* source_url, const char* resource_type);
void adblock_engine_destroy(adblock_engine_t engine);
*/
import "C"
import (
	"unsafe"
)

type Engine struct {
	ptr C.adblock_engine_t
}

func NewEngine(rules string) *Engine {
	cRules := C.CString(rules)
	defer C.free(unsafe.Pointer(cRules))

	ptr := C.adblock_engine_create(cRules)
	if ptr == nil {
		return nil
	}

	return &Engine{ptr: ptr}
}

func (e *Engine) Check(url, sourceURL, resourceType string) bool {
	if e.ptr == nil {
		return false
	}

	cUrl := C.CString(url)
	defer C.free(unsafe.Pointer(cUrl))

	cSourceURL := C.CString(sourceURL)
	defer C.free(unsafe.Pointer(cSourceURL))

	cResourceType := C.CString(resourceType)
	defer C.free(unsafe.Pointer(cResourceType))

	return bool(C.adblock_engine_check(e.ptr, cUrl, cSourceURL, cResourceType))
}

func (e *Engine) Close() {
	if e.ptr != nil {
		C.adblock_engine_destroy(e.ptr)
		e.ptr = nil
	}
}
