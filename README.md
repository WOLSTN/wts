# WTS - TypeScript Frontend for Wolstn

**WTS** (Wolstn TypeScript) is a TypeScript frontend for the [Wolstn](https://github.com/WOLSTN) native compiler. It parses TypeScript code, performs full type checking, and emits a comprehensive typed IR (`.wir`) that preserves all type information for native code generation.

## Overview

WTS is forked from [Microsoft's TypeScript-go](https://github.com/microsoft/TypeScript-go) (tsgo), but adapted to serve as a frontend for native compilation rather than JavaScript emission.

> **WTS is not Wolstn itself. WTS only produces `.wir` IR files — it does not compile to native executables.**
> To compile TypeScript to a native binary, feed the `.wir` output to **[WolstnC](https://github.com/WOLSTN/wolstnC)** (the Wolstn compiler backend).

| | tsgo | WTS |
|---|---|---|
| **Purpose** | TypeScript → JavaScript | TypeScript → Wolstn IR |
| **Output** | JavaScript code | Typed IR (JSON) |
| **Type info** | Used for checking only | **Preserved for codegen** |
| **AST** | Only in encoder binary format | **Full tree in IR** |
| **Control flow** | Not exposed | **CFG with basic blocks** |
| **Runtime** | V8, Node.js, browsers | Native via WolstnC |

## Why WTS?

Unlike Deno and Bun, which transpile TypeScript to JavaScript for execution in a JS runtime, **Wolstn compiles TypeScript directly to native machine code**. This requires preserving full type information for:

- **Memory layout decisions**: `number` → f64, `string` → GC pointer
- **Calling conventions**: Generic instantiation, overload resolution
- **Optimizations**: Inlining, escape analysis, type-driven optimizations

WTS provides this type information through its IR output, which is a complete lossless representation of the typed program.

### Why not just use `tsgo --api`?

This is the most common question. tsgo has an `--api` flag that encodes the AST into a binary format. Why write a whole new frontend instead of just using that?

**Because `tsgo --api` outputs syntax, not semantics.**

| What you get | `tsgo --api` | WTS `emit-ir` |
|---|---|---|
| AST nodes | ✅ Full tree | ✅ Full tree with types |
| Source positions | ✅ Exact | ✅ Exact |
| String table | ✅ Deduplicated | ✅ Inline |
| **Type information** | ❌ **Not included** | ✅ **Full type system** |
| **Symbol table** | ❌ **Not included** | ✅ **All flags, parents, members** |
| **Function signatures** | ❌ **Not included** | ✅ **Parameters, returns, generics** |
| **Class/Interface structure** | ❌ Raw AST nodes only | ✅ **Properties, methods, modifiers** |
| **Function bodies** | ❌ AST only | ✅ **CFG with basic blocks + instructions** |
| **Expression types** | ❌ Not included | ✅ **Every AST node has a type** |
| **Import/Export semantics** | ❌ Raw AST nodes only | ✅ **Semantic records** |
| **Enum values** | ❌ Raw AST nodes only | ✅ **Resolved member values** |
| **Generic instantiations** | ❌ Not included | ✅ **All type arguments** |

To understand the gap concretely: feeding `tsgo --api` output to a backend would mean the backend must **reimplement a full TypeScript type checker** just to figure out what type each expression has. That's thousands of lines of type inference logic, relational algorithms, and constraint solving — exactly what tsgo's checker already does.

**WTS does the type checking once, in the frontend, and serializes the results.** The backend gets types for free.

#### "But you could just link tsgo as a library"

tsgo is **not designed to be used as a library for type data access**:

- All checker fields (type maps, symbol tables, link stores) are **unexported** — there are no getters for most of them
- The checker is designed to be a black box: parse → check → emit JS, with type information discarded after emission
- Exposing all type data would require adding hundreds of exported methods to tsgo's checker — which is the same amount of work as writing WTS's `CollectAllCheckerData`

WTS takes a different approach: it uses **Go reflection** to collect all checker data exhaustively in a single pass ([`CollectAllCheckerData`](internal/checker/allmaps.go)), then processes it into a clean, backend-friendly IR. No need to modify tsgo's checker internals.

#### "But you're forked from tsgo, maintenance burden!"

Yes, WTS is forked. But the fork is **shallow and stable**:

- WTS only uses **frontend components** from tsgo: parser, binder, checker
- These components are **stable** — the TypeScript grammar and type system don't change often
- The IR emitter (`internal/ir/`) is **new code** written by WTS, not modified tsgo code
- When upstream fixes a parser bug or adds a syntax feature, cherry-picking is straightforward because the AST and checker APIs are shared

The alternative — writing a TypeScript frontend from scratch or wrapping tsgo as a library — would be **far more work** with no meaningful benefit.

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

The output `.wir` file is a JSON-formatted IR containing the full typed program: type system, AST with expression-level types, symbol table, and function bodies with a control flow graph.

#### IR Options

| Flag | Description |
|---|---|
| `--prune` | Remove internal noise types (empty sentinels, zero-value stubs that carry no type information) |
| `--compact` | Compact JSON output — no indentation, smaller file size (~40% reduction) |
| `--no-source` | Omit source text from `File` entries — reduces file size and avoids embedding full source |

Combine flags as needed:

```bash
# Minimal, production-ready IR
wts emit-ir --prune --compact --no-source main.ts -o main.wir
```

## IR Format

The IR is a versioned, self-contained representation of the typed program:

```json
{
  "version": 1,
  "files": [
    {
      "path": "main.ts",
      "source": "...",
      "nodes": [{ "kind": "KindFunctionDeclaration", "pos": 0, "end": 50, "type": "t1", "children": [...] }]
    }
  ],
  "types": [{ "id": "t1", "kind": "object", "flags": 1048576, "properties": [...] }],
  "symbols": [{ "id": "s1", "name": "foo", "flags": 16, "kind": "function", "type": "t1" }],
  "signatures": [{ "id": "sig1", "kind": "call", "parameters": [...], "returnType": "t2" }],
  "globals": [{ "name": "foo", "symbol": "s1" }],
  "classes": [{ "name": "Animal", "symbol": "s2", "properties": [...], "methods": [...], "isAbstract": true }],
  "interfaces": [{ "name": "Person", "symbol": "s3", "properties": [...], "extends": [...] }],
  "enums": [{ "name": "Direction", "symbol": "s4", "members": [...], "isConst": false }],
  "typeAliases": [{ "name": "Point", "symbol": "s5", "target": "t10", "typeParams": [...] }],
  "namespaces": [{ "name": "MyNamespace", "symbol": "s6", "members": [...] }],
  "functions": [{
    "name": "foo",
    "symbol": "s1",
    "parameters": [...],
    "returnType": "t2",
    "isAsync": false,
    "body": {
      "blocks": [{
        "id": 1,
        "label": "entry",
        "preds": [],
        "succs": [2],
        "instrs": [
          { "id": "i1", "opcode": "literal", "type": "t3", "value": 42 },
          { "id": "i2", "opcode": "ret", "type": "t1", "operands": ["i1"] }
        ]
      }]
    }
  }],
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
| `functions` | Top-level function declarations with parameters, generics, async/generator flags, and CFG body |
| `variables` | Top-level variable declarations with types and const-ness |
| `imports` | Import declarations from external modules |
| `exports` | Export declarations including named, default, and re-exports |
| `files` | Source file paths, raw source text, and **full AST tree** |

### Function Body IR

Every function with a body exports a control flow graph of basic blocks:

```
Function
  └── body (FuncBody)
        └── blocks[] (BasicBlock)
              ├── id        — block identifier
              ├── label     — human-readable label (entry, if.then, while.cond, etc.)
              ├── preds[]   — predecessor block IDs
              ├── succs[]   — successor block IDs
              └── instrs[]  — instructions in SSA-like form
                    ├── id       — instruction identifier (i1, i2, ...)
                    ├── opcode   — operation (literal, ident, binary, store, ret, call, etc.)
                    ├── type     — result type ID
                    ├── operands — operand instruction IDs
                    └── value    — literal value (if applicable)
```

#### Supported Opcodes

| Opcode | Description | Opcode | Description |
|---|---|---|---|
| `literal` | Constant value | `ident` | Variable reference |
| `binary` | Binary operation | `unary` | Unary operation |
| `store` | Assignment | `alloc` | Variable allocation |
| `call` | Function call | `new` | Constructor call |
| `prop` | Property access | `elem` | Indexed access |
| `ret` | Return | `jmp` | Unconditional jump |
| `br` | Conditional branch | `inc` | Increment (++i / i++) |
| `dec` | Decrement (--i / i--) | `template` | Template string |
| `array` | Array literal | `object` | Object literal |
| `select` | Ternary expression | `func` | Function expression |
| `await` | Await expression | `yield` | Generator yield |
| `cast` | Type assertion | `typeof` | Typeof expression |
| `throw` | Throw statement | `compound` | Compound assignment (+=, -=) |
| `break` | Break statement | `continue` | Continue statement |
| `case` | Switch case start | `this` | This keyword |
| `super` | Super keyword | `spread` | Spread element |

### AST Tree in Files

Each file entry includes the complete AST tree as nested `nodes`, with type information on every node:

- Every AST node has a `type` field pointing to the type system
- Nodes that resolve to a symbol carry a `symbol` reference
- The tree structure mirrors the TypeScript AST (statements → expressions → literals)

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
      │
      ├── type system    → types[], symbols[], signatures[]
      ├── declarations   → classes[], interfaces[], enums[], functions[], variables[]
      ├── AST tree       → files[].nodes[] (every node has a type)
      └── CFG bodies     → functions[].body.blocks[] (basic blocks + instructions)
      │
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
