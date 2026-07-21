package ir

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/wolstn/wts/internal/ast"
	"github.com/wolstn/wts/internal/checker"
	"github.com/wolstn/wts/internal/compiler"
	"github.com/wolstn/wts/internal/jsnum"
)

const Version = 1

type Program struct {
	Version     int          `json:"version"`
	Files       []*File      `json:"files"`
	Types       []*Type      `json:"types"`
	Symbols     []*Symbol    `json:"symbols"`
	Signatures  []*Signature `json:"signatures"`
	Globals     []*Global    `json:"globals"`
	Classes     []*Class     `json:"classes"`
	Interfaces  []*Interface `json:"interfaces"`
	Enums       []*Enum      `json:"enums,omitempty"`
	TypeAliases []*TypeAlias `json:"typeAliases,omitempty"`
	Namespaces  []*Namespace `json:"namespaces,omitempty"`
	Imports     []*Import    `json:"imports,omitempty"`
	Exports     []*Export    `json:"exports,omitempty"`
	Functions   []*Function  `json:"functions,omitempty"`
	Variables   []*Variable  `json:"variables,omitempty"`
	// Runtime, when non-nil, carries the runtime binding descriptor authored on the
	// frontend (see RuntimeDescriptor). Serialized as the top-level "runtime" field
	// consumed by the wolstnc backend. Absent => backend applies its own defaults.
	Runtime *RuntimeDescriptor `json:"runtime,omitempty"`
}

type File struct {
	Path   string     `json:"path"`
	Source string     `json:"source,omitempty"`
	Nodes  []*ASTNode `json:"nodes,omitempty"`
}

type ASTNode struct {
	Kind     string     `json:"kind"`
	Pos      int        `json:"pos"`
	End      int        `json:"end"`
	Type     string     `json:"type,omitempty"`
	Symbol   string     `json:"symbol,omitempty"`
	Children []*ASTNode `json:"children,omitempty"`
}

type Type struct {
	Id           string       `json:"id"`
	Kind         string       `json:"kind"`
	Flags        uint32       `json:"flags"`
	Name         string       `json:"name,omitempty"`
	ObjectFlags  uint32       `json:"objectFlags,omitempty"`
	Value        any          `json:"value,omitempty"`
	Types        []string     `json:"types,omitempty"`
	Target       string       `json:"target,omitempty"`
	Constraint   string       `json:"constraint,omitempty"`
	Default      string       `json:"default,omitempty"`
	TypeArgs     []string     `json:"typeArgs,omitempty"`
	TypeParams   []string     `json:"typeParams,omitempty"`
	ObjectType   string       `json:"objectType,omitempty"`
	IndexType    string       `json:"indexType,omitempty"`
	CheckType    string       `json:"checkType,omitempty"`
	ExtendsType  string       `json:"extendsType,omitempty"`
	TrueType     string       `json:"trueType,omitempty"`
	FalseType    string       `json:"falseType,omitempty"`
	Properties   []*Property  `json:"properties,omitempty"`
	Signatures   []*SigRef    `json:"signatures,omitempty"`
	BaseTypes    []string     `json:"baseTypes,omitempty"`
	IndexInfos   []*IndexInfo `json:"indexInfos,omitempty"`
	IsReadonly   bool         `json:"isReadonly,omitempty"`
	IsThisType   bool         `json:"isThisType,omitempty"`
	Texts        []string     `json:"texts,omitempty"`
	ElementFlags []string     `json:"elementFlags,omitempty"`
}

type Property struct {
	Name       string `json:"name"`
	Symbol     string `json:"symbol,omitempty"`
	Type       string `json:"type,omitempty"`
	IsOptional bool   `json:"isOptional,omitempty"`
	IsReadonly bool   `json:"isReadonly,omitempty"`
}

type SigRef struct {
	Id string `json:"id"`
}

type IndexInfo struct {
	KeyType   string `json:"keyType"`
	ValueType string `json:"valueType"`
	IsReadonly bool  `json:"isReadonly,omitempty"`
}

type Signature struct {
	Id             string   `json:"id"`
	Kind           string   `json:"kind"`
	Parameters     []Param  `json:"parameters,omitempty"`
	ReturnType     string   `json:"returnType,omitempty"`
	TypeParameters []string `json:"typeParameters,omitempty"`
	Declaration    string   `json:"declaration,omitempty"`
	HasRest        bool     `json:"hasRest,omitempty"`
}

type Param struct {
	Name   string `json:"name"`
	Type   string `json:"type,omitempty"`
	Symbol string `json:"symbol,omitempty"`
}

type Symbol struct {
	Id           string   `json:"id"`
	Name         string   `json:"name"`
	Flags        uint64   `json:"flags"`
	CheckFlags   uint32   `json:"checkFlags,omitempty"`
	Kind         string   `json:"kind"`
	Type         string   `json:"type,omitempty"`
	Parent       string   `json:"parent,omitempty"`
	Declarations []string `json:"declarations,omitempty"`
	Members      []string `json:"members,omitempty"`
	Exports      []string `json:"exports,omitempty"`
}

