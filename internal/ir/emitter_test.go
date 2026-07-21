package ir

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestRuntimeDescriptorSerialization verifies a frontend-authored RuntimeDescriptor
// is serialized into the emitted WIR exactly as the wolstnc backend expects: a
// top-level "runtime" object using camelCase keys (matching serde rename_all =
// "camelCase"). This is what makes `emit-ir --runtime-descriptor` work (N.3 / miss-wts-002).
func TestRuntimeDescriptorSerialization(t *testing.T) {
	arc := false
	desc := &RuntimeDescriptor{
		AllocSymbol:       "kernel_alloc",
		RetainSymbol:      "kernel_retain",
		ReleaseSymbol:     "kernel_release",
		ArcEnabled:        &arc,
		PrintNumberSymbol: "kernel_print_number",
		PrintBoolSymbol:   "kernel_print_bool",
		MethodAliases:     map[string]string{"console.log": "wolstn_console_log"},
	}
	prog := &Program{Runtime: desc}
	data, err := json.Marshal(prog)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}

	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	rt, ok := m["runtime"]
	if !ok {
		t.Fatalf("expected top-level 'runtime' field in WIR; got %s", string(data))
	}
	rm := rt.(map[string]any)

	if rm["arcEnabled"] != false {
		t.Fatalf("expected arcEnabled=false, got %v", rm["arcEnabled"])
	}
	if rm["allocSymbol"] != "kernel_alloc" {
		t.Fatalf("expected allocSymbol=kernel_alloc, got %v", rm["allocSymbol"])
	}
	if rm["retainSymbol"] != "kernel_retain" {
		t.Fatalf("expected retainSymbol=kernel_retain, got %v", rm["retainSymbol"])
	}
	if rm["releaseSymbol"] != "kernel_release" {
		t.Fatalf("expected releaseSymbol=kernel_release, got %v", rm["releaseSymbol"])
	}
	if rm["printNumberSymbol"] != "kernel_print_number" {
		t.Fatalf("expected printNumberSymbol=kernel_print_number, got %v", rm["printNumberSymbol"])
	}
	if rm["printBoolSymbol"] != "kernel_print_bool" {
		t.Fatalf("expected printBoolSymbol=kernel_print_bool, got %v", rm["printBoolSymbol"])
	}
	aliases, ok := rm["methodAliases"].(map[string]any)
	if !ok {
		t.Fatalf("expected methodAliases object, got %v", rm["methodAliases"])
	}
	if aliases["console.log"] != "wolstn_console_log" {
		t.Fatalf("expected methodAliases[console.log]=wolstn_console_log, got %v", aliases["console.log"])
	}
}

// TestRuntimeDescriptorOmittedWhenNil verifies the "runtime" field is omitted
// entirely (omitempty) when no descriptor is authored, so the backend applies its
// own defaults and we never emit a null/empty object.
func TestRuntimeDescriptorOmittedWhenNil(t *testing.T) {
	prog := &Program{Version: 1}
	data, err := json.Marshal(prog)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if strings.Contains(string(data), "runtime") {
		t.Fatalf("expected 'runtime' to be omitted when nil; got %s", string(data))
	}
}

// TestRuntimeDescriptorArcEnabledDefaultsToBackend verifies that when the author does
// NOT set arcEnabled, the field is omitted (pointer nil) and the backend default
// (true) applies — i.e. we don't force-disable ARC silently.
func TestRuntimeDescriptorArcEnabledUnsetOmitted(t *testing.T) {
	desc := &RuntimeDescriptor{AllocSymbol: "kernel_alloc"}
	prog := &Program{Runtime: desc}
	data, _ := json.Marshal(prog)
	var m map[string]any
	_ = json.Unmarshal(data, &m)
	rm := m["runtime"].(map[string]any)
	if _, ok := rm["arcEnabled"]; ok {
		t.Fatalf("arcEnabled must be omitted when unset; backend default (true) should apply")
	}
}
