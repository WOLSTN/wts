package ir

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/wolstn/wts/internal/ast"
	"github.com/wolstn/wts/internal/checker"
	"github.com/wolstn/wts/internal/compiler"
)

const Version = 1

type Program struct {
	Version   int         `json:"version"`
	Files     []*File     `json:"files"`
	Types     []*Type     `json:"types"`
	Symbols   []*Symbol   `json:"symbols"`
	Globals   []*Global   `json:"globals"`
	Functions []*Function `json:"functions"`
}

type File struct {
	Path     string `json:"path"`
	Source   string `json:"source,omitempty"`
	RootNode *Node  `json:"rootNode"`
}

type Node struct {
	Kind     string  `json:"kind"`
	Pos      int     `json:"pos"`
	End      int     `json:"end"`
	Children []*Node `json:"children,omitempty"`
	Symbol   string  `json:"symbol,omitempty"`
	Type     string  `json:"type,omitempty"`
	Text     string  `json:"text,omitempty"`
}

type Type struct {
	Id          string   `json:"id"`
	Kind        string   `json:"kind"`
	Flags       uint32   `json:"flags"`
	Name        string   `json:"name,omitempty"`
	Members     []string `json:"members,omitempty"`
	Properties  []string `json:"properties,omitempty"`
	Signatures  []string `json:"signatures,omitempty"`
	TypeArgs    []string `json:"typeArgs,omitempty"`
	Target      string   `json:"target,omitempty"`
	Value       any      `json:"value,omitempty"`
}

type Symbol struct {
	Id           string   `json:"id"`
	Name         string   `json:"name"`
	Flags        uint64   `json:"flags"`
	Kind         string   `json:"kind"`
	Type         string   `json:"type,omitempty"`
	Declarations []string `json:"declarations,omitempty"`
	Members      []string `json:"members,omitempty"`
	Exports      []string `json:"exports,omitempty"`
}

type Global struct {
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
	Type   string `json:"type"`
}

type Function struct {
	Name       string  `json:"name"`
	Symbol     string  `json:"symbol"`
	Signature  string  `json:"signature"`
	Parameters []Param `json:"parameters"`
	ReturnType string  `json:"returnType"`
	Body       *Block  `json:"body,omitempty"`
}

type Param struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Block struct {
	Statements []Statement `json:"statements"`
}

type Statement struct {
	Kind string `json:"kind"`
	Data any    `json:"data"`
}

type Emitter struct {
	program     *compiler.Program
	checker     *checker.Checker
	irProgram   *Program
	typeMap     map[*checker.Type]string
	symbolMap   map[*ast.Symbol]string
	typeIdGen   int
	symbolIdGen int
}

func NewEmitter(program *compiler.Program) *Emitter {
	return &Emitter{
		program:   program,
		irProgram: &Program{Version: Version},
		typeMap:   make(map[*checker.Type]string),
		symbolMap: make(map[*ast.Symbol]string),
	}
}

func (e *Emitter) Emit() (*Program, error) {
	ctx := context.Background()
	checker, release := e.program.GetTypeChecker(ctx)
	defer release()

	e.checker = checker

	for _, sourceFile := range e.program.GetSourceFiles() {
		if sourceFile.IsDeclarationFile {
			continue
		}
		if err := e.emitFile(sourceFile); err != nil {
			return nil, err
		}
	}

	e.emitGlobalSymbols()

	return e.irProgram, nil
}

func (e *Emitter) emitFile(file *ast.SourceFile) error {
	irFile := &File{
		Path:   file.FileName(),
		Source: file.Text(),
	}

	irFile.RootNode = e.emitNode(file.AsNode())
	e.irProgram.Files = append(e.irProgram.Files, irFile)
	return nil
}

func (e *Emitter) emitNode(node *ast.Node) *Node {
	if node == nil {
		return nil
	}

	irNode := &Node{
		Kind: node.Kind.String(),
		Pos:  node.Pos(),
		End:  node.End(),
	}

	if e.checker != nil {
		if sym := e.checker.GetSymbolAtLocation(node); sym != nil {
			irNode.Symbol = e.getOrCreateSymbolId(sym)
		}
		if t := e.checker.GetTypeAtLocation(node); t != nil {
			irNode.Type = e.getOrCreateTypeId(t)
		}
	}

	if node.Kind == ast.KindIdentifier || node.Kind == ast.KindStringLiteral {
		if node.IdentifierText() != "" {
			irNode.Text = node.IdentifierText()
		} else if str := node.StringLiteralText(); str != "" {
			irNode.Text = str
		}
	}

	node.ForEachChild(func(child *ast.Node) bool {
		irNode.Children = append(irNode.Children, e.emitNode(child))
		return false
	})

	return irNode
}

