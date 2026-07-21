# wts — Known Omissions / Technical Debt

This file transparently documents simplified, unimplemented, or deferred behaviors
in the `wts` frontend, per the project's "No Silent Failures / Transparent Omissions"
guideline. Each entry lists the affected file(s), when it should be properly
implemented, and what the correct implementation looks like.

## Resolved

### [miss-wts-001] `-p` project (tsconfig) mode was unusable with the wolstnc backend — RESOLVED

- **Files**: `cmd/wts/main.go` (`createProgramFromArgs`), `internal/ir/emitter.go` (`NewEmitter`).
- **Was**: Passing `-p tsconfig.json` compiled the config file itself as a TS source;
  `-p .` ignored `noLib` and force-included the default TS lib; and project-mode WIR
  emitted `null` for top-level sequence fields (e.g. `classes`, `interfaces`), which
  wolstnc's WIR parser rejects (`invalid type: null, expected a sequence`).
- **Fix**: `createProgramFromArgs` now genuinely parses the tsconfig via
  `tsoptions.ParseConfigFileTextToJson` + `ParseJsonConfigFileContent` (resolving
  `files`/`include`/`exclude` and `compilerOptions` such as `noLib`/`lib`). `NewEmitter`
  pre-seeds every top-level `Program` sequence field (`Files`, `Types`, `Symbols`,
  `Signatures`, `Globals`, `Classes`, `Interfaces`) to an empty non-nil slice so the
  emitted WIR is always a valid JSON array. With `noLib: true` in the tsconfig, the
  default library is no longer injected and the WIR stays minimal.
- **Verified**: `test/6-stdlib-link` now builds and runs end-to-end via
  `wts emit-ir -p tsconfig.json` (Passed: 7, Failed: 0).
