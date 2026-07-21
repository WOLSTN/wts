// WOLSTN kernel lib (lib.kernel.d.ts)
//
// A minimal, host-independent runtime type surface for `no_std` / kernel targets.
//
// This lib intentionally declares NO host types (no DOM, no Node, no Browser). It only
// provides the runtime symbols a kernel author needs, so that `console.log` type-checks
// without pulling in the entire JS host surface.
//
// It is selected via tsconfig `lib: ["kernel"]`, or by `noLib: true` plus a user-authored
// `declare` of the same surface. This decouples the *type* of the runtime from any specific
// host (noStdRoadMap N.3 / de-privileged design, AGENTS.md §8).
//
// The concrete runtime symbols (wolstn_console_log, wolstn_alloc, ...) are bound later,
// either by the backend defaults or by a frontend-authored RuntimeDescriptor
// (emit-ir --runtime-descriptor).

declare namespace console {
    /** Write a line to the kernel console. Accepts any number of arguments; the backend
     * chooses how to render each value (number -> printNumber, bool -> printBool, etc.). */
    function log(...args: any[]): void;
}

// Allow kernel code to opt into the standard `no_std` entry shape without depending on
// the host `lib.es5`/DOM libs. `globalThis` is intentionally left to the backend's
// intrinsic handling; we only declare the kernel-facing `console` here.