func (e *Emitter) emitGlobalSymbols() {
	for _, sourceFile := range e.program.GetSourceFiles() {
		if sourceFile.IsDeclarationFile {
			continue
		}
		if sourceFile.Symbol != nil {
			for name, export := range sourceFile.Symbol.Exports {
				irGlobal := &Global{
					Name:   name,
					Symbol: e.getOrCreateSymbolId(export),
				}
				if e.checker != nil {
					if t := e.checker.GetTypeOfSymbol(export); t != nil {
						irGlobal.Type = e.getOrCreateTypeId(t)
					}
				}
				e.irProgram.Globals = append(e.irProgram.Globals, irGlobal)
			}
		}
	}
}

func (e *Emitter) getOrCreateTypeId(t *checker.Type) string {
	if t == nil {
		return ""
	}
	if id, ok := e.typeMap[t]; ok {
		return id
	}

	e.typeIdGen++
	id := fmt.Sprintf("t%d", e.typeIdGen)
	e.typeMap[t] = id

	irType := &Type{
		Id:    id,
		Flags: uint32(t.Flags()),
		Kind:  e.typeKindToString(t),
	}

	if sym := t.Symbol(); sym != nil {
		irType.Name = sym.Name
	}

	e.irProgram.Types = append(e.irProgram.Types, irType)
	return id
}

func (e *Emitter) getOrCreateSymbolId(sym *ast.Symbol) string {
	if sym == nil {
		return ""
	}
	if id, ok := e.symbolMap[sym]; ok {
		return id
	}

	e.symbolIdGen++
	id := fmt.Sprintf("s%d", e.symbolIdGen)
	e.symbolMap[sym] = id

	irSym := &Symbol{
		Id:    id,
		Name:  sym.Name,
		Flags: uint64(sym.Flags),
		Kind:  e.symbolKindToString(sym),
	}

	if e.checker != nil {
		if t := e.checker.GetTypeOfSymbol(sym); t != nil {
			irSym.Type = e.getOrCreateTypeId(t)
		}
	}

	e.irProgram.Symbols = append(e.irProgram.Symbols, irSym)
	return id
}

func (e *Emitter) typeKindToString(t *checker.Type) string {
	flags := t.Flags()
	switch {
	case flags&checker.TypeFlagsAny != 0:
		return "any"
	case flags&checker.TypeFlagsUnknown != 0:
		return "unknown"
	case flags&checker.TypeFlagsString != 0:
		return "string"
	case flags&checker.TypeFlagsNumber != 0:
		return "number"
	case flags&checker.TypeFlagsBoolean != 0:
		return "boolean"
	case flags&checker.TypeFlagsBigInt != 0:
		return "bigint"
	case flags&checker.TypeFlagsVoid != 0:
		return "void"
	case flags&checker.TypeFlagsUndefined != 0:
		return "undefined"
	case flags&checker.TypeFlagsNull != 0:
		return "null"
	case flags&checker.TypeFlagsNever != 0:
		return "never"
	case flags&checker.TypeFlagsObject != 0:
		return "object"
	case flags&checker.TypeFlagsUnion != 0:
		return "union"
	case flags&checker.TypeFlagsIntersection != 0:
		return "intersection"
	case flags&checker.TypeFlagsTypeParameter != 0:
		return "typeParameter"
	case flags&checker.TypeFlagsLiteral != 0:
		return "literal"
	case flags&checker.TypeFlagsStringLiteral != 0:
		return "stringLiteral"
	case flags&checker.TypeFlagsNumberLiteral != 0:
		return "numberLiteral"
	case flags&checker.TypeFlagsBooleanLiteral != 0:
		return "booleanLiteral"
	default:
		return "unknown"
	}
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
	case flags&ast.SymbolFlagsVariable != 0:
		return "variable"
	case flags&ast.SymbolFlagsConst != 0:
		return "const"
	case flags&ast.SymbolFlagsMethod != 0:
		return "method"
	case flags&ast.SymbolFlagsProperty != 0:
		return "property"
	case flags&ast.SymbolFlagsEnum != 0:
		return "enum"
	case flags&ast.SymbolFlagsModule != 0:
		return "module"
	case flags&ast.SymbolFlagsTypeAlias != 0:
		return "typeAlias"
	case flags&ast.SymbolFlagsAlias != 0:
		return "alias"
	default:
		return "unknown"
	}
}

func (p *Program) ToJSON() ([]byte, error) {
	return json.MarshalIndent(p, "", "  ")
}
