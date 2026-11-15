package gestaltung

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// sync.sh must produce exactly what the Go programs embed.
func TestSyncMatchesCSS(t *testing.T) {
	dir := t.TempDir()
	if out, err := exec.Command("sh", "scripts/sync.sh", dir).CombinedOutput(); err != nil {
		t.Fatalf("sync.sh: %v\n%s", err, out)
	}
	got, err := os.ReadFile(filepath.Join(dir, "gestaltung.css"))
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != CSS {
		t.Fatalf("sync.sh output differs from gestaltung.CSS (%d vs %d bytes)", len(got), len(CSS))
	}
}
