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

## Active / Deferred

### [miss-wts-002] TS 作者侧缺少声明自定义 RuntimeDescriptor 的语法/编译开关 —— DEFERRED (N.3)

- **类型**：`miss(Others)`
- **状态**：待实现（留待路线图 N.3）
- **描述**：wolstnc 后端（N.1，见 `doc/noStdRoadMap.md`）已支持数据驱动的 `RuntimeDescriptor`：
  只要 WIR 顶层带有 `runtime.methodAliases` / `printNumberSymbol` / `printBoolSymbol`，
  `console.log` 即可被重定向到任意用户符号（如 `kernel_puts`），编译器不再对 `console` 做硬编码特判。
  但**目前 wts 前端没有任何语法或编译开关**让 TS 作者声明这套重定向。
  当前只能通过 E2E 测试里的 `inject_runtime.py` 在 WIR JSON 层面手动注入。
- **涉及文件**：`internal/ir/emitter.go`（`NewEmitter` / `Program` 发射）、`cmd/wts/main.go`（`createProgramFromArgs`）。
- **实现条件**：路线图 N.3（内核 lib 类型层解耦）阶段，与去特权化旋钮一并落地。
- **正确实现**：
  1. 提供用户友好的声明方式，例如：
     - 编译选项：`wts emit-ir --runtime-descriptor path/to/runtime.json`（读取与 WIR `runtime` 同构的 JSON）；或
     - TS 级编译指示：`#[runtime({ "methodAliases": { "console.log": "kernel_puts" }, "printNumberSymbol": "kernel_print_number" })]` 顶层装饰器 / top-level 标注；
  2. 前端把这些声明原样写入发射出的 WIR 顶层 `runtime` 字段，不做任何改写或特判；
  3. 校验：若用户声明的符号名与目标 no_std 运行时（`.c`/`.a`）实际导出不一致，链接期由 clang 暴露 `undefined symbol`，而非前端静默通过——保持「无静默失败」。
- **不做的妥协**：不要在前端为 `console.log` 写死任何"内置识别"，全部经由 WIR `runtime` 数据驱动，遵循去特权化原则（AGENTS.md 第 8 条）。
