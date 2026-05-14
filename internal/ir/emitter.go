package ir

import (
	"encoding/json"
	"fmt"

	"github.com/wolstn/wts/internal/ast"
	"github.com/wolstn/wts/internal/checker"
	"github.com/wolstn/wts/internal/compiler"
)

type Emitter struct {
	program    *compiler.Program
	checker    *checker.Checker
	irProgram  *Program
	typeMap    map[*checker.Type]string
	symbolMap  map[*ast.Symbol]string
	typeIdGen  int
	symbolIdGen int
}

func NewEmitter(program *compiler.Program) *Emitter {
	return &Emitter{
		program:    program,
		irProgram:  &Program{Version: Version},
		typeMap:    make(map[*checker.Type]string),
		symbolMap:  make(map[*ast.Symbol]string),
	}
}

func (e *Emitter) Emit() (*Program, error) {
	for _, sourceFile := range e.program.GetSourceFiles() {
		if sourceFile.IsDeclarationFile {
			continue
		}
		if err := e.emitFile(sourceFile); err != nil {
			return nil, err
		}
	}
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

	node.ForEachChild(func(child *ast.Node) bool {
		irNode.Children = append(irNode.Children, e.emitNode(child))
		return false
	})

	return irNode
}

func (e *Emitter) getOrCreateTypeId(t *checker.Type) string {
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
		irType.Target = e.getOrCreateSymbolId(sym)
	}

	e.irProgram.Types = append(e.irProgram.Types, irType)
	return id
}

func (e *Emitter) getOrCreateSymbolId(sym *ast.Symbol) string {
	if id, ok := e.symbolMap[sym]; ok {
		return id
	}

	e.symbolIdGen++
	id := fmt.Sprintf("s%d", e.symbolIdGen)
	e.symbolMap[sym] = id

	irSym := &Symbol{
		Id:    id,
		Name:  sym.Name,
		Flags: uint32(sym.Flags),
		Kind:  e.symbolKindToString(sym),
	}

	e.irProgram.Symbols = append(e.irProgram.Symbols, irSym)
	return id
}

func (e *Emitter) typeKindToString(t *checker.Type) string {
	flags := t.Flags()
	switch {
	case flags&checker.TypeFlagsAny != 0:
		return "any"
	case flags&checker.TypeFlagsString != 0:
		return "string"
	case flags&checker.TypeFlagsNumber != 0:
		return "number"
	case flags&checker.TypeFlagsBoolean != 0:
		return "boolean"
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
	default:
		return "unknown"
	}
}

func (e *Emitter) SetChecker(c *checker.Checker) {
	e.checker = c
}

func (p *Program) ToJSON() ([]byte, error) {
	return json.MarshalIndent(p, "", "  ")
}
