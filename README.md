# WTS - TypeScript Frontend for Wolstn

**WTS** (Wolstn TypeScript) is a TypeScript frontend for the [Wolstn](https://github.com/WOLSTN/wolstn) native compiler. It parses TypeScript code, performs full type checking, and emits a comprehensive typed IR (Intermediate Representation) that preserves all type information for native code generation.

## Overview

WTS is forked from [Microsoft's TypeScript-go](https://github.com/microsoft/TypeScript-go) (tsgo), but adapted to serve as a frontend for native compilation rather than JavaScript emission.

| | tsgo | WTS |
|---|---|---|
| **Purpose** | TypeScript → JavaScript | TypeScript → Wolstn IR |
| **Output** | JavaScript code | Typed IR (JSON) |
| **Type info** | Used for checking only | **Preserved for codegen** |
| **Runtime** | V8, Node.js, browsers | Native (ELF/EXE) via Wolstn |

## Why WTS?

Unlike Deno and Bun, which transpile TypeScript to JavaScript for execution in a JS runtime, **Wolstn compiles TypeScript directly to native machine code**. This requires preserving full type information for:

- **Memory layout decisions**: `number` → f64, `string` → GC pointer
- **Calling conventions**: Generic instantiation, overload resolution
- **Optimizations**: Inlining, escape analysis, type-driven optimizations

WTS provides this type information through its IR output, which is a complete lossless representation of the typed program.

## Installation

```bash
go build -o wts.exe ./cmd/wts
```

## Usage

### Type Check

```bash
wts check main.ts
wts check -p tsconfig.json
```

### Emit IR

```bash
wts emit-ir main.ts -o program.wir
wts emit-ir -p tsconfig.json -o program.wir
```

The output `.wir` file is a JSON-formatted IR containing all typed program information.

## IR Format

The IR is a versioned, self-contained representation of the typed program:

```json
{
  "version": 1,
  "files": [{ "path": "main.ts", "source": "..." }],
  "types": [{ "id": "t1", "kind": "object", "flags": 1048576, "properties": [...] }],
  "symbols": [{ "id": "s1", "name": "foo", "flags": 16, "kind": "function", "type": "t1" }],
  "signatures": [{ "id": "sig1", "kind": "call", "parameters": [...], "returnType": "t2" }],
  "globals": [{ "name": "foo", "symbol": "s1" }],
  "classes": [{ "name": "Animal", "symbol": "s2", "properties": [...], "methods": [...], "isAbstract": true }],
  "interfaces": [{ "name": "Person", "symbol": "s3", "properties": [...], "extends": [...] }],
  "enums": [{ "name": "Direction", "symbol": "s4", "members": [...], "isConst": false }],
  "typeAliases": [{ "name": "Point", "symbol": "s5", "target": "t10", "typeParams": [...] }],
  "namespaces": [{ "name": "MyNamespace", "symbol": "s6", "members": [...] }],
  "functions": [{ "name": "foo", "symbol": "s1", "parameters": [...], "returnType": "t2", "isAsync": false, "isGenerator": false }],
  "variables": [{ "name": "x", "symbol": "s7", "type": "t3", "isConst": true }],
  "imports": [{ "kind": "named", "modulePath": "fs", "symbol": "s8", "specifiers": [...] }],
  "exports": [{ "kind": "default", "name": "foo", "symbol": "s1", "isDefault": true }]
}
```

### IR Collections

| Collection | Description |
|---|---|
| `types` | All resolved types (primitives, objects, unions, generics, etc.) |
| `symbols` | All named declarations with flags and type references |
| `signatures` | Function/constructor/call signatures with parameters and return types |
| `globals` | Top-level declarations mapping names to symbols |
| `classes` | Class declarations with properties, methods, constructor, inheritance |
| `interfaces` | Interface declarations with properties, extends, call/construct signatures |
| `enums` | Enum declarations with members and values |
| `typeAliases` | Type alias declarations with target type and type parameters |
| `namespaces` | Namespace/module declarations with exported members |
| `functions` | Top-level function declarations with parameters, generics, async/generator flags |
| `variables` | Top-level variable declarations with types and const-ness |
| `imports` | Import declarations from external modules |
| `exports` | Export declarations including named, default, and re-exports |
| `files` | Source file paths and raw source text |

### Modifier Flags

Where applicable, IR structures carry modifier information:

- **Properties**: `isOptional`, `isReadonly`
- **Methods**: `isStatic`, `isAbstract`, `isPrivate`, `isProtected`, `isOptional`
- **Classes**: `isAbstract`
- **Enums**: `isConst` (const enum)
- **Functions**: `isAsync`, `isGenerator`

### Type Coverage

The IR covers all TypeScript type variants:

| Type Kind | Description |
|---|---|
| `any`, `unknown` | Top types |
| `string`, `number`, `boolean`, `bigint`, `void`, `undefined`, `null`, `never` | Primitive types |
| `stringLiteral`, `numberLiteral`, `booleanLiteral`, `literal` | Literal types |
| `union`, `intersection` | Compound types |
| `object` | Object types with properties, signatures, index infos |
| `typeParameter` | Generic type parameters with constraints and defaults |
| `index`, `indexedAccess` | Indexed access types |
| `conditional` | Conditional types (check, extends, true, false) |
| `templateLiteral` | Template literal types with text parts and embedded types |
| `stringMapping`, `substitution` | Advanced type manipulation |

## Architecture

```
Source Code (.ts, .tsx)
      ↓
   Parser (parser/)
      ↓
    AST (ast/)
      ↓
   Binder (binder/)
      ↓
  Checker (checker/)
      ↓
Typed AST + Full Type System
      ↓
 IR Emitter (ir/)
      ↓
Wolstn IR (.wir JSON)
```

### Data Collection

The IR emitter uses reflection (`CollectAllCheckerData`) to exhaustively collect all type maps, symbol tables, and link stores from the checker, eliminating manual per-field getter functions:

```
checker.Checker ──reflect──→ CheckerData
  ├── symbolArena       ──→ SymbolArena
  ├── signatureArena    ──→ SignatureArena
  ├── indexInfoArena    ──→ IndexInfoArena
  ├── globals           ──→ Globals
  ├── *map[string]*Type ──→ TypeMaps (15+ type maps)
  └── *LinkStore        ──→ LinkStores (25+ link stores)
```

## Module Structure

| Module | Purpose | Origin |
|--------|---------|--------|
| `ast` | AST definitions, symbols, flags | tsgo (adapted) |
| `parser` | Source code parsing | tsgo |
| `scanner` | Lexical analysis | tsgo |
| `binder` | Symbol binding and scoping | tsgo |
| `checker` | Type checking and type system | tsgo (adapted) |
| `ir` | **IR emission (new)** | WTS |
| `compiler` | Program orchestration | tsgo (adapted) |
| `cmd/wts` | CLI entry point | WTS |

Components removed from upstream tsgo (not needed for native compilation): JS emitter, printer, transformers, LSP server, language service, API server, formatter, fourslash tests.

## Development

```bash
# Build CLI
go build -o wts.exe ./cmd/wts

# Type check a file
wts check test.ts

# Emit IR
wts emit-ir test.ts -o test.wir

# Run tests
go test ./...
```

## License

Apache License 2.0 - see [LICENSE](LICENSE) for details.

## Acknowledgments

- Microsoft TypeScript Team for [TypeScript-go](https://github.com/microsoft/TypeScript-go)
- The Wolstn project
