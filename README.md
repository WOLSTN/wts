# WTS - TypeScript Frontend for Wolstn

**WTS** (Wolstn TypeScript) is a TypeScript frontend for the [Wolstn](https://github.com/WOLSTN/wolstn) native compiler. It parses TypeScript code, performs type checking, and emits a typed IR (Intermediate Representation) that Wolstn's backend can compile to native executables.

## Overview

WTS is forked from [Microsoft's TypeScript-go](https://github.com/microsoft/TypeScript-go) (tsgo), but with significant differences:

| | tsgo | WTS |
|---|---|---|
| **Purpose** | TypeScript → JavaScript | TypeScript → Wolstn IR |
| **Output** | JavaScript code | Typed IR (JSON/binary) |
| **Type info** | Used for checking only | **Preserved for codegen** |
| **Runtime** | V8, Node.js, browsers | Native (ELF/EXE) via Wolstn |

## Why WTS?

Unlike Deno and Bun, which transpile TypeScript to JavaScript for execution in a JS runtime, **Wolstn compiles TypeScript directly to native machine code**. This requires preserving full type information for:

- **Memory layout decisions**: `number` → f64, `string` → GC pointer
- **Calling conventions**: Generic instantiation, overload resolution
- **Optimizations**: Inlining, escape analysis, type-driven optimizations

WTS provides this type information through its IR output.

## Installation

```bash
go build -o wts ./cmd/wts
```

## Usage

### Type Check

```bash
wts check main.ts
wts check --project tsconfig.json
```

### Emit IR

```bash
wts emit-ir main.ts -o main.wir
wts emit-ir --project tsconfig.json -o program.wir
```

The output `.wir` file is a JSON-formatted IR containing:

- **Files**: Source file paths and AST nodes
- **Types**: All resolved types with their properties
- **Symbols**: All symbols (functions, classes, variables, etc.)
- **Functions**: Function signatures and parameters
- **Globals**: Top-level declarations

## IR Format

The IR is versioned and designed for forward compatibility:

```json
{
  "version": 1,
  "files": [...],
  "types": [...],
  "symbols": [...],
  "functions": [...],
  "globals": [...]
}
```

See [internal/ir/types.go](internal/ir/types.go) for the full schema.

## Architecture

```
Source Code (.ts)
      ↓
   Parser
      ↓
    AST
      ↓
   Binder
      ↓
  Checker
      ↓
Typed AST + Types
      ↓
 IR Emitter
      ↓
Wolstn IR (.wir)
```

## Modules

WTS retains only the frontend components from tsgo:

| Module | Purpose |
|--------|---------|
| `ast` | AST definitions |
| `parser` | Source code parsing |
| `scanner` | Lexical analysis |
| `binder` | Symbol binding |
| `checker` | Type checking |
| `ir` | IR emission (new) |

Removed components (not needed for native compilation):

- JS emitter, printer, transformers
- LSP server, language service
- API server
- Format, fourslash tests

## Development

WTS tracks upstream tsgo for parser/checker updates. The frontend components (parser, binder, checker) are stable and rarely change their API.

## License

Apache License 2.0 - see [LICENSE](LICENSE) for details.

## Acknowledgments

- Microsoft TypeScript Team for [TypeScript-go](https://github.com/microsoft/TypeScript-go)
- The Wolstn project