type Global struct {
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

type Class struct {
	Name         string      `json:"name"`
	Symbol       string      `json:"symbol"`
	TypeParams   []string    `json:"typeParams,omitempty"`
	BaseClass    string      `json:"baseClass,omitempty"`
	Implements   []string    `json:"implements,omitempty"`
	Properties   []*Property `json:"properties,omitempty"`
	Methods      []*Method   `json:"methods,omitempty"`
	Constructor  *Method     `json:"constructor,omitempty"`
	IsAbstract   bool        `json:"isAbstract,omitempty"`
}

type Interface struct {
	Name         string      `json:"name"`
	Symbol       string      `json:"symbol"`
	TypeParams   []string    `json:"typeParams,omitempty"`
	Extends      []string    `json:"extends,omitempty"`
	Properties   []*Property `json:"properties,omitempty"`
	Methods      []*Method   `json:"methods,omitempty"`
	CallSigs     []*SigRef   `json:"callSigs,omitempty"`
	ConstructSigs []*SigRef  `json:"constructSigs,omitempty"`
}

type Method struct {
	Name         string  `json:"name"`
	Symbol       string  `json:"symbol"`
	ReturnType   string  `json:"returnType,omitempty"`
	Parameters   []Param `json:"parameters,omitempty"`
	TypeParams   []string `json:"typeParams,omitempty"`
	IsStatic     bool    `json:"isStatic,omitempty"`
	IsAbstract   bool    `json:"isAbstract,omitempty"`
	IsPrivate    bool    `json:"isPrivate,omitempty"`
	IsProtected  bool    `json:"isProtected,omitempty"`
	IsOptional   bool    `json:"isOptional,omitempty"`
	IsAsync      bool    `json:"isAsync,omitempty"`
	IsGenerator  bool    `json:"isGenerator,omitempty"`
	Signature    string  `json:"signature,omitempty"`
	VtableIndex  *int    `json:"vtableIndex,omitempty"`
}

type Enum struct {
	Name       string       `json:"name"`
	Symbol     string       `json:"symbol"`
	Members    []*EnumMember `json:"members,omitempty"`
	IsConst    bool         `json:"isConst,omitempty"`
}

type EnumMember struct {
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
	Value  any    `json:"value,omitempty"`
	Type   string `json:"type,omitempty"`
}

type TypeAlias struct {
	Name       string   `json:"name"`
	Symbol     string   `json:"symbol"`
	TypeParams []string `json:"typeParams,omitempty"`
	Target     string   `json:"target,omitempty"`
}

type Namespace struct {
	Name    string          `json:"name"`
	Symbol  string          `json:"symbol"`
	Members []*NamespaceMember `json:"members,omitempty"`
}

type NamespaceMember struct {
	Kind   string `json:"kind"`
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

type Import struct {
	Kind        string          `json:"kind"`
	ModulePath  string          `json:"modulePath"`
	Symbol      string          `json:"symbol,omitempty"`
	Specifiers  []*ImportSpec   `json:"specifiers,omitempty"`
	Namespace   string          `json:"namespace,omitempty"`
	IsTypeOnly  bool            `json:"isTypeOnly,omitempty"`
}

type ImportSpec struct {
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

type Export struct {
	Kind      string `json:"kind"`
	Name      string `json:"name"`
	Symbol    string `json:"symbol"`
	IsDefault bool   `json:"isDefault,omitempty"`
}

type Function struct {
	Name       string     `json:"name"`
	Symbol     string     `json:"symbol"`
	Parameters []Param    `json:"parameters,omitempty"`
	ReturnType string     `json:"returnType,omitempty"`
	TypeParams []string   `json:"typeParams,omitempty"`
	Signature  string     `json:"signature,omitempty"`
	Body       *FuncBody  `json:"body,omitempty"`
	IsAsync    bool       `json:"isAsync,omitempty"`
	IsGenerator bool      `json:"isGenerator,omitempty"`
	IsExport   bool       `json:"isExport,omitempty"`
	IsExtern   bool       `json:"isExtern,omitempty"`
	IsDeclare  bool       `json:"isDeclare,omitempty"`
}

type FuncBody struct {
	Blocks []*BasicBlock `json:"blocks,omitempty"`
}

type BasicBlock struct {
	Id     int      `json:"id"`
	Label  string   `json:"label,omitempty"`
	Preds  []int    `json:"preds,omitempty"`
	Succs  []int    `json:"succs,omitempty"`
	Instrs []*Instr `json:"instrs,omitempty"`
}

type Instr struct {
	Id       string   `json:"id"`
	Opcode   string   `json:"opcode"`
	Type     string   `json:"type,omitempty"`
	Operands []string `json:"operands,omitempty"`
	Value    any      `json:"value,omitempty"`
}

type Variable struct {
	Name    string `json:"name"`
	Symbol  string `json:"symbol"`
	Type    string `json:"type,omitempty"`
	IsConst bool   `json:"isConst,omitempty"`
}

type EmitOptions struct {
	Prune     bool // Filter out internal noise types
	Compact   bool // Compact JSON output (no indentation, omit empty arrays)
	NoSource  bool // Omit source text from File entries
	TreeShake bool // Only emit types/symbols reachable from user code
	// Runtime, when non-nil, is serialized into the emitted WIR as the top-level
	// "runtime" field. It enables a TS author to declare the runtime binding
	// descriptor on the frontend instead of injecting it into the WIR afterwards
	// (noStdRoadMap N.3 / miss-wts-002).
	Runtime *RuntimeDescriptor
}

// RuntimeDescriptor is the WIR-side "runtime binding" descriptor. It mirrors the
// `RuntimeDescriptor` defined by the wolstnc backend (wolstnc/src/ir/program.rs) exactly
// at the JSON level, so that a TS author can author it on the *frontend* side.
//
// Every field uses camelCase JSON tags identical to the backend's serde(rename_all =
// "camelCase") layout, and `omitempty` so that omitted fields fall back to the backend's
// own defaults (e.g. `arcEnabled` defaults to true, `methodAliases` defaults to
// {"console.log": "wolstn_console_log"}). This lets the author override *only* what they
// need — a core tenet of the de-privileged design (AGENTS.md §8).
//
// `ArcEnabled` is a *pointer* so that "unset" (nil) is distinguishable from an explicit
// `false`, preserving the backend default when the author does not mention it.
type RuntimeDescriptor struct {
	MethodAliases     map[string]string `json:"methodAliases,omitempty"`
	PrintNumberSymbol string            `json:"printNumberSymbol,omitempty"`
	PrintBoolSymbol   string            `json:"printBoolSymbol,omitempty"`
	ArcEnabled        *bool             `json:"arcEnabled,omitempty"`
	AllocSymbol       string            `json:"allocSymbol,omitempty"`
	RetainSymbol      string            `json:"retainSymbol,omitempty"`
	ReleaseSymbol     string            `json:"releaseSymbol,omitempty"`
}

type typeState struct {
	id        string
	resolving bool
}

type Emitter struct {
	program       *compiler.Program
	checker       *checker.Checker
	checkerData   *checker.CheckerData
	irProgram     *Program
	typeState     map[any]*typeState
	sigState      map[*checker.Signature]*typeState
	symbolMap     map[*ast.Symbol]string
	symTypeCache  map[*ast.Symbol]string
	emittedSyms    map[*ast.Symbol]bool
	emittedClasses map[*ast.Symbol]bool
	typeIdGen      int
	symbolIdGen    int
	sigIdGen       int
	options        EmitOptions
}

func NewEmitter(program *compiler.Program, opts EmitOptions) *Emitter {
	return &Emitter{
		program:     program,
		irProgram: &Program{
			Version: Version,
			// Initialize every top-level sequence field to an empty (non-nil) slice.
			// The wolstnc backend's WIR parser requires these to be JSON arrays;
			// a nil slice marshals as `null` and is rejected with
			// "invalid type: null, expected a sequence". This happens whenever a
			// program emits no classes/interfaces/etc. (e.g. `-p` project mode with
			// `noLib: true`), so the slices must be pre-seeded rather than relying on
			// emit* helpers to populate them.
			Files:      []*File{},
			Types:      []*Type{},
			Symbols:    []*Symbol{},
			Signatures: []*Signature{},
			Globals:    []*Global{},
			Classes:    []*Class{},
			Interfaces: []*Interface{},
		},
		typeState:   make(map[any]*typeState),
		sigState:    make(map[*checker.Signature]*typeState),
		symbolMap:   make(map[*ast.Symbol]string),
		symTypeCache:   make(map[*ast.Symbol]string),
		emittedSyms:    make(map[*ast.Symbol]bool),
		emittedClasses: make(map[*ast.Symbol]bool),
		options:        opts,
	}
}

func (e *Emitter) Emit() (*Program, error) {
	ctx := context.Background()
	chk, release := e.program.GetTypeChecker(ctx)
	defer release()
	e.checker = chk

	e.checkerData = checker.CollectAllCheckerData(e.checker)

	if e.checkerData.Globals != nil {
		for name, sym := range e.checkerData.Globals {
			e.irProgram.Globals = append(e.irProgram.Globals, &Global{
				Name:   name,
				Symbol: e.getOrCreateSymbolId(sym),
			})
		}
	}

	e.emitSourceFiles()

	if e.options.Prune {
		e.pruneTypes()
	}

	if e.options.TreeShake {
		e.treeShake()
	}

	// Propagate a frontend-authored runtime descriptor into the emitted WIR so the
	// backend consumes it directly (no post-injection needed). nil => backend defaults.
	if e.options.Runtime != nil {
		e.irProgram.Runtime = e.options.Runtime
	}

	return e.irProgram, nil
}

func (e *Emitter) pruneTypes() {
	isNoise := func(t *Type) bool {
		// Empty PseudoBigInt sentinel (Base10Value: "")
		if t.Kind == "literal" {
			if pbi, ok := t.Value.(jsnum.PseudoBigInt); ok && !pbi.Negative && pbi.Base10Value == "" {
				return true
			}
		}
		// Empty template literal with no meaningful content
		if t.Kind == "templateLiteral" && len(t.Texts) == 2 && t.Texts[0] == "" && t.Texts[1] == "" {
			hasOnlyNumberType := len(t.Types) == 1
			if hasOnlyNumberType {
				return true
			}
		}
		return false
	}

	noiseIds := make(map[string]bool)
	for _, t := range e.irProgram.Types {
		if isNoise(t) {
			noiseIds[t.Id] = true
		}
	}

	if len(noiseIds) == 0 {
		return
	}

	filtered := make([]*Type, 0, len(e.irProgram.Types))
	for _, t := range e.irProgram.Types {
		if !noiseIds[t.Id] {
			filtered = append(filtered, t)
		}
	}
	e.irProgram.Types = filtered
}

func (e *Emitter) emitSourceFiles() {
	var allTopLevelStmts []*ast.Node
	for _, sf := range e.program.GetSourceFiles() {
		if sf.IsDeclarationFile {
			continue
		}
		file := &File{Path: sf.FileName()}
		if !e.options.NoSource {
			file.Source = string(sf.Text())
		}
		e.irProgram.Files = append(e.irProgram.Files, file)
		file.Nodes = e.emitNodeTree(sf)
		e.emitFileImports(sf)
		e.emitFileExports(sf)
		e.emitFileDeclarations(sf)
		
		// Collect top level statements from this file
		for _, stmt := range sf.Statements.Nodes {
			switch stmt.Kind {
			case ast.KindExpressionStatement, ast.KindVariableStatement:
				allTopLevelStmts = append(allTopLevelStmts, stmt)
			}
		}
	}

	e.emitAllTopLevelStatements(allTopLevelStmts)

	for _, sym := range e.checkerData.Globals {
		if e.emittedSyms[sym] {
			continue
		}
		e.emittedSyms[sym] = true
		flags := sym.Flags

		switch {
		case flags&ast.SymbolFlagsClass != 0:
			e.emitClassFromSymbol(sym)
		case flags&ast.SymbolFlagsInterface != 0:
			e.emitInterfaceFromSymbol(sym)
		case flags&ast.SymbolFlagsEnum != 0:
			e.emitEnumFromSymbol(sym)
		case flags&ast.SymbolFlagsTypeAlias != 0:
			e.emitTypeAliasFromSymbol(sym, sym.Name)
		case flags&ast.SymbolFlagsModule != 0:
			if sym.Name != "globalThis" {
				e.emitNamespaceFromSymbol(sym)
			}
		case flags&ast.SymbolFlagsFunction != 0:
			e.emitFunctionFromSymbol(sym)
		case flags&ast.SymbolFlagsVariable != 0:
			e.emitVariableFromSymbol(sym)
		}
	}
}

func (e *Emitter) emitAllTopLevelStatements(topLevelStmts []*ast.Node) {
	if len(topLevelStmts) == 0 {
		return
	}

	var voidTypeId string
	for _, t := range e.irProgram.Types {
		if t.Kind == "void" {
			voidTypeId = t.Id
			break
		}
	}

	hasMainCall := false
	for _, stmt := range topLevelStmts {
		if stmt.Kind != ast.KindExpressionStatement {
			continue
		}
		expr := stmt.Expression()
		if expr != nil && expr.Kind == ast.KindCallExpression {
			callee := expr.Expression()
			if callee != nil && callee.Kind == ast.KindIdentifier && callee.Text() == "main" {
				hasMainCall = true
				break
			}
		}
	}

	var originalMain *Function
	if hasMainCall {
		for i, f := range e.irProgram.Functions {
			if f.Name == "main" && f.Body != nil {
				originalMain = f
				e.irProgram.Functions = append(e.irProgram.Functions[:i], e.irProgram.Functions[i+1:]...)
				break
			}
		}
		if originalMain != nil {
			userMain := &Function{
				Name:       "__wolstn_user_main",
				Symbol:     "s_user_main",
				Parameters: originalMain.Parameters,
				ReturnType: originalMain.ReturnType,
				Body:       originalMain.Body,
			}
			e.irProgram.Functions = append(e.irProgram.Functions, userMain)
			e.irProgram.Symbols = append(e.irProgram.Symbols, &Symbol{
				Id:   "s_user_main",
				Name: "__wolstn_user_main",
				Kind: "function",
			})
		}
	}

	be := &bodyEmitter{e: e}
	entry := be.addBlock("entry")
	be.curBlock = entry

	if hasMainCall && originalMain != nil {
		userMainIdentId := be.addInstr("ident", "", "__wolstn_user_main", nil)
		be.addInstr("call", voidTypeId, nil, []string{userMainIdentId})
	}

	for _, stmt := range topLevelStmts {
		if stmt.Kind == ast.KindExpressionStatement {
			expr := stmt.Expression()
			if expr != nil && expr.Kind == ast.KindCallExpression {
				callee := expr.Expression()
				if callee != nil && callee.Kind == ast.KindIdentifier && callee.Text() == "main" {
					continue
				}
			}
		}
		be.emitStatement(stmt)
		if be.curBlock == nil {
			break
		}
	}

	if len(be.body.Blocks) == 0 {
		return
	}

	fn := &Function{
		Name:       "__wolstn_toplevel",
		Symbol:     "s_toplevel",
		Body:       be.body,
		ReturnType: voidTypeId,
	}
	e.irProgram.Functions = append(e.irProgram.Functions, fn)

	be2 := &bodyEmitter{e: e}
	entry2 := be2.addBlock("entry")
	be2.curBlock = entry2
	toplevelIdentId := be2.addInstr("ident", "", "__wolstn_toplevel", nil)
	be2.addInstr("call", voidTypeId, nil, []string{toplevelIdentId})

	entryFn := &Function{
		Name:       "__wolstn_entry",
		Symbol:     "s_entry",
		ReturnType: voidTypeId,
		Body:       be2.body,
	}
	e.irProgram.Functions = append(e.irProgram.Functions, entryFn)

	e.irProgram.Symbols = append(e.irProgram.Symbols, &Symbol{
		Id:   "s_toplevel",
		Name: "__wolstn_toplevel",
		Kind: "function",
	})
	e.irProgram.Symbols = append(e.irProgram.Symbols, &Symbol{
		Id:   "s_entry",
		Name: "__wolstn_entry",
		Kind: "function",
	})
}

func (e *Emitter) emitFileImports(sf *ast.SourceFile) {
	for _, stmt := range sf.Statements.Nodes {
		if stmt.Kind == ast.KindImportDeclaration {
			importDecl := stmt.AsImportDeclaration()
			modulePath := ""
			if importDecl.ModuleSpecifier != nil {
				if sl := importDecl.ModuleSpecifier.AsStringLiteral(); sl != nil {
					modulePath = sl.Text
				}
			}
			if importDecl.ImportClause != nil {
				ic := importDecl.ImportClause.AsImportClause()
				isTypeOnly := ic.PhaseModifier == ast.KindTypeKeyword
				if ic.Name() != nil {
					sym := ic.AsNode().Symbol()
					e.irProgram.Imports = append(e.irProgram.Imports, &Import{
						Kind:       "default",
						ModulePath: modulePath,
						Symbol:     e.getOrCreateSymbolId(sym),
						IsTypeOnly: isTypeOnly,
					})
				}
				if ic.NamedBindings != nil {
					nb := ic.NamedBindings
					switch nb.Kind {
					case ast.KindNamespaceImport:
						ns := nb.AsNamespaceImport()
						sym := ns.AsNode().Symbol()
						e.irProgram.Imports = append(e.irProgram.Imports, &Import{
							Kind:       "namespace",
							ModulePath: modulePath,
							Symbol:     e.getOrCreateSymbolId(sym),
							Namespace:  ns.Name().Text(),
							IsTypeOnly: isTypeOnly,
						})
					case ast.KindNamedImports:
						named := nb.AsNamedImports()
						for _, elem := range named.Elements.Nodes {
							spec := elem.AsImportSpecifier()
							sym := spec.AsNode().Symbol()
							specName := spec.Name().Text()
							if spec.PropertyName != nil {
								specName = spec.PropertyName.Text()
							}
							e.irProgram.Imports = append(e.irProgram.Imports, &Import{
								Kind:       "named",
								ModulePath: modulePath,
								Symbol:     e.getOrCreateSymbolId(sym),
								IsTypeOnly: isTypeOnly || spec.IsTypeOnly,
								Specifiers: []*ImportSpec{{
									Name:   specName,
									Symbol: e.getOrCreateSymbolId(sym),
								}},
							})
						}
					}
				}
			}
		}
	}
}

func (e *Emitter) emitFileExports(sf *ast.SourceFile) {
	for _, stmt := range sf.Statements.Nodes {
		switch stmt.Kind {
		case ast.KindExportDeclaration:
			exportDecl := stmt.AsExportDeclaration()
			if exportDecl.ExportClause != nil {
				ec := exportDecl.ExportClause
				if ec.Kind == ast.KindNamedExports {
					named := ec.AsNamedExports()
					for _, elem := range named.Elements.Nodes {
						spec := elem.AsExportSpecifier()
						sym := spec.AsNode().Symbol()
						e.irProgram.Exports = append(e.irProgram.Exports, &Export{
							Kind:      "named",
							Name:      spec.Name().Text(),
							Symbol:    e.getOrCreateSymbolId(sym),
							IsDefault: false,
						})
					}
				}
			} else if exportDecl.ModuleSpecifier != nil {
				modulePath := ""
				if sl := exportDecl.ModuleSpecifier.AsStringLiteral(); sl != nil {
					modulePath = sl.Text
				}
				e.irProgram.Exports = append(e.irProgram.Exports, &Export{
					Kind:   "reexport",
					Name:   modulePath,
					Symbol: "",
				})
			}
		case ast.KindExportAssignment:
			exportAssign := stmt.AsExportAssignment()
			if !exportAssign.IsExportEquals {
				sym := stmt.Symbol()
				name := "default"
				if exportAssign.Expression != nil {
					if id := exportAssign.Expression.AsIdentifier(); id != nil {
						name = id.Text
					}
				}
				e.irProgram.Exports = append(e.irProgram.Exports, &Export{
					Kind:      "default",
					Name:      name,
					Symbol:    e.getOrCreateSymbolId(sym),
					IsDefault: true,
				})
			}
		}
	}
}

func (e *Emitter) emitFileDeclarations(sf *ast.SourceFile) {
	for _, stmt := range sf.Statements.Nodes {
		if stmt.Kind == ast.KindVariableStatement {
			e.emitVariableStatementDeclarations(stmt)
			continue
		}
		sym := stmt.Symbol()
		if sym == nil {
			sym = e.checker.GetSymbolAtLocation(stmt)
		}
		if sym == nil && stmt.Kind == ast.KindFunctionDeclaration {
			funcDecl := stmt.AsFunctionDeclaration()
			if name := funcDecl.Name(); name != nil {
				sym = e.checker.GetSymbolAtLocation(name)
			}
		}
		if sym == nil {
			continue
		}
		if e.emittedSyms[sym] {
			continue
		}
		e.emittedSyms[sym] = true
		flags := sym.Flags
		switch {
		case flags&ast.SymbolFlagsTypeAlias != 0:
			e.emitTypeAliasFromSymbol(sym, sym.Name)
		case flags&ast.SymbolFlagsClass != 0:
			e.emitClassFromSymbol(sym)
		case flags&ast.SymbolFlagsInterface != 0:
			e.emitInterfaceFromSymbol(sym)
		case flags&ast.SymbolFlagsEnum != 0:
			e.emitEnumFromSymbol(sym)
		case flags&ast.SymbolFlagsFunction != 0:
			e.emitFunctionFromSymbol(sym)
		case flags&ast.SymbolFlagsVariable != 0:
			e.emitVariableFromSymbol(sym)
		case flags&ast.SymbolFlagsModule != 0:
			if sym.Name != "globalThis" {
				e.emitNamespaceFromSymbol(sym)
			}
		}
	}

	if sf.Locals != nil {
		for _, sym := range sf.Locals {
			if e.emittedSyms[sym] {
				continue
			}
			e.emittedSyms[sym] = true
			flags := sym.Flags
			switch {
			case flags&ast.SymbolFlagsFunction != 0:
				e.emitFunctionFromSymbol(sym)
			case flags&ast.SymbolFlagsVariable != 0:
				e.emitVariableFromSymbol(sym)
			case flags&ast.SymbolFlagsClass != 0:
				e.emitClassFromSymbol(sym)
			case flags&ast.SymbolFlagsInterface != 0:
				e.emitInterfaceFromSymbol(sym)
			case flags&ast.SymbolFlagsEnum != 0:
				e.emitEnumFromSymbol(sym)
			case flags&ast.SymbolFlagsTypeAlias != 0:
				e.emitTypeAliasFromSymbol(sym, sym.Name)
			}
		}
	}

	if sf.Symbol != nil && sf.Symbol.Exports != nil {
		for _, sym := range sf.Symbol.Exports {
			if e.emittedSyms[sym] {
				continue
			}
			e.emittedSyms[sym] = true
			flags := sym.Flags
			switch {
			case flags&ast.SymbolFlagsFunction != 0:
				e.emitFunctionFromSymbol(sym)
			case flags&ast.SymbolFlagsVariable != 0:
				e.emitVariableFromSymbol(sym)
			case flags&ast.SymbolFlagsClass != 0:
				e.emitClassFromSymbol(sym)
			case flags&ast.SymbolFlagsInterface != 0:
				e.emitInterfaceFromSymbol(sym)
			case flags&ast.SymbolFlagsEnum != 0:
				e.emitEnumFromSymbol(sym)
			case flags&ast.SymbolFlagsTypeAlias != 0:
				e.emitTypeAliasFromSymbol(sym, sym.Name)
			}
		}
	}

	for _, sym := range e.checkerData.Symbols {
		if sym == nil {
			continue
		}
		if e.emittedSyms[sym] {
			continue
		}
		if len(sym.Declarations) == 0 {
			continue
		}
		decl := sym.Declarations[0]
		if decl == nil {
			continue
		}
		sourceFile := ast.GetSourceFileOfNode(decl)
		if sourceFile == nil {
			if sf.Symbol != nil && sym.Parent == sf.Symbol {
				sourceFile = sf
			} else {
				for _, d := range sym.Declarations {
					if d != nil && d.Pos() >= 0 && d.End() <= len(string(sf.Text())) {
						sourceFile = sf
						break
					}
				}
			}
		}
		if sourceFile != sf {
			continue
		}
		e.emittedSyms[sym] = true
		flags := sym.Flags
		switch {
		case flags&ast.SymbolFlagsFunction != 0:
			e.emitFunctionFromSymbol(sym)
		case flags&ast.SymbolFlagsVariable != 0:
			e.emitVariableFromSymbol(sym)
		case flags&ast.SymbolFlagsClass != 0:
			e.emitClassFromSymbol(sym)
		case flags&ast.SymbolFlagsInterface != 0:
			e.emitInterfaceFromSymbol(sym)
		case flags&ast.SymbolFlagsEnum != 0:
			e.emitEnumFromSymbol(sym)
		case flags&ast.SymbolFlagsTypeAlias != 0:
			e.emitTypeAliasFromSymbol(sym, sym.Name)
		}
	}
}

func (e *Emitter) emitFunctionDeclaration(stmt *ast.Node) {
	sym := stmt.Symbol()
	if sym == nil {
		sym = e.checker.GetSymbolAtLocation(stmt)
	}
	if sym == nil {
		funcDecl := stmt.AsFunctionDeclaration()
		name := funcDecl.Name()
		if name != nil {
			sym = e.checker.GetSymbolAtLocation(name)
		}
	}
	if sym == nil {
		return
	}
	if e.emittedSyms[sym] {
		return
	}
	e.emittedSyms[sym] = true
	e.emitFunctionFromSymbol(sym)
}

func (e *Emitter) emitVariableStatementDeclarations(stmt *ast.Node) {
	stmt.ForEachChild(func(child *ast.Node) bool {
		if child.Kind == ast.KindVariableDeclarationList {
			child.ForEachChild(func(decl *ast.Node) bool {
				if decl.Kind == ast.KindVariableDeclaration {
					sym := decl.Symbol()
					if sym != nil && !e.emittedSyms[sym] {
						e.emittedSyms[sym] = true
						e.emitVariableFromSymbol(sym)
					}
				}
				return false
			})
		}
		return false
	})
}

func (e *Emitter) emitClassFromSymbol(sym *ast.Symbol) {
	if e.emittedClasses[sym] {
		return
	}
	e.emittedClasses[sym] = true

	irClass := &Class{
		Name:   sym.Name,
		Symbol: e.getOrCreateSymbolId(sym),
	}

	declaredType := e.checker.GetDeclaredTypeOfSymbol(sym)
	var baseTypes []*checker.Type
	if declaredType != nil {
		typeParams := e.checker.GetLocalTypeParametersOfClassOrInterfaceOrTypeAlias(sym)
		for _, tp := range typeParams {
			irClass.TypeParams = append(irClass.TypeParams, e.getOrCreateTypeId(tp))
		}

		baseTypes = e.checker.GetBaseTypes(declaredType)
		for _, bt := range baseTypes {
			if bt.Symbol() != nil {
				if bt.Symbol().Flags&ast.SymbolFlagsClass != 0 {
					irClass.BaseClass = e.getOrCreateSymbolId(bt.Symbol())
				} else if bt.Symbol().Flags&ast.SymbolFlagsInterface != 0 {
					irClass.Implements = append(irClass.Implements, e.getOrCreateSymbolId(bt.Symbol()))
				}
			}
		}

		props := e.checker.GetPropertiesOfType(declaredType)
		for _, prop := range props {
			e.emitClassProperty(prop, irClass)
		}

		if irClass.Constructor == nil {
			irClass.Constructor = e.emitConstructor(sym, declaredType)
		}

		implSet := make(map[string]bool)
		for _, id := range irClass.Implements {
			implSet[id] = true
		}
		for _, decl := range sym.Declarations {
			if decl.Kind == ast.KindClassDeclaration {
				for _, impl := range ast.GetImplementsHeritageClauseElements(decl) {
					implType := e.checker.GetTypeFromTypeNode(impl.AsNode())
					if implType != nil && implType.Symbol() != nil {
						id := e.getOrCreateSymbolId(implType.Symbol())
						if !implSet[id] {
							implSet[id] = true
							irClass.Implements = append(irClass.Implements, id)
						}
					}
				}
			}
		}
	}

	if sym.Flags&ast.SymbolFlagsClass != 0 {
		for _, decl := range sym.Declarations {
			if decl.Kind == ast.KindClassDeclaration {
				classDecl := decl.AsClassDeclaration()
				if mods := classDecl.Modifiers(); mods != nil {
					for _, mod := range mods.Nodes {
						if mod.Kind == ast.KindAbstractKeyword {
							irClass.IsAbstract = true
						}
					}
				}
			}
		}
	}

	// Calculate vtable index for non-static methods
	var parentClass *Class
	for _, bt := range baseTypes {
		if bt.Symbol() != nil && bt.Symbol().Flags&ast.SymbolFlagsClass != 0 {
			parentSym := bt.Symbol()
			if !e.emittedSyms[parentSym] {
				e.emitClassFromSymbol(parentSym)
			}
			parentSymId := e.getOrCreateSymbolId(parentSym)
			for _, c := range e.irProgram.Classes {
				if c.Symbol == parentSymId {
					parentClass = c
					break
				}
			}
		}
	}

	parentVtable := make(map[string]int)
	nextIndex := 0
	if parentClass != nil {
		for _, m := range parentClass.Methods {
			if !m.IsStatic && m.VtableIndex != nil {
				parentVtable[m.Name] = *m.VtableIndex
				if *m.VtableIndex >= nextIndex {
					nextIndex = *m.VtableIndex + 1
				}
			}
		}
	}

	for _, m := range irClass.Methods {
		if m.IsStatic {
			continue
		}
		if parentIndex, ok := parentVtable[m.Name]; ok {
			idx := parentIndex
			m.VtableIndex = &idx
		} else {
			idx := nextIndex
			m.VtableIndex = &idx
			nextIndex++
		}
	}

	e.irProgram.Classes = append(e.irProgram.Classes, irClass)
}

func (e *Emitter) emitConstructorFromNode(ctorNode *ast.Node, classSym *ast.Symbol, method *Method) {
	if sym := ctorNode.Symbol(); sym != nil {
		method.Symbol = e.getOrCreateSymbolId(sym)
	}

	ctorDecl := ctorNode.AsConstructorDeclaration()
	if ctorDecl == nil {
		return
	}

	for _, param := range ctorDecl.Parameters.Nodes {
		if param.Kind == ast.KindParameter {
			pDecl := param.AsParameterDeclaration()
			if pDecl == nil {
				continue
			}
			p := Param{Name: pDecl.Name().Text(), Symbol: e.getOrCreateSymbolId(pDecl.Symbol)}
			if pDecl.Type != nil {
				t := e.checker.GetTypeFromTypeNode(pDecl.Type)
				if t != nil {
					p.Type = e.getOrCreateTypeId(t)
				}
			}
			method.Parameters = append(method.Parameters, p)
		}
	}

	sig := e.checker.GetSignatureFromDeclaration(ctorNode)
	if sig != nil {
		method.Signature = e.getOrCreateSignatureId(sig)
		if ret := e.checker.GetReturnTypeOfSignature(sig); ret != nil {
			method.ReturnType = e.getOrCreateTypeId(ret)
		}
		for _, tp := range sig.TypeParameters() {
			method.TypeParams = append(method.TypeParams, e.getOrCreateTypeId(tp))
		}
	}

	body := e.emitFunctionBody(ctorNode)
	if body != nil && len(body.Blocks) > 0 {
		fullName := classSym.Name + ".constructor"
		fn := &Function{
			Name:       fullName,
			Symbol:     method.Symbol,
			Body:       body,
			ReturnType: method.ReturnType,
		}
		for _, p := range method.Parameters {
			fn.Parameters = append(fn.Parameters, Param{Name: p.Name, Type: p.Type, Symbol: p.Symbol})
		}
		e.irProgram.Functions = append(e.irProgram.Functions, fn)
	}
}

func (e *Emitter) emitConstructor(classSym *ast.Symbol, declaredType *checker.Type) *Method {
	method := &Method{
		Name: "constructor",
	}

	for _, decl := range classSym.Declarations {
		if decl.Kind != ast.KindClassDeclaration {
			continue
		}
		classDecl := decl.AsClassDeclaration()
		if mods := classDecl.Modifiers(); mods != nil {
			for _, mod := range mods.Nodes {
				if mod.Kind == ast.KindAbstractKeyword {
					method.IsAbstract = true
				}
			}
		}
		if classDecl.Members != nil && classDecl.Members.Nodes != nil {
			for _, member := range classDecl.Members.Nodes {
				if member.Kind == ast.KindConstructor {
					e.emitConstructorFromNode(member, classSym, method)
					break
				}
			}
		}
	}

	return method
}

func (e *Emitter) emitClassProperty(prop *ast.Symbol, irClass *Class) {
	flags := prop.Flags

	switch {
	case flags&ast.SymbolFlagsMethod != 0:
		method := e.emitMethod(prop)
		irClass.Methods = append(irClass.Methods, method)
	case flags&ast.SymbolFlagsConstructor != 0:
		irClass.Constructor = e.emitMethod(prop)
	default:
		propType := e.checker.GetTypeOfSymbol(prop)
		irProp := &Property{
			Name:       prop.Name,
			Symbol:     e.getOrCreateSymbolId(prop),
			IsOptional: prop.Flags&ast.SymbolFlagsOptional != 0,
		}
		if propType != nil {
			irProp.Type = e.getOrCreateTypeId(propType)
		}
		if len(prop.Declarations) > 0 {
			irProp.IsReadonly = ast.HasSyntacticModifier(prop.Declarations[0], ast.ModifierFlagsReadonly)
		}
		irClass.Properties = append(irClass.Properties, irProp)
	}
}

func (e *Emitter) emitMethod(sym *ast.Symbol) *Method {
	method := &Method{
		Name:   sym.Name,
		Symbol: e.getOrCreateSymbolId(sym),
	}

	if len(sym.Declarations) > 0 {
		decl := sym.Declarations[0]
		method.IsStatic = ast.HasSyntacticModifier(decl, ast.ModifierFlagsStatic)
		method.IsAbstract = ast.HasSyntacticModifier(decl, ast.ModifierFlagsAbstract)
		method.IsPrivate = ast.HasSyntacticModifier(decl, ast.ModifierFlagsPrivate)
		method.IsProtected = ast.HasSyntacticModifier(decl, ast.ModifierFlagsProtected)
		method.IsAsync = ast.HasSyntacticModifier(decl, ast.ModifierFlagsAsync)
		if decl.Kind == ast.KindMethodDeclaration {
			method.IsGenerator = decl.AsMethodDeclaration().AsteriskToken != nil
		}
	}

	if len(sym.Declarations) > 0 {
		decl := sym.Declarations[0]
		if decl.Kind == ast.KindMethodDeclaration || decl.Kind == ast.KindGetAccessor || decl.Kind == ast.KindSetAccessor || decl.Kind == ast.KindConstructor {
			funcType := e.checker.GetTypeOfSymbol(sym)
			if funcType != nil {
				sigs := e.checker.GetCallSignatures(funcType)
				if len(sigs) > 0 {
					sig := sigs[0]
					method.Signature = e.getOrCreateSignatureId(sig)
					for _, param := range sig.Parameters() {
						p := Param{Name: param.Name, Symbol: e.getOrCreateSymbolId(param)}
						if t := e.checker.GetTypeOfSymbol(param); t != nil {
							p.Type = e.getOrCreateTypeId(t)
						}
						method.Parameters = append(method.Parameters, p)
					}
					if ret := e.checker.GetReturnTypeOfSignature(sig); ret != nil {
						method.ReturnType = e.getOrCreateTypeId(ret)
					}
					for _, tp := range sig.TypeParameters() {
						method.TypeParams = append(method.TypeParams, e.getOrCreateTypeId(tp))
					}
				}
			}
		}
	}

	if len(sym.Declarations) > 0 {
		decl := sym.Declarations[0]
		body := e.emitFunctionBody(decl)
		if body != nil && len(body.Blocks) > 0 {
			className := ""
			if sym.Parent != nil {
				className = sym.Parent.Name
			}
			fullName := sym.Name
			if className != "" {
				fullName = className + "." + sym.Name
			}
			fn := &Function{
				Name:       fullName,
				Symbol:     method.Symbol,
				Body:       body,
				ReturnType: method.ReturnType,
				IsAsync:    method.IsAsync,
				IsGenerator: method.IsGenerator,
			}
			for _, p := range method.Parameters {
				fn.Parameters = append(fn.Parameters, Param{Name: p.Name, Type: p.Type, Symbol: p.Symbol})
			}
			e.irProgram.Functions = append(e.irProgram.Functions, fn)
		}
	}

	return method
}

func (e *Emitter) emitInterfaceFromSymbol(sym *ast.Symbol) {
	irIntf := &Interface{
		Name:   sym.Name,
		Symbol: e.getOrCreateSymbolId(sym),
	}

	declaredType := e.checker.GetDeclaredTypeOfSymbol(sym)
	if declaredType != nil {
		typeParams := e.checker.GetLocalTypeParametersOfClassOrInterfaceOrTypeAlias(sym)
		for _, tp := range typeParams {
			irIntf.TypeParams = append(irIntf.TypeParams, e.getOrCreateTypeId(tp))
		}

		baseTypes := e.checker.GetBaseTypes(declaredType)
		for _, bt := range baseTypes {
			if bt.Symbol() != nil && bt.Symbol().Flags&ast.SymbolFlagsInterface != 0 {
				irIntf.Extends = append(irIntf.Extends, e.getOrCreateSymbolId(bt.Symbol()))
			}
		}

		props := e.checker.GetPropertiesOfType(declaredType)
		for _, prop := range props {
			propType := e.checker.GetTypeOfSymbol(prop)
			irProp := &Property{
				Name:       prop.Name,
				Symbol:     e.getOrCreateSymbolId(prop),
				IsOptional: prop.Flags&ast.SymbolFlagsOptional != 0,
			}
			if propType != nil {
				irProp.Type = e.getOrCreateTypeId(propType)
			}
			if len(prop.Declarations) > 0 {
				irProp.IsReadonly = ast.HasSyntacticModifier(prop.Declarations[0], ast.ModifierFlagsReadonly)
			}
			irIntf.Properties = append(irIntf.Properties, irProp)
		}

		callSigs := e.checker.GetCallSignatures(declaredType)
		for _, sig := range callSigs {
			irIntf.CallSigs = append(irIntf.CallSigs, &SigRef{Id: e.getOrCreateSignatureId(sig)})
		}

		constructSigs := e.checker.GetConstructSignatures(declaredType)
		for _, sig := range constructSigs {
			irIntf.ConstructSigs = append(irIntf.ConstructSigs, &SigRef{Id: e.getOrCreateSignatureId(sig)})
		}
	}

	e.irProgram.Interfaces = append(e.irProgram.Interfaces, irIntf)
}

func (e *Emitter) emitEnumFromSymbol(sym *ast.Symbol) {
	irEnum := &Enum{
		Name:   sym.Name,
		Symbol: e.getOrCreateSymbolId(sym),
	}

	if len(sym.Declarations) > 0 {
		irEnum.IsConst = ast.IsEnumConst(sym.Declarations[0])
	}

	if sym.Members != nil {
		for name, member := range sym.Members {
			memberType := e.checker.GetTypeOfSymbol(member)
			irMember := &EnumMember{
				Name:   name,
				Symbol: e.getOrCreateSymbolId(member),
			}
			if memberType != nil {
				irMember.Type = e.getOrCreateTypeId(memberType)
			}
			if len(member.Declarations) > 0 {
				irMember.Value = e.checker.GetEnumMemberValue(member.Declarations[0])
			}
			irEnum.Members = append(irEnum.Members, irMember)
		}
	}

	e.irProgram.Enums = append(e.irProgram.Enums, irEnum)
}

func (e *Emitter) emitTypeAliasFromSymbol(sym *ast.Symbol, name string) {
	irAlias := &TypeAlias{
		Name:   name,
		Symbol: e.getOrCreateSymbolId(sym),
	}

	typeParams := e.checker.GetLocalTypeParametersOfClassOrInterfaceOrTypeAlias(sym)
	for _, tp := range typeParams {
		irAlias.TypeParams = append(irAlias.TypeParams, e.getOrCreateTypeId(tp))
	}

	if len(sym.Declarations) > 0 {
		decl := sym.Declarations[0]
		if decl.Kind == ast.KindTypeAliasDeclaration {
			typeAliasDecl := decl.AsTypeAliasDeclaration()
			if typeAliasDecl.Type != nil {
				t := e.checker.GetTypeFromTypeNode(typeAliasDecl.Type)
				if t != nil {
					irAlias.Target = e.getOrCreateTypeId(t)
				}
			}
		}
	}

	e.irProgram.TypeAliases = append(e.irProgram.TypeAliases, irAlias)
}

func (e *Emitter) emitNamespaceFromSymbol(sym *ast.Symbol) {
	irNs := &Namespace{
		Name:   sym.Name,
		Symbol: e.getOrCreateSymbolId(sym),
	}

	if sym.Exports != nil {
		for name, exported := range sym.Exports {
			member := &NamespaceMember{
				Name:   name,
				Symbol: e.getOrCreateSymbolId(exported),
			}
			switch {
			case exported.Flags&ast.SymbolFlagsFunction != 0:
				member.Kind = "function"
			case exported.Flags&ast.SymbolFlagsClass != 0:
				member.Kind = "class"
			case exported.Flags&ast.SymbolFlagsInterface != 0:
				member.Kind = "interface"
			case exported.Flags&ast.SymbolFlagsEnum != 0:
				member.Kind = "enum"
			case exported.Flags&ast.SymbolFlagsModule != 0:
				member.Kind = "namespace"
			case exported.Flags&ast.SymbolFlagsTypeAlias != 0:
				member.Kind = "typeAlias"
			case exported.Flags&ast.SymbolFlagsVariable != 0:
				member.Kind = "variable"
			default:
				member.Kind = "unknown"
			}
			irNs.Members = append(irNs.Members, member)
		}
	}

	e.irProgram.Namespaces = append(e.irProgram.Namespaces, irNs)
}

func (e *Emitter) emitFunctionFromSymbol(sym *ast.Symbol) {
	fn := &Function{
		Name:   sym.Name,
		Symbol: e.getOrCreateSymbolId(sym),
	}

	if len(sym.Declarations) > 0 {
		decl := sym.Declarations[0]
		fn.Body = e.emitFunctionBody(decl)

		if decl.Kind == ast.KindFunctionDeclaration || decl.Kind == ast.KindFunctionExpression {
			if decl.Kind == ast.KindFunctionDeclaration {
				fn.IsAsync = ast.HasSyntacticModifier(decl, ast.ModifierFlagsAsync)
				fn.IsGenerator = decl.AsFunctionDeclaration().AsteriskToken != nil
				fn.IsExport = ast.HasSyntacticModifier(decl, ast.ModifierFlagsExport)
				fn.IsExtern = ast.HasSyntacticModifier(decl, ast.ModifierFlagsAmbient)
			}
			funcType := e.checker.GetTypeOfSymbol(sym)
			if funcType != nil {
				sigs := e.checker.GetCallSignatures(funcType)
				if len(sigs) > 0 {
					sig := sigs[0]
					fn.Signature = e.getOrCreateSignatureId(sig)
					for _, param := range sig.Parameters() {
						p := Param{Name: param.Name, Symbol: e.getOrCreateSymbolId(param)}
						if t := e.checker.GetTypeOfSymbol(param); t != nil {
							p.Type = e.getOrCreateTypeId(t)
						}
						fn.Parameters = append(fn.Parameters, p)
					}
					if ret := e.checker.GetReturnTypeOfSignature(sig); ret != nil {
						fn.ReturnType = e.getOrCreateTypeId(ret)
					}
					for _, tp := range sig.TypeParameters() {
						fn.TypeParams = append(fn.TypeParams, e.getOrCreateTypeId(tp))
					}
				}
			}
		}
	}

	e.irProgram.Functions = append(e.irProgram.Functions, fn)
}

type bodyEmitter struct {
	e          *Emitter
	body       *FuncBody
	curBlock   *BasicBlock
	blockIdGen int
	instrIdGen int
	blockStack []*BasicBlock // for if/while/loop target blocks
}

func (be *bodyEmitter) newBlockId() int {
	be.blockIdGen++
	return be.blockIdGen
}

func (be *bodyEmitter) newInstrId() string {
	be.instrIdGen++
	return fmt.Sprintf("i%d", be.instrIdGen)
}

func (be *bodyEmitter) addBlock(label string) *BasicBlock {
	block := &BasicBlock{
		Id:    be.newBlockId(),
		Label: label,
	}
	if be.body == nil {
		be.body = &FuncBody{}
	}
	be.body.Blocks = append(be.body.Blocks, block)
	return block
}

func (be *bodyEmitter) addInstr(opcode string, typ string, value any, operands []string) string {
	id := be.newInstrId()
	instr := &Instr{
		Id:       id,
		Opcode:   opcode,
		Type:     typ,
		Value:    value,
		Operands: operands,
	}
	be.curBlock.Instrs = append(be.curBlock.Instrs, instr)
	return id
}

func (be *bodyEmitter) connectBlocks(pred, succ *BasicBlock) {
	if pred != nil && succ != nil {
		pred.Succs = append(pred.Succs, succ.Id)
		succ.Preds = append(succ.Preds, pred.Id)
	}
}

func (e *Emitter) emitFunctionBody(decl *ast.Node) *FuncBody {
	bodyNode := decl.Body()
	if bodyNode == nil {
		return nil
	}

	be := &bodyEmitter{e: e}
	entry := be.addBlock("entry")
	be.curBlock = entry

	be.emitBlockStatements(bodyNode)

	if len(be.body.Blocks) == 0 {
		return nil
	}
	return be.body
}

func (be *bodyEmitter) emitBlockStatements(block *ast.Node) {
	if block == nil {
		return
	}
	stmtList := block.Statements()
	if stmtList == nil {
		return
	}
	for _, stmt := range stmtList {
		be.emitStatement(stmt)
		if be.curBlock == nil {
			return
		}
	}
}

func (be *bodyEmitter) emitStatement(stmt *ast.Node) {
	if stmt == nil {
		return
	}
	switch stmt.Kind {
	case ast.KindVariableStatement:
		be.emitVariableStatement(stmt)
	case ast.KindReturnStatement:
		be.emitReturnStatement(stmt)
	case ast.KindExpressionStatement:
		be.emitExpressionStatement(stmt)
	case ast.KindIfStatement:
		be.emitIfStatement(stmt)
	case ast.KindWhileStatement:
		be.emitWhileStatement(stmt)
	case ast.KindDoStatement:
		be.emitDoStatement(stmt)
	case ast.KindForStatement:
		be.emitForStatement(stmt)
	case ast.KindForInStatement:
		be.emitForInStatement(stmt)
	case ast.KindForOfStatement:
		be.emitForOfStatement(stmt)
	case ast.KindSwitchStatement:
		be.emitSwitchStatement(stmt)
	case ast.KindBlock:
		be.emitBlockStatements(stmt)
	case ast.KindBreakStatement:
		be.addInstr("break", "", nil, nil)
	case ast.KindContinueStatement:
		be.addInstr("continue", "", nil, nil)
	case ast.KindEmptyStatement:
		// nothing
	case ast.KindTryStatement:
		be.emitTryStatement(stmt)
	case ast.KindThrowStatement:
		be.emitThrowStatement(stmt)
	case ast.KindDebuggerStatement:
		// nothing
	default:
		be.emitExpression(stmt)
	}
}

func (be *bodyEmitter) emitVariableStatement(stmt *ast.Node) {
	stmt.ForEachChild(func(child *ast.Node) bool {
		if child.Kind == ast.KindVariableDeclarationList {
			declList := child.AsVariableDeclarationList()
			for _, decl := range declList.Declarations.Nodes {
				be.emitVariableDeclaration(decl)
			}
		}
		return false
	})
}

func (be *bodyEmitter) emitVariableDeclaration(decl *ast.Node) {
	varDecl := decl.AsVariableDeclaration()
	name := varDecl.Name()
	if name == nil {
		return
	}
	sym := decl.Symbol()
	targetId := name.Text()
	if sym != nil {
		targetId = be.e.getOrCreateSymbolId(sym)
	}
	initializer := varDecl.Initializer
	id := ""
	if initializer != nil {
		valId := be.emitExpression(initializer)
		id = be.addInstr("store", be.e.getNodeType(decl), targetId, []string{valId})
	} else {
		id = be.addInstr("alloc", be.e.getNodeType(decl), targetId, nil)
	}
	_ = id
}

func (be *bodyEmitter) emitReturnStatement(stmt *ast.Node) {
	expr := stmt.Expression()
	if expr != nil {
		valId := be.emitExpression(expr)
		be.addInstr("ret", be.e.getNodeType(stmt), nil, []string{valId})
	} else {
		be.addInstr("ret", "", nil, nil)
	}
}

func (be *bodyEmitter) emitExpressionStatement(stmt *ast.Node) {
	expr := stmt.Expression()
	if expr != nil {
		be.emitExpression(expr)
	}
}

func (be *bodyEmitter) emitThrowStatement(stmt *ast.Node) {
	expr := stmt.Expression()
	if expr != nil {
		valId := be.emitExpression(expr)
		be.addInstr("throw", be.e.getNodeType(stmt), nil, []string{valId})
	} else {
		be.addInstr("throw", "", nil, nil)
	}
}

func (be *bodyEmitter) emitIfStatement(stmt *ast.Node) {
	cond := stmt.Expression()
	thenBlock := stmt.AsIfStatement().ThenStatement
	elseBlock := stmt.AsIfStatement().ElseStatement

	condId := be.emitExpression(cond)

	thenBB := be.addBlock("if.then")
	elseBB := be.addBlock("if.else")
	endBB := be.addBlock("if.end")

	be.addInstr("br", "", nil, []string{condId, thenBB.Label, elseBB.Label})
	be.connectBlocks(be.curBlock, thenBB)
	be.connectBlocks(be.curBlock, elseBB)

	be.curBlock = thenBB
	if thenBlock != nil {
		if thenBlock.Kind == ast.KindBlock {
			be.emitBlockStatements(thenBlock)
		} else {
			be.emitStatement(thenBlock)
		}
	}
	be.addInstr("jmp", "", nil, []string{endBB.Label})
	be.connectBlocks(be.curBlock, endBB)

	be.curBlock = elseBB
	if elseBlock != nil {
		if elseBlock.Kind == ast.KindBlock {
			be.emitBlockStatements(elseBlock)
		} else {
			be.emitStatement(elseBlock)
		}
	}
	be.addInstr("jmp", "", nil, []string{endBB.Label})
	be.connectBlocks(be.curBlock, endBB)

	be.curBlock = endBB
}

func (be *bodyEmitter) emitWhileStatement(stmt *ast.Node) {
	condBB := be.addBlock("while.cond")
	bodyBB := be.addBlock("while.body")
	endBB := be.addBlock("while.end")

	be.addInstr("jmp", "", nil, []string{condBB.Label})
	be.connectBlocks(be.curBlock, condBB)

	be.curBlock = condBB
	condId := be.emitExpression(stmt.Expression())
	be.addInstr("br", "", nil, []string{condId, bodyBB.Label, endBB.Label})
	be.connectBlocks(condBB, bodyBB)
	be.connectBlocks(condBB, endBB)

	be.curBlock = bodyBB
	body := stmt.AsWhileStatement().Statement
	if body != nil {
		if body.Kind == ast.KindBlock {
			be.emitBlockStatements(body)
		} else {
			be.emitStatement(body)
		}
	}
	be.addInstr("jmp", "", nil, []string{condBB.Label})
	be.connectBlocks(bodyBB, condBB)

	be.curBlock = endBB
}

func (be *bodyEmitter) emitDoStatement(stmt *ast.Node) {
	bodyBB := be.addBlock("do.body")
	condBB := be.addBlock("do.cond")
	endBB := be.addBlock("do.end")

	be.addInstr("jmp", "", nil, []string{bodyBB.Label})
	be.connectBlocks(be.curBlock, bodyBB)

	be.curBlock = bodyBB
	body := stmt.AsDoStatement().Statement
	if body != nil {
		if body.Kind == ast.KindBlock {
			be.emitBlockStatements(body)
		} else {
			be.emitStatement(body)
		}
	}
	be.addInstr("jmp", "", nil, []string{condBB.Label})
	be.connectBlocks(bodyBB, condBB)

	be.curBlock = condBB
	condId := be.emitExpression(stmt.Expression())
	be.addInstr("br", "", nil, []string{condId, bodyBB.Label, endBB.Label})
	be.connectBlocks(condBB, bodyBB)
	be.connectBlocks(condBB, endBB)

	be.curBlock = endBB
}

func (be *bodyEmitter) emitForStatement(stmt *ast.Node) {
	forStmt := stmt.AsForStatement()
	initBB := be.addBlock("for.init")
	condBB := be.addBlock("for.cond")
	bodyBB := be.addBlock("for.body")
	incrBB := be.addBlock("for.increment")
	endBB := be.addBlock("for.end")

	be.addInstr("jmp", "", nil, []string{initBB.Label})
	be.connectBlocks(be.curBlock, initBB)

	be.curBlock = initBB
	if initializer := forStmt.Initializer; initializer != nil {
		if initializer.Kind == ast.KindVariableDeclarationList {
			declList := initializer.AsVariableDeclarationList()
			for _, decl := range declList.Declarations.Nodes {
				be.emitVariableDeclaration(decl)
			}
		} else {
			be.emitExpression(initializer)
		}
	}
	be.addInstr("jmp", "", nil, []string{condBB.Label})
	be.connectBlocks(initBB, condBB)

	be.curBlock = condBB
	if cond := forStmt.Condition; cond != nil {
		condId := be.emitExpression(cond)
		be.addInstr("br", "", nil, []string{condId, bodyBB.Label, endBB.Label})
	} else {
		be.addInstr("jmp", "", nil, []string{bodyBB.Label})
	}
	be.connectBlocks(condBB, bodyBB)
	be.connectBlocks(condBB, endBB)

	be.curBlock = bodyBB
	body := forStmt.Statement
	if body != nil {
		if body.Kind == ast.KindBlock {
			be.emitBlockStatements(body)
		} else {
			be.emitStatement(body)
		}
	}
	be.addInstr("jmp", "", nil, []string{incrBB.Label})
	be.connectBlocks(bodyBB, incrBB)

	be.curBlock = incrBB
	if incrementor := forStmt.Incrementor; incrementor != nil {
		be.emitExpression(incrementor)
	}
	be.addInstr("jmp", "", nil, []string{condBB.Label})
	be.connectBlocks(incrBB, condBB)

	be.curBlock = endBB
}

func (be *bodyEmitter) emitSwitchStatement(stmt *ast.Node) {
	condId := be.emitExpression(stmt.Expression())
	endBB := be.addBlock("switch.end")

	switchBlock := stmt.AsSwitchStatement().CaseBlock.AsCaseBlock()
	for _, clause := range switchBlock.Clauses.Nodes {
		bb := be.addBlock("switch.case")
		be.addInstr("case", "", nil, []string{condId})
		if clause.Kind == ast.KindCaseClause {
			caseClause := clause.AsCaseOrDefaultClause()
			caseValId := be.emitExpression(caseClause.Expression)
			be.addInstr("case.match", "", nil, []string{caseValId})
		}
		for _, cs := range clause.Statements() {
			be.emitStatement(cs)
		}
		be.connectBlocks(be.curBlock, bb)
		be.curBlock = bb
	}
	be.addInstr("jmp", "", nil, []string{endBB.Label})
	be.connectBlocks(be.curBlock, endBB)
	be.curBlock = endBB
}

func (be *bodyEmitter) emitForInStatement(stmt *ast.Node) {
	forIn := stmt.AsForInOrOfStatement()

	initBB := be.addBlock("forin.init")
	condBB := be.addBlock("forin.cond")
	bodyBB := be.addBlock("forin.body")
	endBB := be.addBlock("forin.end")

	be.addInstr("jmp", "", nil, []string{initBB.Label})
	be.connectBlocks(be.curBlock, initBB)

	be.curBlock = initBB

	if forIn.Initializer != nil {
		if forIn.Initializer.Kind == ast.KindVariableDeclarationList {
			declList := forIn.Initializer.AsVariableDeclarationList()
			for _, decl := range declList.Declarations.Nodes {
				be.emitVariableDeclaration(decl)
			}
		}
	}

	objId := be.emitExpression(forIn.Expression)
	iterId := be.addInstr("iter.keys", be.e.getNodeType(forIn.Expression), nil, []string{objId})

	be.addInstr("jmp", "", nil, []string{condBB.Label})
	be.connectBlocks(initBB, condBB)

	be.curBlock = condBB
	hasNextId := be.addInstr("iter.has_next", "", nil, []string{iterId})
	be.addInstr("br", "", nil, []string{hasNextId, bodyBB.Label, endBB.Label})
	be.connectBlocks(condBB, bodyBB)
	be.connectBlocks(condBB, endBB)

	be.curBlock = bodyBB
	keyId := be.addInstr("iter.next", "", nil, []string{iterId})

	if forIn.Initializer != nil {
		if forIn.Initializer.Kind == ast.KindVariableDeclarationList {
			declList := forIn.Initializer.AsVariableDeclarationList()
			for _, decl := range declList.Declarations.Nodes {
				name := decl.AsVariableDeclaration().Name()
				if name != nil {
					nameText := name.Text()
					be.addInstr("store", be.e.getNodeType(decl), nameText, []string{keyId})
				}
			}
		}
	}

	body := forIn.Statement
	if body != nil {
		if body.Kind == ast.KindBlock {
			be.emitBlockStatements(body)
		} else {
			be.emitStatement(body)
		}
	}
	be.addInstr("jmp", "", nil, []string{condBB.Label})
	be.connectBlocks(bodyBB, condBB)

	be.curBlock = endBB
}

func (be *bodyEmitter) emitForOfStatement(stmt *ast.Node) {
	forOf := stmt.AsForInOrOfStatement()

	initBB := be.addBlock("forof.init")
	condBB := be.addBlock("forof.cond")
	bodyBB := be.addBlock("forof.body")
	endBB := be.addBlock("forof.end")

	be.addInstr("jmp", "", nil, []string{initBB.Label})
	be.connectBlocks(be.curBlock, initBB)

	be.curBlock = initBB

	if forOf.Initializer != nil {
		if forOf.Initializer.Kind == ast.KindVariableDeclarationList {
			declList := forOf.Initializer.AsVariableDeclarationList()
			for _, decl := range declList.Declarations.Nodes {
				be.emitVariableDeclaration(decl)
			}
		}
	}

	objId := be.emitExpression(forOf.Expression)

	var iterTyp string
	if forOf.AwaitModifier != nil {
		iterId := be.addInstr("iter.async_values", be.e.getNodeType(forOf.Expression), nil, []string{objId})
		iterTyp = iterId
	} else {
		iterTyp = be.addInstr("iter.values", be.e.getNodeType(forOf.Expression), nil, []string{objId})
	}

	be.addInstr("jmp", "", nil, []string{condBB.Label})
	be.connectBlocks(initBB, condBB)

	be.curBlock = condBB
	hasNextId := be.addInstr("iter.has_next", "", nil, []string{iterTyp})
	be.addInstr("br", "", nil, []string{hasNextId, bodyBB.Label, endBB.Label})
	be.connectBlocks(condBB, bodyBB)
	be.connectBlocks(condBB, endBB)

	be.curBlock = bodyBB
	valId := be.addInstr("iter.next", "", nil, []string{iterTyp})

	if forOf.Initializer != nil {
		if forOf.Initializer.Kind == ast.KindVariableDeclarationList {
			declList := forOf.Initializer.AsVariableDeclarationList()
			for _, decl := range declList.Declarations.Nodes {
				name := decl.AsVariableDeclaration().Name()
				if name != nil {
					nameText := name.Text()
					be.addInstr("store", be.e.getNodeType(decl), nameText, []string{valId})
				}
			}
		}
	}

	body := forOf.Statement
	if body != nil {
		if body.Kind == ast.KindBlock {
			be.emitBlockStatements(body)
		} else {
			be.emitStatement(body)
		}
	}
	be.addInstr("jmp", "", nil, []string{condBB.Label})
	be.connectBlocks(bodyBB, condBB)

	be.curBlock = endBB
}

func (be *bodyEmitter) emitTryStatement(stmt *ast.Node) {
	tryStmt := stmt.AsTryStatement()

	tryBB := be.addBlock("try.body")
	catchBB := be.addBlock("try.catch")
	finallyBB := be.addBlock("try.finally")
	endBB := be.addBlock("try.end")

	hasCatch := tryStmt.CatchClause != nil
	hasFinally := tryStmt.FinallyBlock != nil

	var tryTargets []string
	if hasCatch {
		tryTargets = append(tryTargets, catchBB.Label)
	} else if hasFinally {
		tryTargets = append(tryTargets, finallyBB.Label)
	}
	be.addInstr("try.begin", "", nil, tryTargets)
	be.addInstr("jmp", "", nil, []string{tryBB.Label})
	be.connectBlocks(be.curBlock, tryBB)

	be.curBlock = tryBB
	if tryStmt.TryBlock != nil {
		be.emitBlockStatements(tryStmt.TryBlock)
	}
	be.addInstr("try.end", "", nil, nil)

	if hasFinally {
		be.addInstr("jmp", "", nil, []string{finallyBB.Label})
		be.connectBlocks(be.curBlock, finallyBB)
	} else {
		be.addInstr("jmp", "", nil, []string{endBB.Label})
		be.connectBlocks(be.curBlock, endBB)
	}

	if hasCatch {
		be.curBlock = catchBB
		be.addInstr("catch.begin", "", nil, nil)

		catchClause := tryStmt.CatchClause.AsCatchClause()
		if catchClause.VariableDeclaration != nil {
			be.emitVariableDeclaration(catchClause.VariableDeclaration)
		}
		if catchClause.Block != nil {
			be.emitBlockStatements(catchClause.Block)
		}

		be.addInstr("catch.end", "", nil, nil)
		if hasFinally {
			be.addInstr("jmp", "", nil, []string{finallyBB.Label})
			be.connectBlocks(be.curBlock, finallyBB)
		} else {
			be.addInstr("jmp", "", nil, []string{endBB.Label})
			be.connectBlocks(be.curBlock, endBB)
		}
	}

	if hasFinally {
		be.curBlock = finallyBB
		be.addInstr("finally.begin", "", nil, nil)
		be.emitBlockStatements(tryStmt.FinallyBlock)
		be.addInstr("finally.end", "", nil, nil)
		be.addInstr("jmp", "", nil, []string{endBB.Label})
		be.connectBlocks(be.curBlock, endBB)
	}

	be.curBlock = endBB
}

func (be *bodyEmitter) emitExpression(node *ast.Node) string {
	if node == nil {
		return ""
	}

	typ := be.e.getNodeType(node)

	switch node.Kind {
	case ast.KindNumericLiteral:
		return be.addInstr("literal", typ, node.AsNumericLiteral().Text, nil)
	case ast.KindStringLiteral:
		return be.addInstr("literal", typ, node.Text(), nil)
	case ast.KindFalseKeyword:
		return be.addInstr("literal", typ, false, nil)
	case ast.KindTrueKeyword:
		return be.addInstr("literal", typ, true, nil)
	case ast.KindNullKeyword:
		return be.addInstr("literal", typ, nil, nil)
	case ast.KindIdentifier:
		text := node.Text()
		sym := be.e.checker.GetResolvedSymbolOfNode(node)
		if sym != nil {
			text = be.e.getOrCreateSymbolId(sym)
		}
		return be.addInstr("ident", typ, text, nil)
	case ast.KindThisKeyword:
		return be.addInstr("this", typ, nil, nil)
	case ast.KindSuperKeyword:
		return be.addInstr("super", typ, nil, nil)
	case ast.KindBinaryExpression:
		return be.emitBinaryExpression(node, typ)
	case ast.KindPrefixUnaryExpression:
		return be.emitUnaryExpression(node, typ)
	case ast.KindPostfixUnaryExpression:
		return be.emitPostfixExpression(node, typ)
	case ast.KindCallExpression:
		return be.emitCallExpression(node, typ)
	case ast.KindPropertyAccessExpression:
		return be.emitPropertyAccess(node, typ)
	case ast.KindElementAccessExpression:
		return be.emitElementAccess(node, typ)
	case ast.KindNewExpression:
		return be.emitNewExpression(node, typ)
	case ast.KindFunctionExpression:
		return be.addInstr("func", typ, nil, nil)
	case ast.KindArrowFunction:
		return be.addInstr("func", typ, nil, nil)
	case ast.KindArrayLiteralExpression:
		return be.emitArrayLiteral(node, typ)
	case ast.KindObjectLiteralExpression:
		return be.emitObjectLiteral(node, typ)
	case ast.KindTemplateExpression:
		return be.emitTemplateExpression(node, typ)
	case ast.KindConditionalExpression:
		return be.emitConditionalExpression(node, typ)
	case ast.KindTypeOfExpression:
		exprId := be.emitExpression(node.AsTypeOfExpression().Expression)
		return be.addInstr("typeof", typ, nil, []string{exprId})
	case ast.KindAwaitExpression:
		exprId := be.emitExpression(node.AsAwaitExpression().Expression)
		return be.addInstr("await", typ, nil, []string{exprId})
	case ast.KindYieldExpression:
		yield := node.AsYieldExpression()
		if yield.Expression != nil {
			exprId := be.emitExpression(yield.Expression)
			return be.addInstr("yield", typ, nil, []string{exprId})
		}
		return be.addInstr("yield", typ, nil, nil)
	case ast.KindVoidExpression:
		exprId := be.emitExpression(node.AsVoidExpression().Expression)
		return be.addInstr("void", typ, nil, []string{exprId})
	case ast.KindDeleteExpression:
		exprId := be.emitExpression(node.AsDeleteExpression().Expression)
		return be.addInstr("delete", typ, nil, []string{exprId})
	case ast.KindSpreadElement:
		exprId := be.emitExpression(node.AsSpreadElement().Expression)
		return be.addInstr("spread", typ, nil, []string{exprId})
	case ast.KindAsExpression, ast.KindTypeAssertionExpression, ast.KindSatisfiesExpression:
		expr := node.Expression()
		exprId := be.emitExpression(expr)
		return be.addInstr("cast", typ, nil, []string{exprId})
	case ast.KindNonNullExpression:
		exprId := be.emitExpression(node.Expression())
		return be.addInstr("notnull", typ, nil, []string{exprId})
	case ast.KindParenthesizedExpression:
		return be.emitExpression(node.AsParenthesizedExpression().Expression)
	case ast.KindTypeReference:
		return be.addInstr("type", typ, nil, nil)
	case ast.KindRegularExpressionLiteral:
		return be.addInstr("literal", typ, node.Text(), nil)
	case ast.KindNoSubstitutionTemplateLiteral:
		return be.addInstr("literal", typ, node.Text(), nil)
	case ast.KindTemplateHead, ast.KindTemplateMiddle, ast.KindTemplateTail:
		return be.addInstr("literal", typ, node.Text(), nil)
	default:
		return be.addInstr("expr", typ, node.Kind.String(), nil)
	}
}

func (be *bodyEmitter) emitBinaryExpression(node *ast.Node, typ string) string {
	bin := node.AsBinaryExpression()
	leftId := be.emitExpression(bin.Left)
	rightId := be.emitExpression(bin.Right)
	op := bin.OperatorToken.Kind.String()
	if bin.OperatorToken.Kind == ast.KindEqualsToken {
		return be.addInstr("store", typ, nil, []string{leftId, rightId})
	}
	if bin.OperatorToken.Kind == ast.KindPlusEqualsToken || bin.OperatorToken.Kind == ast.KindMinusEqualsToken {
		return be.addInstr("compound", typ, op, []string{leftId, rightId})
	}
	return be.addInstr("binary", typ, op, []string{leftId, rightId})
}

func (be *bodyEmitter) emitUnaryExpression(node *ast.Node, typ string) string {
	unary := node.AsPrefixUnaryExpression()
	operandId := be.emitExpression(unary.Operand)
	if unary.Operator == ast.KindPlusPlusToken {
		return be.addInstr("inc", typ, nil, []string{operandId})
	}
	if unary.Operator == ast.KindMinusMinusToken {
		return be.addInstr("dec", typ, nil, []string{operandId})
	}
	return be.addInstr("unary", typ, unary.Operator.String(), []string{operandId})
}

func (be *bodyEmitter) emitPostfixExpression(node *ast.Node, typ string) string {
	postfix := node.AsPostfixUnaryExpression()
	operandId := be.emitExpression(postfix.Operand)
	if postfix.Operator == ast.KindPlusPlusToken {
		return be.addInstr("inc", typ, nil, []string{operandId})
	}
	if postfix.Operator == ast.KindMinusMinusToken {
		return be.addInstr("dec", typ, nil, []string{operandId})
	}
	return be.addInstr("unary", typ, postfix.Operator.String(), []string{operandId})
}

func (be *bodyEmitter) emitCallExpression(node *ast.Node, typ string) string {
	call := node.AsCallExpression()
	calleeId := be.emitExpression(call.Expression)
	var argIds []string
	for _, arg := range call.Arguments.Nodes {
		argIds = append(argIds, be.emitExpression(arg))
	}
	allOperands := append([]string{calleeId}, argIds...)
	return be.addInstr("call", typ, nil, allOperands)
}

func (be *bodyEmitter) emitPropertyAccess(node *ast.Node, typ string) string {
	pa := node.AsPropertyAccessExpression()
	objId := be.emitExpression(pa.Expression)
	propName := ""
	name := ast.GetElementOrPropertyAccessName(node)
	if name != nil {
		propName = name.Text()
	}
	return be.addInstr("prop", typ, propName, []string{objId})
}

func (be *bodyEmitter) emitElementAccess(node *ast.Node, typ string) string {
	ea := node.AsElementAccessExpression()
	objId := be.emitExpression(ea.Expression)
	idxId := be.emitExpression(ea.ArgumentExpression)
	return be.addInstr("elem", typ, nil, []string{objId, idxId})
}

func (be *bodyEmitter) emitNewExpression(node *ast.Node, typ string) string {
	ne := node.AsNewExpression()
	calleeId := be.emitExpression(ne.Expression)
	var argIds []string
	if ne.Arguments != nil {
		for _, arg := range ne.Arguments.Nodes {
			argIds = append(argIds, be.emitExpression(arg))
		}
	}
	allOperands := append([]string{calleeId}, argIds...)
	return be.addInstr("new", typ, nil, allOperands)
}

func (be *bodyEmitter) emitArrayLiteral(node *ast.Node, typ string) string {
	arr := node.AsArrayLiteralExpression()
	var elemIds []string
	if arr.Elements != nil {
		for _, elem := range arr.Elements.Nodes {
			elemIds = append(elemIds, be.emitExpression(elem))
		}
	}
	return be.addInstr("array", typ, nil, elemIds)
}

func (be *bodyEmitter) emitObjectLiteral(node *ast.Node, typ string) string {
	objId := be.addInstr("object", typ, nil, nil)
	ole := node.AsObjectLiteralExpression()
	if ole.Properties != nil {
		for _, propNode := range ole.Properties.Nodes {
			if propNode.Kind == ast.KindPropertyAssignment {
				pa := propNode.AsPropertyAssignment()
				propName := pa.Name().Text()
				valId := be.emitExpression(pa.Initializer)
				
				valTy := be.e.getNodeType(pa.Initializer)
				propId := be.addInstr("prop", valTy, propName, []string{objId})
				be.addInstr("store", valTy, nil, []string{propId, valId})
			} else if propNode.Kind == ast.KindMethodDeclaration {
				sym := propNode.Symbol()
				if sym != nil {
					be.e.emitMethod(sym)
					propName := sym.Name
					valTy := be.e.getNodeType(propNode)
					propId := be.addInstr("prop", valTy, propName, []string{objId})
					methodSymId := be.e.getOrCreateSymbolId(sym)
					methodValId := be.addInstr("ident", valTy, methodSymId, nil)
					be.addInstr("store", valTy, nil, []string{propId, methodValId})
				}
			}
		}
	}
	return objId
}

func (be *bodyEmitter) emitTemplateExpression(node *ast.Node, typ string) string {
	tmpl := node.AsTemplateExpression()
	var parts []string
	parts = append(parts, be.emitExpression(tmpl.Head))
	for _, span := range tmpl.TemplateSpans.Nodes {
		parts = append(parts, be.emitExpression(span.AsTemplateSpan().Expression))
		parts = append(parts, be.emitExpression(span.AsTemplateSpan().Literal))
	}
	return be.addInstr("template", typ, nil, parts)
}

func (be *bodyEmitter) emitConditionalExpression(node *ast.Node, typ string) string {
	cond := node.AsConditionalExpression()
	condId := be.emitExpression(cond.Condition)
	whenTrueId := be.emitExpression(cond.WhenTrue)
	whenFalseId := be.emitExpression(cond.WhenFalse)
	return be.addInstr("select", typ, nil, []string{condId, whenTrueId, whenFalseId})
}

func (e *Emitter) getNodeType(node *ast.Node) string {
	if t := e.checker.GetTypeOfNode(node); t != nil {
		return e.getOrCreateTypeId(t)
	}
	return ""
}

func (e *Emitter) emitVariableFromSymbol(sym *ast.Symbol) {
	v := &Variable{
		Name:   sym.Name,
		Symbol: e.getOrCreateSymbolId(sym),
	}

	if t := e.checker.GetTypeOfSymbol(sym); t != nil {
		v.Type = e.getOrCreateTypeId(t)
	}

	if len(sym.Declarations) > 0 {
		v.IsConst = ast.HasSyntacticModifier(sym.Declarations[0], ast.ModifierFlagsConst)
	}

	e.irProgram.Variables = append(e.irProgram.Variables, v)
}

func (e *Emitter) emitNodeTree(sf *ast.SourceFile) []*ASTNode {
	var nodes []*ASTNode
	for _, stmt := range sf.Statements.Nodes {
		if n := e.emitASTNode(stmt); n != nil {
			nodes = append(nodes, n)
		}
	}
	return nodes
}

func (e *Emitter) emitASTNode(node *ast.Node) *ASTNode {
	if node == nil {
		return nil
	}

	n := &ASTNode{
		Kind: node.Kind.String(),
		Pos:  node.Loc.Pos(),
		End:  node.Loc.End(),
	}

	if t := e.checker.GetTypeOfNode(node); t != nil {
		n.Type = e.getOrCreateTypeId(t)
	}
	if s := e.checker.GetResolvedSymbolOfNode(node); s != nil {
		n.Symbol = e.getOrCreateSymbolId(s)
	}

	node.ForEachChild(func(child *ast.Node) bool {
		if c := e.emitASTNode(child); c != nil {
			n.Children = append(n.Children, c)
		}
		return false
	})

	return n
}

func (e *Emitter) getOrCreateTypeId(t *checker.Type) string {
	if t == nil {
		return ""
	}

	identity := checker.GetRecursionIdentity(t)

	if state, ok := e.typeState[identity]; ok {
		return state.id
	}

	e.typeIdGen++
	id := fmt.Sprintf("t%d", e.typeIdGen)

	state := &typeState{id: id, resolving: true}
	e.typeState[identity] = state

	irType := &Type{
		Id:      id,
		Flags:   uint32(t.Flags()),
	}

	flags := t.Flags()

	switch {
	case e.checker.IsArrayType(t):
		// 任务 2: 数组一等引用类型。
		// 数组在 WIR 中以 kind="array" + elementFlags=[元素类型id] 表示，
		// 后端据此将其 codegen 为单一指针（ARC 管理的堆对象）。
		irType.Kind = "array"
		if elemTy := e.checker.GetElementTypeOfArrayType(t); elemTy != nil {
			irType.ElementFlags = append(irType.ElementFlags, e.getOrCreateTypeId(elemTy))
		}
	case flags&checker.TypeFlagsAny != 0:
		irType.Kind = "any"
	case flags&checker.TypeFlagsUnknown != 0:
		irType.Kind = "unknown"
	case flags&checker.TypeFlagsString != 0:
		irType.Kind = "string"
	case flags&checker.TypeFlagsNumber != 0:
		irType.Kind = "number"
	case flags&checker.TypeFlagsBoolean != 0:
		irType.Kind = "boolean"
	case flags&checker.TypeFlagsBigInt != 0:
		irType.Kind = "bigint"
	case flags&checker.TypeFlagsVoid != 0:
		irType.Kind = "void"
	case flags&checker.TypeFlagsUndefined != 0:
		irType.Kind = "undefined"
	case flags&checker.TypeFlagsNull != 0:
		irType.Kind = "null"
	case flags&checker.TypeFlagsNever != 0:
		irType.Kind = "never"
	case flags&checker.TypeFlagsUnion != 0:
		irType.Kind = "union"
		e.emitUnionOrIntersection(t, irType)
	case flags&checker.TypeFlagsIntersection != 0:
		irType.Kind = "intersection"
		e.emitUnionOrIntersection(t, irType)
	case flags&checker.TypeFlagsTypeParameter != 0:
		irType.Kind = "typeParameter"
		e.emitTypeParameter(t, irType)
	case flags&checker.TypeFlagsStringLiteral != 0:
		irType.Kind = "string"
		e.emitLiteralType(t, irType)
	case flags&checker.TypeFlagsLiteral != 0:
		irType.Kind = "literal"
		e.emitLiteralType(t, irType)
	case flags&checker.TypeFlagsNumberLiteral != 0:
		irType.Kind = "numberLiteral"
		e.emitLiteralType(t, irType)
	case flags&checker.TypeFlagsBooleanLiteral != 0:
		irType.Kind = "booleanLiteral"
		e.emitLiteralType(t, irType)
	case flags&checker.TypeFlagsObject != 0:
		irType.Kind = "object"
		e.emitObjectType(t, irType)
	case flags&checker.TypeFlagsIndex != 0:
		irType.Kind = "index"
		e.emitIndexType(t, irType)
	case flags&checker.TypeFlagsIndexedAccess != 0:
		irType.Kind = "indexedAccess"
		e.emitIndexedAccessType(t, irType)
	case flags&checker.TypeFlagsConditional != 0:
		irType.Kind = "conditional"
		e.emitConditionalType(t, irType)
	case flags&checker.TypeFlagsTemplateLiteral != 0:
		irType.Kind = "templateLiteral"
		e.emitTemplateLiteralType(t, irType)
	case flags&checker.TypeFlagsStringMapping != 0:
		irType.Kind = "stringMapping"
		e.emitStringMappingType(t, irType)
	case flags&checker.TypeFlagsSubstitution != 0:
		irType.Kind = "substitution"
		e.emitSubstitutionType(t, irType)
	default:
		irType.Kind = "unknown"
	}

	if t.Symbol() != nil {
		irType.Name = t.Symbol().Name
	}

	e.irProgram.Types = append(e.irProgram.Types, irType)
	state.resolving = false
	return id
}

func (e *Emitter) emitUnionOrIntersection(t *checker.Type, irType *Type) {
	types := t.Types()
	for _, constituent := range types {
		irType.Types = append(irType.Types, e.getOrCreateTypeId(constituent))
	}
}

func (e *Emitter) emitTypeParameter(t *checker.Type, irType *Type) {
	if constraint := e.checker.GetConstraintOfTypeParameter(t); constraint != nil {
		irType.Constraint = e.getOrCreateTypeId(constraint)
	}
	if defaultType := e.checker.GetDefaultFromTypeParameter(t); defaultType != nil {
		irType.Default = e.getOrCreateTypeId(defaultType)
	}
	tp := t.AsTypeParameter()
	irType.IsThisType = tp.IsThisType()
}

func (e *Emitter) emitLiteralType(t *checker.Type, irType *Type) {
	lt := t.AsLiteralType()
	irType.Value = lt.Value()
}

func (e *Emitter) emitObjectType(t *checker.Type, irType *Type) {
	irType.ObjectFlags = uint32(t.ObjectFlags())

	objectFlags := t.ObjectFlags()
	if objectFlags&checker.ObjectFlagsReference != 0 {
		if typeArgs := e.checker.GetTypeArguments(t); len(typeArgs) > 0 {
			for _, ta := range typeArgs {
				irType.TypeArgs = append(irType.TypeArgs, e.getOrCreateTypeId(ta))
			}
		}
	}

	if objectFlags&checker.ObjectFlagsClassOrInterface != 0 {
		if baseTypes := e.checker.GetBaseTypes(t); len(baseTypes) > 0 {
			for _, bt := range baseTypes {
				irType.BaseTypes = append(irType.BaseTypes, e.getOrCreateTypeId(bt))
			}
		}
	}

	if props := e.checker.GetPropertiesOfType(t); len(props) > 0 {
		for _, prop := range props {
			irProp := &Property{
				Name:       prop.Name,
				Symbol:     e.getOrCreateSymbolId(prop),
				IsOptional: prop.Flags&ast.SymbolFlagsOptional != 0,
			}
			if cachedType, ok := e.symTypeCache[prop]; ok {
				irProp.Type = cachedType
			} else if propType := e.checker.GetTypeOfSymbol(prop); propType != nil {
				propTypeId := e.getOrCreateTypeId(propType)
				e.symTypeCache[prop] = propTypeId
				irProp.Type = propTypeId
			}
			irType.Properties = append(irType.Properties, irProp)
		}
	}

	if callSigs := e.checker.GetCallSignatures(t); len(callSigs) > 0 {
		for _, sig := range callSigs {
			irType.Signatures = append(irType.Signatures, &SigRef{Id: e.getOrCreateSignatureId(sig)})
		}
	}

	if constructSigs := e.checker.GetConstructSignatures(t); len(constructSigs) > 0 {
		for _, sig := range constructSigs {
			irType.Signatures = append(irType.Signatures, &SigRef{Id: e.getOrCreateSignatureId(sig)})
		}
	}

	if indexInfos := e.checker.GetIndexInfosOfType(t); len(indexInfos) > 0 {
		for _, info := range indexInfos {
			irType.IndexInfos = append(irType.IndexInfos, &IndexInfo{
				KeyType:    e.getOrCreateTypeId(info.KeyType()),
				ValueType:  e.getOrCreateTypeId(info.ValueType()),
				IsReadonly: info.IsReadonly(),
			})
		}
	}

	if t.IsTupleType() {
		tt := t.TargetTupleType()
		if tt != nil {
			irType.IsReadonly = tt.IsReadonly()
			for _, info := range tt.ElementInfos() {
				irType.ElementFlags = append(irType.ElementFlags, fmt.Sprintf("%d", info.TupleElementFlags()))
			}
		}
	}
}

func (e *Emitter) emitIndexType(t *checker.Type, irType *Type) {
	it := t.AsIndexType()
	irType.Target = e.getOrCreateTypeId(it.Target())
}

func (e *Emitter) emitIndexedAccessType(t *checker.Type, irType *Type) {
	ia := t.AsIndexedAccessType()
	irType.ObjectType = e.getOrCreateTypeId(ia.ObjectType())
	irType.IndexType = e.getOrCreateTypeId(ia.IndexType())
}

func (e *Emitter) emitConditionalType(t *checker.Type, irType *Type) {
	ct := t.AsConditionalType()
	irType.CheckType = e.getOrCreateTypeId(ct.CheckType())
	irType.ExtendsType = e.getOrCreateTypeId(ct.ExtendsType())
}

func (e *Emitter) emitTemplateLiteralType(t *checker.Type, irType *Type) {
	tl := t.AsTemplateLiteralType()
	irType.Texts = tl.Texts()
	for _, ty := range tl.Types() {
		irType.Types = append(irType.Types, e.getOrCreateTypeId(ty))
	}
}

func (e *Emitter) emitStringMappingType(t *checker.Type, irType *Type) {
	sm := t.AsStringMappingType()
	irType.Target = e.getOrCreateTypeId(sm.Target())
}

func (e *Emitter) emitSubstitutionType(t *checker.Type, irType *Type) {
	st := t.AsSubstitutionType()
	irType.Target = e.getOrCreateTypeId(st.BaseType())
	irType.Constraint = e.getOrCreateTypeId(st.SubstConstraint())
}

func (e *Emitter) getOrCreateSymbolId(sym *ast.Symbol) string {
	if sym == nil {
		return ""
	}
	if exp := e.checker.GetExportSymbolOfSymbol(sym); exp != nil && exp != sym {
		sym = exp
	}
	if sym.Flags&ast.SymbolFlagsAlias != 0 {
		if resolved := e.checker.SkipAlias(sym); resolved != nil && resolved != sym {
			sym = resolved
		}
	}
	if id, ok := e.symbolMap[sym]; ok {
		return id
	}

	e.symbolIdGen++
	id := fmt.Sprintf("s%d", e.symbolIdGen)
	e.symbolMap[sym] = id

	flags := uint64(sym.Flags)
	if sym.Flags&ast.SymbolFlagsBlockScopedVariable != 0 {
		isConst := false
		for _, decl := range sym.Declarations {
			if ast.IsVarConst(decl) {
				isConst = true
				break
			}
		}
		if !isConst {
			flags &^= uint64(ast.SymbolFlagsBlockScopedVariable)
		}
	}

	irSym := &Symbol{
		Id:         id,
		Name:       sym.Name,
		Flags:      flags,
		CheckFlags: uint32(sym.CheckFlags),
		Kind:       e.symbolKindToString(sym),
	}

	if len(sym.Declarations) > 0 {
		for _, decl := range sym.Declarations {
			irSym.Declarations = append(irSym.Declarations, fmt.Sprintf("%d:%d", decl.Pos(), decl.End()))
		}
	}

	if sym.Parent != nil {
		irSym.Parent = e.getOrCreateSymbolId(sym.Parent)
	}

	if sym.Members != nil {
		for _, m := range sym.Members {
			irSym.Members = append(irSym.Members, e.getOrCreateSymbolId(m))
		}
	}

	if sym.Exports != nil {
		for _, ex := range sym.Exports {
			irSym.Exports = append(irSym.Exports, e.getOrCreateSymbolId(ex))
		}
	}

	e.irProgram.Symbols = append(e.irProgram.Symbols, irSym)
	return id
}

func (e *Emitter) getOrCreateSignatureId(sig *checker.Signature) string {
	if sig == nil {
		return ""
	}

	if state, ok := e.sigState[sig]; ok {
		return state.id
	}

	e.sigIdGen++
	id := fmt.Sprintf("sig%d", e.sigIdGen)

	state := &typeState{id: id, resolving: true}
	e.sigState[sig] = state

	irSig := &Signature{
		Id:      id,
		Kind:    "call",
		HasRest: sig.HasRestParameter(),
	}

	for _, param := range sig.Parameters() {
		p := Param{Name: param.Name, Symbol: e.getOrCreateSymbolId(param)}
		if cachedType, ok := e.symTypeCache[param]; ok {
			p.Type = cachedType
		} else if t := e.checker.GetTypeOfSymbol(param); t != nil {
			typeId := e.getOrCreateTypeId(t)
			e.symTypeCache[param] = typeId
			p.Type = typeId
		}
		irSig.Parameters = append(irSig.Parameters, p)
	}

	if ret := e.checker.GetReturnTypeOfSignature(sig); ret != nil {
		irSig.ReturnType = e.getOrCreateTypeId(ret)
	}

	for _, tp := range sig.TypeParameters() {
		irSig.TypeParameters = append(irSig.TypeParameters, e.getOrCreateTypeId(tp))
	}

	if sig.Declaration() != nil {
		irSig.Declaration = fmt.Sprintf("%d:%d", sig.Declaration().Pos(), sig.Declaration().End())
	}

	e.irProgram.Signatures = append(e.irProgram.Signatures, irSig)
	state.resolving = false
	return id
}

func (e *Emitter) symbolKindToString(sym *ast.Symbol) string {
	flags := sym.Flags
	switch {
	case flags&ast.SymbolFlagsFunction != 0:
		return "function"
	case flags&ast.SymbolFlagsClass != 0:
		return "class"
	case flags&ast.SymbolFlagsInterface != 0:
		return "interface"
	case flags&ast.SymbolFlagsVariable != 0, flags&ast.SymbolFlagsBlockScopedVariable != 0:
		return "variable"
	case flags&ast.SymbolFlagsMethod != 0:
		return "method"
	case flags&ast.SymbolFlagsProperty != 0:
		return "property"
	case flags&ast.SymbolFlagsEnum != 0:
		return "enum"
	case flags&ast.SymbolFlagsEnumMember != 0:
		return "enumMember"
	case flags&ast.SymbolFlagsModule != 0:
		return "module"
	case flags&ast.SymbolFlagsTypeAlias != 0:
		return "typeAlias"
	case flags&ast.SymbolFlagsAlias != 0:
		return "alias"
	case flags&ast.SymbolFlagsConstructor != 0:
		return "constructor"
	case flags&ast.SymbolFlagsSignature != 0:
		return "signature"
	case flags&ast.SymbolFlagsTypeLiteral != 0:
		return "typeLiteral"
	case flags&ast.SymbolFlagsObjectLiteral != 0:
		return "objectLiteral"
	default:
		return "unknown"
	}
}

func (e *Emitter) MarshalJSON() ([]byte, error) {
	if e.options.Compact {
		return json.Marshal(e.irProgram)
	}
	return json.MarshalIndent(e.irProgram, "", "  ")
}

func (e *Emitter) Serialize() ([]byte, error) {
	return e.MarshalJSON()
}

func (e *Emitter) treeShake() {
	reachableTypes := make(map[string]bool)
	reachableSymbols := make(map[string]bool)
	reachableSignatures := make(map[string]bool)

	userFiles := make(map[string]bool)
	for _, sf := range e.program.GetSourceFiles() {
		if !sf.IsDeclarationFile {
			userFiles[sf.FileName()] = true
		}
	}

	libDeclPositions := make(map[string]bool)
	for _, sf := range e.program.GetSourceFiles() {
		if sf.IsDeclarationFile {
			textLen := len(string(sf.Text()))
			for i := 0; i <= textLen; i++ {
				libDeclPositions[fmt.Sprintf("%d", i)] = true
			}
		}
	}

	isLibSymbol := func(s *Symbol) bool {
		for _, decl := range s.Declarations {
			if libDeclPositions[decl] {
				return true
			}
		}
		return false
	}

	var traceType func(typeId string, depth int)
	var traceSymbol func(symbolId string, depth int)
	var traceSignature func(sigId string, depth int)

	maxDepth := 3

	traceType = func(typeId string, depth int) {
		if typeId == "" || reachableTypes[typeId] {
			return
		}
		reachableTypes[typeId] = true

		if depth >= maxDepth {
			return
		}

		for _, t := range e.irProgram.Types {
			if t.Id == typeId {
				for _, subType := range t.Types {
					traceType(subType, depth+1)
				}
				if t.Target != "" {
					traceType(t.Target, depth+1)
				}
				for _, ta := range t.TypeArgs {
					traceType(ta, depth+1)
				}
				if t.ObjectType != "" {
					traceType(t.ObjectType, depth+1)
				}
				if t.IndexType != "" {
					traceType(t.IndexType, depth+1)
				}
				break
			}
		}
	}

	traceSymbol = func(symbolId string, depth int) {
		if symbolId == "" || reachableSymbols[symbolId] {
			return
		}

		var sym *Symbol
		for _, s := range e.irProgram.Symbols {
			if s.Id == symbolId {
				sym = s
				break
			}
		}

		if sym != nil && isLibSymbol(sym) && depth > 0 {
			reachableSymbols[symbolId] = true
			return
		}

		reachableSymbols[symbolId] = true

		if depth >= maxDepth {
			return
		}

		if sym != nil {
			if sym.Parent != "" {
				traceSymbol(sym.Parent, depth+1)
			}
		}
	}

	traceSignature = func(sigId string, depth int) {
		if sigId == "" || reachableSignatures[sigId] {
			return
		}
		reachableSignatures[sigId] = true

		if depth >= maxDepth {
			return
		}

		for _, sig := range e.irProgram.Signatures {
			if sig.Id == sigId {
				for _, p := range sig.Parameters {
					if p.Type != "" {
						traceType(p.Type, depth+1)
					}
				}
				if sig.ReturnType != "" {
					traceType(sig.ReturnType, depth+1)
				}
				break
			}
		}
	}

	for _, file := range e.irProgram.Files {
		if !strings.HasSuffix(file.Path, ".d.ts") {
			var traceNode func(node *ASTNode)
			traceNode = func(node *ASTNode) {
				if node.Type != "" {
					traceType(node.Type, 0)
				}
				if node.Symbol != "" {
					traceSymbol(node.Symbol, 0)
				}
				for _, child := range node.Children {
					traceNode(child)
				}
			}
			for _, node := range file.Nodes {
				traceNode(node)
			}
		}
	}

	for _, imp := range e.irProgram.Imports {
		traceSymbol(imp.Symbol, 0)
		for _, spec := range imp.Specifiers {
			traceSymbol(spec.Symbol, 0)
		}
	}

	for _, fn := range e.irProgram.Functions {
		traceSymbol(fn.Symbol, 0)
		for _, p := range fn.Parameters {
			if p.Type != "" {
				traceType(p.Type, 0)
			}
		}
		if fn.ReturnType != "" {
			traceType(fn.ReturnType, 0)
		}
		if fn.Signature != "" {
			traceSignature(fn.Signature, 0)
		}
	}

	for _, v := range e.irProgram.Variables {
		traceSymbol(v.Symbol, 0)
		if v.Type != "" {
			traceType(v.Type, 0)
		}
	}

	filteredTypes := make([]*Type, 0, len(reachableTypes))
	for _, t := range e.irProgram.Types {
		if reachableTypes[t.Id] {
			filteredTypes = append(filteredTypes, t)
		}
	}
	e.irProgram.Types = filteredTypes

	filteredSymbols := make([]*Symbol, 0, len(reachableSymbols))
	for _, s := range e.irProgram.Symbols {
		if reachableSymbols[s.Id] {
			filteredSymbols = append(filteredSymbols, s)
		}
	}
	e.irProgram.Symbols = filteredSymbols

	filteredSignatures := make([]*Signature, 0, len(reachableSignatures))
	for _, sig := range e.irProgram.Signatures {
		if reachableSignatures[sig.Id] {
			filteredSignatures = append(filteredSignatures, sig)
		}
	}
	e.irProgram.Signatures = filteredSignatures

	filteredGlobals := make([]*Global, 0)
	for _, g := range e.irProgram.Globals {
		if reachableSymbols[g.Symbol] {
			filteredGlobals = append(filteredGlobals, g)
		}
	}
	e.irProgram.Globals = filteredGlobals

	// Connect class constructor/references to class instance symbols to prevent incorrect tree-shaking
	for {
		changed := false
		for _, c := range e.irProgram.Classes {
			if !reachableSymbols[c.Symbol] {
				for _, s := range e.irProgram.Symbols {
					if reachableSymbols[s.Id] && s.Name == c.Name {
						reachableSymbols[c.Symbol] = true
						changed = true
						break
					}
				}
			}
		}
		if !changed {
			break
		}
	}

	filteredClasses := make([]*Class, 0)
	for _, c := range e.irProgram.Classes {
		if reachableSymbols[c.Symbol] {
			filteredClasses = append(filteredClasses, c)
		}
	}
	e.irProgram.Classes = filteredClasses

	filteredInterfaces := make([]*Interface, 0)
	for _, i := range e.irProgram.Interfaces {
		if reachableSymbols[i.Symbol] {
			filteredInterfaces = append(filteredInterfaces, i)
		}
	}
	e.irProgram.Interfaces = filteredInterfaces

	filteredEnums := make([]*Enum, 0)
	for _, en := range e.irProgram.Enums {
		if reachableSymbols[en.Symbol] {
			filteredEnums = append(filteredEnums, en)
		}
	}
	e.irProgram.Enums = filteredEnums

	filteredTypeAliases := make([]*TypeAlias, 0)
	for _, ta := range e.irProgram.TypeAliases {
		if reachableSymbols[ta.Symbol] {
			filteredTypeAliases = append(filteredTypeAliases, ta)
		}
	}
	e.irProgram.TypeAliases = filteredTypeAliases

	filteredNamespaces := make([]*Namespace, 0)
	for _, ns := range e.irProgram.Namespaces {
		if reachableSymbols[ns.Symbol] {
			filteredNamespaces = append(filteredNamespaces, ns)
		}
	}
	e.irProgram.Namespaces = filteredNamespaces
}
