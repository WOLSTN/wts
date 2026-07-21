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

### [miss-wts-002] TS 作者侧声明自定义 RuntimeDescriptor 的编译开关 —— RESOLVED (N.3)

- **类型**：`miss(Others)`
- **状态**：已实现（路线图 N.3）
- **描述**：wolstnc 后端（N.1）支持数据驱动的 `RuntimeDescriptor`（`runtime.methodAliases` 等），
  但 wts 前端原本没有任何语法/开关让 TS 作者声明这套重定向，只能靠 E2E 的 `inject_runtime.py` 后注入。
- **实现**：
  1. `internal/ir/emitter.go` 新增 `RuntimeDescriptor` 结构（与 WIR `runtime` 同构，camelCase；
     `ArcEnabled *bool` 用指针保留后端默认值），`Program` 增加 `Runtime *RuntimeDescriptor`（JSON `runtime`），
     `EmitOptions` 增加 `Runtime *RuntimeDescriptor`；`Emit()` 在发射前 `e.irProgram.Runtime = e.options.Runtime`。
  2. `cmd/wts/main.go` 新增 `--runtime-descriptor <path>` 编译开关，读取与 WIR `runtime` 同构的 JSON
     并赋给 `EmitOptions.Runtime`；无特判、无写死 `console`/`wolstn_*`。
  3. 校验：符号名与目标 no_std 运行时（`.c`/`.a`）导出不一致时，由 clang 链接期暴露 `undefined symbol`，
     前端不静默通过——保持「无静默失败」。
- **涉及文件**：`internal/ir/emitter.go`、`cmd/wts/main.go`。
- **验收**：`test/12-nostd-kernel`（`wts emit-ir -p tsconfig.json --runtime-descriptor runtime.json` + `noLib` 变体）
  Passed: 26, Failed: 0；WIR 顶层确含 `runtime` 且 `console.log -> kernel_puts`，无任何 `wolstn_*` 后门符号。
- **正确实现**，已落地为编译选项（TS 级 `#[runtime(...)]` 装饰器未做，留待用户需求驱动；当前开关已满足去特权化）。
