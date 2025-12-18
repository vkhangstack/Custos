package adblock

import (
	"testing"
)

func TestAdblockEngine(t *testing.T) {
	rules := `||ads.example.com^
||doubleclick.net^$domain=example.com
`
	engine := NewEngine(rules)
	if engine == nil {
		t.Fatal("Failed to create adblock engine")
	}
	defer engine.Close()

	tests := []struct {
		url          string
		sourceURL    string
		resourceType string
		wantBlocked  bool
	}{
		{"http://ads.example.com/banner.gif", "http://example.com", "image", true},
		{"http://example.com/index.html", "http://example.com", "document", false},
		{"http://doubleclick.net/ad.js", "http://example.com", "script", true},
		{"http://doubleclick.net/ad.js", "http://other.com", "script", false},
		{"http://pagead2.googlesyndication.com/", "http://example.com", "script", false},
	}

	for _, tt := range tests {
		got := engine.Check(tt.url, tt.sourceURL, tt.resourceType)
		if got != tt.wantBlocked {
			t.Errorf("Check(%q, %q, %q) = %v; want %v", tt.url, tt.sourceURL, tt.resourceType, got, tt.wantBlocked)
		}
	}
}
