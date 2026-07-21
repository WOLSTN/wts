package bundled

import (
	"strings"
	"testing"
)

// TestKernelLibBundled verifies lib.kernel.d.ts ships in the embedded bundle and
// provides a minimal, host-independent `console` type surface (noStdRoadMap N.3).
func TestKernelLibBundled(t *testing.T) {
	// It must be present in the registered lib names.
	found := false
	for _, n := range LibNames {
		if n == "lib.kernel.d.ts" {
			found = true
			break
		}
	}
	if !found {
		t.Fatalf("lib.kernel.d.ts missing from bundled LibNames")
	}

	content, ok := embeddedContents["libs/lib.kernel.d.ts"]
	if !ok {
		t.Fatalf("expected lib.kernel.d.ts to be retrievable from embedded bundle")
	}

	// Minimal runtime surface: a console namespace with log().
	if !strings.Contains(content, "declare namespace console") {
		t.Fatalf("kernel lib must declare 'console' namespace; got:\n%s", content)
	}
	if !strings.Contains(content, "function log") {
		t.Fatalf("kernel lib must declare console.log")
	}

	// Host-independence: must NOT ship DOM/Node host types. A kernel/no_std target
	// relies on this lib precisely to AVOID pulling in the JS host surface.
	hostMarkers := []string{
		"declare var document",
		"declare function setTimeout",
		"interface Window",
		"declare var process",
		"lib.dom",
	}
	for _, m := range hostMarkers {
		if strings.Contains(content, m) {
			t.Fatalf("kernel lib must not ship host types; found %q", m)
		}
	}
}
