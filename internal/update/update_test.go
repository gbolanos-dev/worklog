package update

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var tests = []struct {
	name    string
	current string
	latest  string
	isNewer bool
}{
	{"patch bump", "v1.0.0", "v1.0.1", true},
	{"major bump", "v1.0.0", "v2.0.0", true},
	{"older version", "v1.0.0", "v0.9.9", false},
	{"same version", "v1.0.0", "v1.0.0", false},
	{"missing v prefix", "1.0.0", "v1.0.1", true},
}

func TestNewer(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newer(tt.current, tt.latest)
			if got != tt.isNewer {
				t.Errorf("newer(%s, %s) = %t, want %t", tt.current, tt.latest, got, tt.isNewer)
			}
		})
	}
}

func TestCacheReadWrite(t *testing.T) {
	dir := t.TempDir()

	c := cache{
		CheckedAt: time.Now().UTC(),
		Latest:    "v1.1.0",
	}

	data, _ := json.Marshal(c)
	_ = os.WriteFile(filepath.Join(dir, "update-check.json"), data, 0666)

	// Check should find the cache and compare versions
	result := Check("v1.0.0", dir)
	if result == nil {
		t.Fatal("Check() returned nil")
	}
	if !result.Available {
		t.Error("expected Available to be true")
	}
	if result.Latest != "v1.1.0" {
		t.Errorf("expected Latest to be v1.0.0, got %s", result.Latest)
	}
}
