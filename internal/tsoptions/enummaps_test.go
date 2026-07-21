package tsoptions

import "testing"

// TestLibMapKernel verifies the WOLSTN-specific "kernel" lib resolves to
// lib.kernel.d.ts (noStdRoadMap N.3). Resolution must be case-insensitive to match
// the host lib behavior.
func TestLibMapKernel(t *testing.T) {
	name, ok := GetLibFileName("kernel")
	if !ok {
		t.Fatalf("expected 'kernel' to resolve to a lib file, got ok=false")
	}
	if name != "lib.kernel.d.ts" {
		t.Fatalf("expected 'lib.kernel.d.ts', got %q", name)
	}

	name2, ok2 := GetLibFileName("Kernel")
	if !ok2 || name2 != "lib.kernel.d.ts" {
		t.Fatalf("expected case-insensitive 'kernel' resolution, got %q ok=%v", name2, ok2)
	}

	// The kernel lib must also be reachable as a file name directly.
	if _, ok := GetLibFileName("lib.kernel.d.ts"); !ok {
		t.Fatalf("expected 'lib.kernel.d.ts' to be a valid lib file name")
	}
}
