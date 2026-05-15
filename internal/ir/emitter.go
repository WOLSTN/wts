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
	Version     int          `json:"version"`
	Files       []*File      `json:"files"`
	Types       []*Type      `json:"types"`
	Symbols     []*Symbol    `json:"symbols"`
	Globals     []*Global    `json:"globals"`
	Functions   []*Function  `json:"functions"`
	Classes     []*Class     `json:"classes"`
	Interfaces  []*Interface `json:"interfaces"`
	Enums       []*Enum      `json:"enums,omitempty"`
	TypeAliases []*TypeAlias `json:"typeAliases,omitempty"`
	Namespaces  []*Namespace `json:"namespaces,omitempty"`
	Imports     []*Import    `json:"imports,omitempty"`
	Exports     []*Export    `json:"exports,omitempty"`
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
	Id             string           `json:"id"`
	Kind           string           `json:"kind"`
	Flags          uint32           `json:"flags"`
	Name           string           `json:"name,omitempty"`
	Members        []string         `json:"members,omitempty"`
	Properties     []*Property      `json:"properties,omitempty"`
	Signatures     []*TypeSignature `json:"signatures,omitempty"`
	TypeArgs       []string         `json:"typeArgs,omitempty"`
	TypeParams     []string         `json:"typeParams,omitempty"`
	Target         string           `json:"target,omitempty"`
	Value          any              `json:"value,omitempty"`
	ReturnType     string           `json:"returnType,omitempty"`
	Constraint     string           `json:"constraint,omitempty"`
	Default        string           `json:"default,omitempty"`
}

type Property struct {
	Name       string `json:"name"`
	Symbol     string `json:"symbol,omitempty"`
	Type       string `json:"type,omitempty"`
	IsOptional bool   `json:"isOptional,omitempty"`
	IsReadonly bool   `json:"isReadonly,omitempty"`
}

type TypeSignature struct {
	Kind           string   `json:"kind"`
	Parameters     []Param  `json:"parameters,omitempty"`
	ReturnType     string   `json:"returnType,omitempty"`
	TypeParameters []string `json:"typeParameters,omitempty"`
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

type Class struct {
	Name        string            `json:"name"`
	Symbol      string            `json:"symbol"`
	TypeParams  []string          `json:"typeParams,omitempty"`
	BaseClass   string            `json:"baseClass,omitempty"`
	Implements  []string          `json:"implements,omitempty"`
	Constructor *ClassMethod      `json:"constructor,omitempty"`
	Properties  []*ClassProperty  `json:"properties,omitempty"`
	Methods     []*ClassMethod    `json:"methods,omitempty"`
	IsAbstract  bool              `json:"isAbstract,omitempty"`
}

type ClassProperty struct {
	Name       string     `json:"name"`
	Symbol     string     `json:"symbol"`
	Type       string     `json:"type"`
	Visibility string     `json:"visibility,omitempty"`
	IsStatic   bool       `json:"isStatic,omitempty"`
	IsReadonly bool       `json:"isReadonly,omitempty"`
	IsOptional bool       `json:"isOptional,omitempty"`
	Init       *Expression `json:"init,omitempty"`
}

type ClassMethod struct {
	Name       string  `json:"name"`
	Symbol     string  `json:"symbol"`
	Parameters []Param `json:"parameters,omitempty"`
	ReturnType string  `json:"returnType,omitempty"`
	Visibility string  `json:"visibility,omitempty"`
	IsStatic   bool    `json:"isStatic,omitempty"`
	IsAbstract bool    `json:"isAbstract,omitempty"`
	Body       *Block  `json:"body,omitempty"`
}

type Interface struct {
	Name           string            `json:"name"`
	Symbol         string            `json:"symbol"`
	TypeParams     []string          `json:"typeParams,omitempty"`
	Extends        []string          `json:"extends,omitempty"`
	Properties     []*Property       `json:"properties,omitempty"`
	Methods        []*InterfaceMethod `json:"methods,omitempty"`
	CallSignatures []*TypeSignature  `json:"callSignatures,omitempty"`
}

type InterfaceMethod struct {
	Name       string `json:"name"`
	Symbol     string `json:"symbol"`
	Parameters []Param `json:"parameters,omitempty"`
	ReturnType string  `json:"returnType,omitempty"`
}

type Import struct {
	Kind       string        `json:"kind"`
	ModulePath string        `json:"modulePath"`
	IsTypeOnly bool          `json:"isTypeOnly,omitempty"`
	IsDefault  bool          `json:"isDefault,omitempty"`
	Namespace  string        `json:"namespace,omitempty"`
	Specifiers []*ImportSpec `json:"specifiers,omitempty"`
}

type ImportSpec struct {
	Name     string `json:"name"`
	Property string `json:"property,omitempty"`
	Symbol   string `json:"symbol,omitempty"`
	IsType   bool   `json:"isType,omitempty"`
}

type Export struct {
	Kind        string       `json:"kind"`
	Name        string       `json:"name,omitempty"`
	Symbol      string       `json:"symbol,omitempty"`
	IsTypeOnly  bool         `json:"isTypeOnly,omitempty"`
	IsDefault   bool         `json:"isDefault,omitempty"`
	ModulePath  string       `json:"modulePath,omitempty"`
	Declaration *Declaration `json:"declaration,omitempty"`
}

type Enum struct {
	Name    string        `json:"name"`
	Symbol  string        `json:"symbol"`
	Members []*EnumMember `json:"members"`
	IsConst bool          `json:"isConst,omitempty"`
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
	Target     string   `json:"target"`
}

type Namespace struct {
	Name    string             `json:"name"`
	Symbol  string             `json:"symbol"`
	Members []*NamespaceMember `json:"members,omitempty"`
}

type NamespaceMember struct {
	Kind   string `json:"kind"`
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

type Param struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Block struct {
	Statements []Statement `json:"statements"`
}

type Statement struct {
	Kind        string       `json:"kind"`
	Expression  *Expression  `json:"expression,omitempty"`
	Declaration *Declaration `json:"declaration,omitempty"`
	Return      *ReturnStmt  `json:"return,omitempty"`
	If          *IfStmt      `json:"if,omitempty"`
	While       *WhileStmt   `json:"while,omitempty"`
	For         *ForStmt     `json:"for,omitempty"`
	Block       *Block       `json:"block,omitempty"`
}

type Expression struct {
	Kind       string        `json:"kind"`
	Text       string        `json:"text,omitempty"`
	Type       string        `json:"type,omitempty"`
	Symbol     string        `json:"symbol,omitempty"`
	Value      any           `json:"value,omitempty"`
	Operator   string        `json:"operator,omitempty"`
	Left       *Expression   `json:"left,omitempty"`
	Right      *Expression   `json:"right,omitempty"`
	Operand    *Expression   `json:"operand,omitempty"`
	Callee     *Expression   `json:"callee,omitempty"`
	Arguments  []Expression  `json:"arguments,omitempty"`
	Object     *Expression   `json:"object,omitempty"`
	Property   *Expression   `json:"property,omitempty"`
	Elements   []Expression  `json:"elements,omitempty"`
	Properties []PropertyInit `json:"properties,omitempty"`
	Condition  *Expression   `json:"condition,omitempty"`
	Consequent *Expression   `json:"consequent,omitempty"`
	Alternate  *Expression   `json:"alternate,omitempty"`
	Parameters []Param       `json:"parameters,omitempty"`
	Body       *Block        `json:"body,omitempty"`
}

type PropertyInit struct {
	Name  string     `json:"name"`
	Value *Expression `json:"value,omitempty"`
}

type Declaration struct {
	Name       string     `json:"name"`
	Symbol     string     `json:"symbol,omitempty"`
	Type       string     `json:"type,omitempty"`
	Init       *Expression `json:"init,omitempty"`
	Params     []Param    `json:"params,omitempty"`
	ReturnType string     `json:"returnType,omitempty"`
	Body       *Block     `json:"body,omitempty"`
}

type ReturnStmt struct {
	Value *Expression `json:"value,omitempty"`
}

type IfStmt struct {
	Condition *Expression `json:"condition"`
	Then      *Block      `json:"then"`
	Else      *Block      `json:"else,omitempty"`
}

type WhileStmt struct {
	Condition *Expression `json:"condition"`
	Body      *Block      `json:"body"`
}

type ForStmt struct {
	Init      *Statement  `json:"init,omitempty"`
	Condition *Expression `json:"condition,omitempty"`
	Update    *Statement  `json:"update,omitempty"`
	Body      *Block      `json:"body"`
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
		e.emitSourceFileSymbols(sourceFile)
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

func (e *Emitter) emitSourceFileSymbols(file *ast.SourceFile) {
	if file.Symbol == nil {
		return
	}

	for name, sym := range file.Symbol.Exports {
		if sym == nil {
			continue
		}
		
		e.emitGlobal(name, sym)
		e.emitSymbolDeclaration(sym)
	}
}

func (e *Emitter) emitGlobal(name string, sym *ast.Symbol) {
	irGlobal := &Global{
		Name:   name,
		Symbol: e.getOrCreateSymbolId(sym),
	}

	if e.checker != nil {
		if t := e.checker.GetTypeOfSymbol(sym); t != nil {
			irGlobal.Type = e.getOrCreateTypeId(t)
		}
	}

	e.irProgram.Globals = append(e.irProgram.Globals, irGlobal)
}

func (e *Emitter) emitSymbolDeclaration(sym *ast.Symbol) {
	if sym == nil || len(sym.Declarations) == 0 {
		return
	}

	for _, decl := range sym.Declarations {
		if decl == nil {
			continue
		}

		switch decl.Kind {
		case ast.KindFunctionDeclaration:
			e.emitFunctionFromSymbol(sym, decl)
		case ast.KindClassDeclaration:
			e.emitClassFromSymbol(sym, decl)
		case ast.KindInterfaceDeclaration:
			e.emitInterfaceFromSymbol(sym, decl)
		case ast.KindEnumDeclaration:
			e.emitEnumFromSymbol(sym, decl)
		case ast.KindTypeAliasDeclaration:
			e.emitTypeAliasFromSymbol(sym, decl)
		case ast.KindModuleDeclaration:
			e.emitNamespaceFromSymbol(sym, decl)
		case ast.KindImportDeclaration, ast.KindImportEqualsDeclaration:
			e.emitImportFromDeclaration(decl)
		case ast.KindExportDeclaration, ast.KindExportAssignment:
			e.emitExportFromDeclaration(decl)
		}
	}
}

func (e *Emitter) emitFunctionFromSymbol(sym *ast.Symbol, decl *ast.Node) {
	irFunc := &Function{
		Name:   sym.Name,
		Symbol: e.getOrCreateSymbolId(sym),
	}

	if e.checker != nil {
		sig := e.checker.GetSignatureFromDeclaration(decl)
		if sig != nil {
			for _, param := range sig.Parameters() {
				irParam := Param{Name: param.Name}
				if paramType := e.checker.GetTypeOfSymbol(param); paramType != nil {
					irParam.Type = e.getOrCreateTypeId(paramType)
				}
				irFunc.Parameters = append(irFunc.Parameters, irParam)
			}

			if returnType := e.checker.GetReturnTypeOfSignature(sig); returnType != nil {
				irFunc.ReturnType = e.getOrCreateTypeId(returnType)
			}
		}
	}

	funcDecl := decl.AsFunctionDeclaration()
	if funcDecl != nil && funcDecl.Body != nil {
		irFunc.Body = e.emitBlock(funcDecl.Body)
	}

	e.irProgram.Functions = append(e.irProgram.Functions, irFunc)
}

func (e *Emitter) emitClassFromSymbol(sym *ast.Symbol, decl *ast.Node) {
	irClass := &Class{
		Name:   sym.Name,
		Symbol: e.getOrCreateSymbolId(sym),
	}

	classDecl := decl.AsClassDeclaration()
	if classDecl == nil {
		e.irProgram.Classes = append(e.irProgram.Classes, irClass)
		return
	}

	if typeParams := decl.TypeParameters(); len(typeParams) > 0 {
		for _, tp := range typeParams {
			tpSym := e.checker.GetSymbolAtLocation(tp.Name())
			if tpSym != nil {
				irClass.TypeParams = append(irClass.TypeParams, e.getOrCreateSymbolId(tpSym))
			}
		}
	}

	if heritageClauses := classDecl.HeritageClauses; heritageClauses != nil {
		for _, clause := range heritageClauses.Nodes {
			if clause.Kind == ast.KindHeritageClause {
				hc := clause.AsHeritageClause()
				if hc.Token == ast.KindExtendsKeyword {
					for _, typeExprNode := range hc.Types.Nodes {
						typeExpr := typeExprNode.AsExpressionWithTypeArguments()
						if typeExpr != nil && typeExpr.Expression != nil {
							if typeSym := e.checker.GetSymbolAtLocation(typeExpr.Expression); typeSym != nil {
								irClass.BaseClass = e.getOrCreateSymbolId(typeSym)
							}
						}
					}
				} else if hc.Token == ast.KindImplementsKeyword {
					for _, typeExprNode := range hc.Types.Nodes {
						typeExpr := typeExprNode.AsExpressionWithTypeArguments()
						if typeExpr != nil && typeExpr.Expression != nil {
							if typeSym := e.checker.GetSymbolAtLocation(typeExpr.Expression); typeSym != nil {
								irClass.Implements = append(irClass.Implements, e.getOrCreateSymbolId(typeSym))
							}
						}
					}
				}
			}
		}
	}

	if modifiers := decl.Modifiers(); modifiers != nil {
		for _, mod := range modifiers.Nodes {
			if mod.Kind == ast.KindAbstractKeyword {
				irClass.IsAbstract = true
			}
		}
	}

	for _, member := range classDecl.Members.Nodes {
		switch member.Kind {
		case ast.KindConstructor:
			irClass.Constructor = e.emitClassMethod(member)
		case ast.KindPropertyDeclaration:
			prop := e.emitClassProperty(member)
			if prop != nil {
				irClass.Properties = append(irClass.Properties, prop)
			}
		case ast.KindMethodDeclaration:
			method := e.emitClassMethod(member)
			if method != nil {
				irClass.Methods = append(irClass.Methods, method)
			}
		case ast.KindGetAccessor, ast.KindSetAccessor:
			method := e.emitClassMethod(member)
			if method != nil {
				irClass.Methods = append(irClass.Methods, method)
			}
		}
	}

	e.irProgram.Classes = append(e.irProgram.Classes, irClass)
}

func (e *Emitter) emitClassProperty(node *ast.Node) *ClassProperty {
	if node == nil {
		return nil
	}

	nameNode := node.Name()
	if nameNode == nil {
		return nil
	}

	prop := &ClassProperty{
		Name: nameNode.Text(),
	}

	if sym := e.checker.GetSymbolAtLocation(nameNode); sym != nil {
		prop.Symbol = e.getOrCreateSymbolId(sym)
	}

	if e.checker != nil {
		if t := e.checker.GetTypeAtLocation(node); t != nil {
			prop.Type = e.getOrCreateTypeId(t)
		}
	}

	if modifiers := node.Modifiers(); modifiers != nil {
		for _, mod := range modifiers.Nodes {
			switch mod.Kind {
			case ast.KindPublicKeyword:
				prop.Visibility = "public"
			case ast.KindPrivateKeyword:
				prop.Visibility = "private"
			case ast.KindProtectedKeyword:
				prop.Visibility = "protected"
			case ast.KindStaticKeyword:
				prop.IsStatic = true
			case ast.KindReadonlyKeyword:
				prop.IsReadonly = true
			}
		}
	}

	if sym := e.checker.GetSymbolAtLocation(nameNode); sym != nil {
		if sym.Flags&ast.SymbolFlagsOptional != 0 {
			prop.IsOptional = true
		}
	}

	propDecl := node.AsPropertyDeclaration()
	if propDecl != nil && propDecl.Initializer != nil {
		prop.Init = e.emitExpression(propDecl.Initializer)
	}

	return prop
}

func (e *Emitter) emitClassMethod(node *ast.Node) *ClassMethod {
	if node == nil {
		return nil
	}

	nameNode := node.Name()
	if nameNode == nil {
		return nil
	}

	method := &ClassMethod{
		Name: nameNode.Text(),
	}

	if sym := e.checker.GetSymbolAtLocation(nameNode); sym != nil {
		method.Symbol = e.getOrCreateSymbolId(sym)
	}

	if e.checker != nil {
		sig := e.checker.GetSignatureFromDeclaration(node)
		if sig != nil {
			for _, param := range sig.Parameters() {
				irParam := Param{Name: param.Name}
				if paramType := e.checker.GetTypeOfSymbol(param); paramType != nil {
					irParam.Type = e.getOrCreateTypeId(paramType)
				}
				method.Parameters = append(method.Parameters, irParam)
			}

			if returnType := e.checker.GetReturnTypeOfSignature(sig); returnType != nil {
				method.ReturnType = e.getOrCreateTypeId(returnType)
			}
		}
	}

	if modifiers := node.Modifiers(); modifiers != nil {
		for _, mod := range modifiers.Nodes {
			switch mod.Kind {
			case ast.KindPublicKeyword:
				method.Visibility = "public"
			case ast.KindPrivateKeyword:
				method.Visibility = "private"
			case ast.KindProtectedKeyword:
				method.Visibility = "protected"
			case ast.KindStaticKeyword:
				method.IsStatic = true
			case ast.KindAbstractKeyword:
				method.IsAbstract = true
			}
		}
	}

	switch node.Kind {
	case ast.KindConstructor:
		ctor := node.AsConstructorDeclaration()
		if ctor != nil && ctor.Body != nil {
			method.Body = e.emitBlock(ctor.Body)
		}
	case ast.KindMethodDeclaration:
		methodDecl := node.AsMethodDeclaration()
		if methodDecl != nil && methodDecl.Body != nil {
			method.Body = e.emitBlock(methodDecl.Body)
		}
	case ast.KindGetAccessor:
		getAccessor := node.AsGetAccessorDeclaration()
		if getAccessor != nil && getAccessor.Body != nil {
			method.Body = e.emitBlock(getAccessor.Body)
		}
	case ast.KindSetAccessor:
		setAccessor := node.AsSetAccessorDeclaration()
		if setAccessor != nil && setAccessor.Body != nil {
			method.Body = e.emitBlock(setAccessor.Body)
		}
	}

	return method
}

func (e *Emitter) emitInterfaceFromSymbol(sym *ast.Symbol, decl *ast.Node) {
	irInterface := &Interface{
		Name:   sym.Name,
		Symbol: e.getOrCreateSymbolId(sym),
	}

	if typeParams := decl.TypeParameters(); len(typeParams) > 0 {
		for _, tp := range typeParams {
			tpSym := e.checker.GetSymbolAtLocation(tp.Name())
			if tpSym != nil {
				irInterface.TypeParams = append(irInterface.TypeParams, e.getOrCreateSymbolId(tpSym))
			}
		}
	}

	ifaceDecl := decl.AsInterfaceDeclaration()
	if ifaceDecl != nil && ifaceDecl.HeritageClauses != nil {
		for _, clause := range ifaceDecl.HeritageClauses.Nodes {
			if clause.Kind == ast.KindHeritageClause {
				hc := clause.AsHeritageClause()
				if hc.Token == ast.KindExtendsKeyword {
					for _, typeExprNode := range hc.Types.Nodes {
						typeExpr := typeExprNode.AsExpressionWithTypeArguments()
						if typeExpr != nil && typeExpr.Expression != nil {
							if typeSym := e.checker.GetSymbolAtLocation(typeExpr.Expression); typeSym != nil {
								irInterface.Extends = append(irInterface.Extends, e.getOrCreateSymbolId(typeSym))
							}
						}
					}
				}
			}
		}
	}

	if sym.Members != nil {
		for memberName, memberSym := range sym.Members {
			if memberSym == nil || len(memberSym.Declarations) == 0 {
				continue
			}

			memberDecl := memberSym.Declarations[0]
			switch memberDecl.Kind {
			case ast.KindPropertySignature:
				prop := e.emitInterfaceProperty(memberSym, memberDecl)
				if prop != nil {
					irInterface.Properties = append(irInterface.Properties, prop)
				}
			case ast.KindMethodSignature:
				method := e.emitInterfaceMethod(memberSym, memberDecl)
				if method != nil {
					irInterface.Methods = append(irInterface.Methods, method)
				}
			case ast.KindCallSignature:
				sig := e.emitTypeSignature(memberDecl)
				if sig != nil {
					irInterface.CallSignatures = append(irInterface.CallSignatures, sig)
				}
			}
			_ = memberName
		}
	}

	e.irProgram.Interfaces = append(e.irProgram.Interfaces, irInterface)
}

func (e *Emitter) emitInterfaceProperty(sym *ast.Symbol, decl *ast.Node) *Property {
	prop := &Property{
		Name:   sym.Name,
		Symbol: e.getOrCreateSymbolId(sym),
	}

	if e.checker != nil {
		if t := e.checker.GetTypeOfSymbol(sym); t != nil {
			prop.Type = e.getOrCreateTypeId(t)
		}
	}

	if sym.Flags&ast.SymbolFlagsOptional != 0 {
		prop.IsOptional = true
	}

	if modifiers := decl.Modifiers(); modifiers != nil {
		for _, mod := range modifiers.Nodes {
			if mod.Kind == ast.KindReadonlyKeyword {
				prop.IsReadonly = true
			}
		}
	}

	return prop
}

func (e *Emitter) emitInterfaceMethod(sym *ast.Symbol, decl *ast.Node) *InterfaceMethod {
	method := &InterfaceMethod{
		Name:   sym.Name,
		Symbol: e.getOrCreateSymbolId(sym),
	}

	if e.checker != nil {
		sig := e.checker.GetSignatureFromDeclaration(decl)
		if sig != nil {
			for _, param := range sig.Parameters() {
				irParam := Param{Name: param.Name}
				if paramType := e.checker.GetTypeOfSymbol(param); paramType != nil {
					irParam.Type = e.getOrCreateTypeId(paramType)
				}
				method.Parameters = append(method.Parameters, irParam)
			}

			if returnType := e.checker.GetReturnTypeOfSignature(sig); returnType != nil {
				method.ReturnType = e.getOrCreateTypeId(returnType)
			}
		}
	}

	return method
}

func (e *Emitter) emitEnumFromSymbol(sym *ast.Symbol, decl *ast.Node) {
	irEnum := &Enum{
		Name:   sym.Name,
		Symbol: e.getOrCreateSymbolId(sym),
	}

	enumDecl := decl.AsEnumDeclaration()
	if enumDecl != nil {
		if modifiers := decl.Modifiers(); modifiers != nil {
			for _, mod := range modifiers.Nodes {
				if mod.Kind == ast.KindConstKeyword {
					irEnum.IsConst = true
					break
				}
			}
		}

		for _, member := range enumDecl.Members.Nodes {
			if member.Kind == ast.KindEnumMember {
				irMember := e.emitEnumMember(member)
				if irMember != nil {
					irEnum.Members = append(irEnum.Members, irMember)
				}
			}
		}
	}

	e.irProgram.Enums = append(e.irProgram.Enums, irEnum)
}

func (e *Emitter) emitEnumMember(node *ast.Node) *EnumMember {
	if node == nil {
		return nil
	}

	nameNode := node.Name()
	if nameNode == nil {
		return nil
	}

	sym := e.checker.GetSymbolAtLocation(nameNode)
	if sym == nil {
		return nil
	}

	member := &EnumMember{
		Name:   nameNode.Text(),
		Symbol: e.getOrCreateSymbolId(sym),
	}

	if t := e.checker.GetTypeOfSymbol(sym); t != nil {
		member.Type = e.getOrCreateTypeId(t)
	}

	enumMember := node.AsEnumMember()
	if enumMember != nil && enumMember.Initializer != nil {
		init := enumMember.Initializer
		switch init.Kind {
		case ast.KindNumericLiteral:
			member.Value = init.Text()
		case ast.KindStringLiteral:
			member.Value = init.Text()
		case ast.KindTrueKeyword:
			member.Value = true
		case ast.KindFalseKeyword:
			member.Value = false
		}
	}

	return member
}

func (e *Emitter) emitTypeAliasFromSymbol(sym *ast.Symbol, decl *ast.Node) {
	irTypeAlias := &TypeAlias{
		Name:   sym.Name,
		Symbol: e.getOrCreateSymbolId(sym),
	}

	if typeParams := decl.TypeParameters(); len(typeParams) > 0 {
		for _, tp := range typeParams {
			tpSym := e.checker.GetSymbolAtLocation(tp.Name())
			if tpSym != nil {
				irTypeAlias.TypeParams = append(irTypeAlias.TypeParams, e.getOrCreateSymbolId(tpSym))
			}
		}
	}

	typeAliasDecl := decl.AsTypeAliasDeclaration()
	if typeAliasDecl != nil && typeAliasDecl.Type != nil {
		if t := e.checker.GetTypeFromTypeNode(typeAliasDecl.Type); t != nil {
			irTypeAlias.Target = e.getOrCreateTypeId(t)
		}
	}

	e.irProgram.TypeAliases = append(e.irProgram.TypeAliases, irTypeAlias)
}

func (e *Emitter) emitNamespaceFromSymbol(sym *ast.Symbol, decl *ast.Node) {
	irNamespace := &Namespace{
		Name:   sym.Name,
		Symbol: e.getOrCreateSymbolId(sym),
	}

	moduleDecl := decl.AsModuleDeclaration()
	if moduleDecl != nil && moduleDecl.Body != nil {
		if block := moduleDecl.Body.AsModuleBlock(); block != nil {
			for _, stmt := range block.Statements.Nodes {
				irMember := e.emitNamespaceMember(stmt)
				if irMember != nil {
					irNamespace.Members = append(irNamespace.Members, irMember)
				}
			}
		}
	}

	e.irProgram.Namespaces = append(e.irProgram.Namespaces, irNamespace)
}

func (e *Emitter) emitNamespaceMember(node *ast.Node) *NamespaceMember {
	if node == nil {
		return nil
	}

	nameNode := node.Name()
	if nameNode == nil {
		return nil
	}

	sym := e.checker.GetSymbolAtLocation(nameNode)
	if sym == nil {
		return nil
	}

	return &NamespaceMember{
		Kind:   node.Kind.String(),
		Name:   sym.Name,
		Symbol: e.getOrCreateSymbolId(sym),
	}
}

func (e *Emitter) emitImportFromDeclaration(node *ast.Node) {
	if node == nil {
		return
	}

	irImport := &Import{
		Kind: node.Kind.String(),
	}

	switch node.Kind {
	case ast.KindImportDeclaration:
		decl := node.AsImportDeclaration()
		if decl == nil {
			return
		}

		if decl.ModuleSpecifier != nil {
			irImport.ModulePath = decl.ModuleSpecifier.Text()
		}

		if decl.ImportClause != nil {
			ic := decl.ImportClause.AsImportClause()
			if ic != nil {
				if name := ic.Name(); name != nil {
					irImport.IsDefault = true
					if sym := e.checker.GetSymbolAtLocation(name); sym != nil {
						irImport.Specifiers = append(irImport.Specifiers, &ImportSpec{
							Name:   name.Text(),
							Symbol: e.getOrCreateSymbolId(sym),
						})
					}
				}

				if ic.NamedBindings != nil {
					if ic.PhaseModifier == ast.KindTypeKeyword {
						irImport.IsTypeOnly = true
					}

					switch ic.NamedBindings.Kind {
					case ast.KindNamespaceImport:
						nsImport := ic.NamedBindings.AsNamespaceImport()
						if nsImport != nil {
							if name := nsImport.Name(); name != nil {
								if sym := e.checker.GetSymbolAtLocation(name); sym != nil {
									irImport.Namespace = sym.Name
								}
							}
						}
					case ast.KindNamedImports:
						for _, spec := range ic.NamedBindings.Elements() {
							if spec.Kind == ast.KindImportSpecifier {
								importSpec := spec.AsImportSpecifier()
								if importSpec != nil {
									if name := importSpec.Name(); name != nil {
										if sym := e.checker.GetSymbolAtLocation(name); sym != nil {
											irImport.Specifiers = append(irImport.Specifiers, &ImportSpec{
												Name:   name.Text(),
												Symbol: e.getOrCreateSymbolId(sym),
												IsType: importSpec.IsTypeOnly,
											})
										}
									}
								}
							}
						}
					}
				}
			}
		}

	case ast.KindImportEqualsDeclaration:
		decl := node.AsImportEqualsDeclaration()
		if decl == nil {
			return
		}

		if decl.ModuleReference != nil {
			irImport.ModulePath = decl.ModuleReference.Text()
		}

		if name := decl.Name(); name != nil {
			if sym := e.checker.GetSymbolAtLocation(name); sym != nil {
				irImport.Namespace = sym.Name
			}
		}
	}

	e.irProgram.Imports = append(e.irProgram.Imports, irImport)
}

func (e *Emitter) emitExportFromDeclaration(node *ast.Node) {
	if node == nil {
		return
	}

	switch node.Kind {
	case ast.KindExportDeclaration:
		e.emitExportDeclaration(node)
	case ast.KindExportAssignment:
		e.emitExportAssignment(node)
	default:
		nameNode := node.Name()
		if nameNode == nil {
			return
		}

		sym := e.checker.GetSymbolAtLocation(nameNode)
		if sym == nil {
			return
		}

		irExport := &Export{
			Kind:   node.Kind.String(),
			Name:   sym.Name,
			Symbol: e.getOrCreateSymbolId(sym),
		}

		if modifiers := node.Modifiers(); modifiers != nil {
			for _, mod := range modifiers.Nodes {
				if mod.Kind == ast.KindDefaultKeyword {
					irExport.IsDefault = true
				}
			}
		}

		e.irProgram.Exports = append(e.irProgram.Exports, irExport)
	}
}

func (e *Emitter) emitExportDeclaration(node *ast.Node) {
	if node == nil {
		return
	}

	decl := node.AsExportDeclaration()
	if decl == nil {
		return
	}

	irExport := &Export{
		Kind: node.Kind.String(),
	}

	if decl.IsTypeOnly {
		irExport.IsTypeOnly = true
	}

	if decl.ExportClause != nil {
		if decl.ExportClause.Kind == ast.KindNamedExports {
			for _, spec := range decl.ExportClause.Elements() {
				if spec.Kind == ast.KindExportSpecifier {
					name := spec.Name()
					if name != nil {
						exportName := name.Text()
						if sym := e.checker.GetSymbolAtLocation(name); sym != nil {
							e.irProgram.Exports = append(e.irProgram.Exports, &Export{
								Kind:   "export_specifier",
								Name:   exportName,
								Symbol: e.getOrCreateSymbolId(sym),
							})
						}
					}
				}
			}
		}
	}

	if decl.ModuleSpecifier != nil {
		irExport.ModulePath = decl.ModuleSpecifier.Text()
	}

	if irExport.Name != "" || irExport.ModulePath != "" {
		e.irProgram.Exports = append(e.irProgram.Exports, irExport)
	}
}

func (e *Emitter) emitExportAssignment(node *ast.Node) {
	if node == nil {
		return
	}

	decl := node.AsExportAssignment()
	if decl == nil {
		return
	}

	irExport := &Export{
		Kind:      node.Kind.String(),
		IsDefault: !decl.IsExportEquals,
	}

	if decl.Expression != nil {
		if sym := e.checker.GetSymbolAtLocation(decl.Expression); sym != nil {
			irExport.Name = sym.Name
			irExport.Symbol = e.getOrCreateSymbolId(sym)
		}
	}

	e.irProgram.Exports = append(e.irProgram.Exports, irExport)
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

	skipTypeCheck := node.Kind == ast.KindImportDeclaration ||
		node.Kind == ast.KindImportClause ||
		node.Kind == ast.KindImportSpecifier ||
		node.Kind == ast.KindImportEqualsDeclaration ||
		node.Kind == ast.KindExportDeclaration ||
		node.Kind == ast.KindExportAssignment ||
		node.Kind == ast.KindExportSpecifier ||
		node.Kind == ast.KindNamespaceImport ||
		node.Kind == ast.KindNamedImports

	if e.checker != nil && !skipTypeCheck {
		if sym := e.checker.GetSymbolAtLocation(node); sym != nil {
			irNode.Symbol = e.getOrCreateSymbolId(sym)
		}
		if t := e.checker.GetTypeAtLocation(node); t != nil {
			irNode.Type = e.getOrCreateTypeId(t)
		}
	}

	if node.Kind == ast.KindIdentifier || node.Kind == ast.KindStringLiteral {
		irNode.Text = node.Text()
	}

	node.ForEachChild(func(child *ast.Node) bool {
		irNode.Children = append(irNode.Children, e.emitNode(child))
		return false
	})

	return irNode
}

func (e *Emitter) emitBlock(node *ast.Node) *Block {
	if node == nil {
		return nil
	}

	block := &Block{}

	node.ForEachChild(func(child *ast.Node) bool {
		if stmt := e.emitStatement(child); stmt != nil {
			block.Statements = append(block.Statements, *stmt)
		}
		return false
	})

	return block
}

func (e *Emitter) emitStatement(node *ast.Node) *Statement {
	if node == nil {
		return nil
	}

	stmt := &Statement{Kind: node.Kind.String()}

	switch node.Kind {
	case ast.KindReturnStatement:
		ret := &ReturnStmt{}
		if expr := node.Expression(); expr != nil {
			ret.Value = e.emitExpression(expr)
		}
		stmt.Return = ret

	case ast.KindVariableStatement:
		vs := node.AsVariableStatement()
		if vs != nil && vs.DeclarationList != nil {
			declList := vs.DeclarationList.AsVariableDeclarationList()
			if declList != nil && declList.Declarations != nil {
				for _, decl := range declList.Declarations.Nodes {
					if decl.Kind == ast.KindVariableDeclaration {
						stmt.Declaration = e.emitVariableDeclaration(decl)
					}
				}
			}
		}

	case ast.KindExpressionStatement:
		stmt.Expression = e.emitExpression(node.Expression())

	case ast.KindIfStatement:
		stmt.If = e.emitIfStatement(node)

	case ast.KindWhileStatement:
		stmt.While = e.emitWhileStatement(node)

	case ast.KindForStatement:
		stmt.For = e.emitForStatement(node)

	case ast.KindBlock:
		stmt.Block = e.emitBlock(node)
	}

	return stmt
}

func (e *Emitter) emitVariableDeclaration(node *ast.Node) *Declaration {
	if node == nil {
		return nil
	}

	nameNode := node.Name()
	if nameNode == nil {
		return nil
	}

	decl := &Declaration{}

	switch nameNode.Kind {
	case ast.KindIdentifier:
		decl.Name = nameNode.Text()
		if sym := e.checker.GetSymbolAtLocation(nameNode); sym != nil {
			decl.Symbol = e.getOrCreateSymbolId(sym)
		}
	}

	if e.checker != nil {
		if t := e.checker.GetTypeAtLocation(node); t != nil {
			decl.Type = e.getOrCreateTypeId(t)
		}
	}

	if init := node.Initializer(); init != nil {
		decl.Init = e.emitExpression(init)
	}

	return decl
}

func (e *Emitter) emitExpression(node *ast.Node) *Expression {
	if node == nil {
		return nil
	}

	expr := &Expression{
		Kind: node.Kind.String(),
	}

	if e.checker != nil {
		if t := e.checker.GetTypeAtLocation(node); t != nil {
			expr.Type = e.getOrCreateTypeId(t)
		}
		if sym := e.checker.GetSymbolAtLocation(node); sym != nil {
			expr.Symbol = e.getOrCreateSymbolId(sym)
		}
	}

	switch node.Kind {
	case ast.KindIdentifier:
		expr.Text = node.Text()

	case ast.KindStringLiteral:
		expr.Text = node.Text()
		expr.Value = node.Text()

	case ast.KindNumericLiteral:
		expr.Text = node.Text()

	case ast.KindTrueKeyword:
		expr.Text = "true"
		expr.Value = true

	case ast.KindFalseKeyword:
		expr.Text = "false"
		expr.Value = false

	case ast.KindNullKeyword:
		expr.Text = "null"
		expr.Value = nil

	case ast.KindBinaryExpression:
		bin := node.AsBinaryExpression()
		expr.Operator = bin.OperatorToken.Kind.String()
		expr.Left = e.emitExpression(bin.Left)
		expr.Right = e.emitExpression(bin.Right)

	case ast.KindPrefixUnaryExpression:
		unary := node.AsPrefixUnaryExpression()
		expr.Operator = unary.Operator.String()
		expr.Operand = e.emitExpression(unary.Operand)

	case ast.KindPostfixUnaryExpression:
		unary := node.AsPostfixUnaryExpression()
		expr.Operator = unary.Operator.String()
		expr.Operand = e.emitExpression(unary.Operand)

	case ast.KindCallExpression:
		expr.Callee = e.emitExpression(node.Expression())
		for _, arg := range node.Arguments() {
			expr.Arguments = append(expr.Arguments, *e.emitExpression(arg))
		}

	case ast.KindPropertyAccessExpression:
		expr.Object = e.emitExpression(node.Expression())
		expr.Property = e.emitExpression(node.Name())

	case ast.KindElementAccessExpression:
		elem := node.AsElementAccessExpression()
		expr.Object = e.emitExpression(elem.Expression)
		expr.Property = e.emitExpression(elem.ArgumentExpression)

	case ast.KindObjectLiteralExpression:
		for _, prop := range node.Properties() {
			propInit := PropertyInit{}
			name := prop.Name()
			if name != nil {
				propInit.Name = name.Text()
			}
			switch prop.Kind {
			case ast.KindPropertyAssignment:
				pa := prop.AsPropertyAssignment()
				if pa.Initializer != nil {
					propInit.Value = e.emitExpression(pa.Initializer)
				} else if name != nil {
					propInit.Value = e.emitExpression(name)
				}
			case ast.KindShorthandPropertyAssignment:
				spa := prop.AsShorthandPropertyAssignment()
				if spa.ObjectAssignmentInitializer != nil {
					propInit.Value = e.emitExpression(spa.ObjectAssignmentInitializer)
				} else if name != nil {
					propInit.Value = e.emitExpression(name)
				}
			}
			expr.Properties = append(expr.Properties, propInit)
		}

	case ast.KindArrayLiteralExpression:
		for _, elem := range node.Elements() {
			expr.Elements = append(expr.Elements, *e.emitExpression(elem))
		}

	case ast.KindNewExpression:
		expr.Callee = e.emitExpression(node.Expression())
		for _, arg := range node.Arguments() {
			expr.Arguments = append(expr.Arguments, *e.emitExpression(arg))
		}

	case ast.KindParenthesizedExpression:
		return e.emitExpression(node.Expression())

	case ast.KindTemplateExpression:
		template := node.AsTemplateExpression()
		if template != nil && template.Head != nil {
			expr.Text = template.Head.Text()
		}

	case ast.KindArrowFunction:
		arrow := node.AsArrowFunction()
		for _, param := range arrow.Parameters.Nodes {
			irParam := Param{}
			if name := param.Name(); name != nil {
				irParam.Name = name.Text()
			}
			if e.checker != nil {
				if t := e.checker.GetTypeAtLocation(param); t != nil {
					irParam.Type = e.getOrCreateTypeId(t)
				}
			}
			expr.Parameters = append(expr.Parameters, irParam)
		}
		if arrow.Body != nil {
			if block := arrow.Body.AsBlock(); block != nil {
				expr.Body = e.emitBlock(arrow.Body)
			} else {
				expr.Body = &Block{
					Statements: []Statement{{
						Kind:   "ReturnStatement",
						Return: &ReturnStmt{Value: e.emitExpression(arrow.Body)},
					}},
				}
			}
		}

	case ast.KindConditionalExpression:
		cond := node.AsConditionalExpression()
		expr.Condition = e.emitExpression(cond.Condition)
		expr.Consequent = e.emitExpression(cond.WhenTrue)
		expr.Alternate = e.emitExpression(cond.WhenFalse)

	case ast.KindSpreadElement:
		expr.Operand = e.emitExpression(node.Expression())

	case ast.KindSpreadAssignment:
		expr.Operand = e.emitExpression(node.Expression())

	case ast.KindTypeOfExpression:
		typeOf := node.AsTypeOfExpression()
		expr.Operand = e.emitExpression(typeOf.Expression)

	case ast.KindAwaitExpression:
		await := node.AsAwaitExpression()
		expr.Operand = e.emitExpression(await.Expression)

	case ast.KindYieldExpression:
		yield := node.AsYieldExpression()
		if yield.Expression != nil {
			expr.Operand = e.emitExpression(yield.Expression)
		}

	case ast.KindMetaProperty:
		meta := node.AsMetaProperty()
		if meta.KeywordToken == ast.KindNewKeyword {
			if name := meta.Name(); name != nil {
				expr.Text = "new." + name.Text()
			}
		} else if meta.KeywordToken == ast.KindImportKeyword {
			if name := meta.Name(); name != nil {
				expr.Text = "import." + name.Text()
			}
		}

	case ast.KindDeleteExpression:
		deleteExpr := node.AsDeleteExpression()
		expr.Operand = e.emitExpression(deleteExpr.Expression)

	case ast.KindVoidExpression:
		voidExpr := node.AsVoidExpression()
		expr.Operand = e.emitExpression(voidExpr.Expression)

	case ast.KindNonNullExpression:
		nonNull := node.AsNonNullExpression()
		expr.Operand = e.emitExpression(nonNull.Expression)

	case ast.KindAsExpression:
		asExpr := node.AsAsExpression()
		expr.Operand = e.emitExpression(asExpr.Expression)

	case ast.KindTypeAssertionExpression:
		typeAssert := node.AsTypeAssertion()
		expr.Operand = e.emitExpression(typeAssert.Expression)

	case ast.KindTaggedTemplateExpression:
		tagged := node.AsTaggedTemplateExpression()
		expr.Callee = e.emitExpression(tagged.Tag)
	}

	return expr
}

func (e *Emitter) emitIfStatement(node *ast.Node) *IfStmt {
	if node == nil {
		return nil
	}

	ifNode := node.AsIfStatement()
	ifStmt := &IfStmt{
		Condition: e.emitExpression(ifNode.Expression),
		Then:      e.emitBlock(ifNode.ThenStatement),
	}

	if ifNode.ElseStatement != nil {
		if ifNode.ElseStatement.Kind == ast.KindBlock {
			ifStmt.Else = e.emitBlock(ifNode.ElseStatement)
		} else {
			ifStmt.Else = &Block{
				Statements: []Statement{*e.emitStatement(ifNode.ElseStatement)},
			}
		}
	}

	return ifStmt
}

func (e *Emitter) emitWhileStatement(node *ast.Node) *WhileStmt {
	if node == nil {
		return nil
	}

	whileNode := node.AsWhileStatement()
	return &WhileStmt{
		Condition: e.emitExpression(whileNode.Expression),
		Body:      e.emitBlock(whileNode.Statement),
	}
}

func (e *Emitter) emitForStatement(node *ast.Node) *ForStmt {
	if node == nil {
		return nil
	}

	forNode := node.AsForStatement()
	forStmt := &ForStmt{}

	if forNode.Initializer != nil {
		forStmt.Init = e.emitStatement(forNode.Initializer)
	}
	if forNode.Condition != nil {
		forStmt.Condition = e.emitExpression(forNode.Condition)
	}
	if forNode.Incrementor != nil {
		forStmt.Update = e.emitStatement(forNode.Incrementor)
	}
	if forNode.Statement != nil {
		forStmt.Body = e.emitBlock(forNode.Statement)
	}

	return forStmt
}

func (e *Emitter) emitTypeSignature(node *ast.Node) *TypeSignature {
	if node == nil {
		return nil
	}

	sig := e.checker.GetSignatureFromDeclaration(node)
	if sig == nil {
		return nil
	}

	irSig := &TypeSignature{
		Kind: node.Kind.String(),
	}

	for _, param := range sig.Parameters() {
		irParam := Param{Name: param.Name}
		if paramType := e.checker.GetTypeOfSymbol(param); paramType != nil {
			irParam.Type = e.getOrCreateTypeId(paramType)
		}
		irSig.Parameters = append(irSig.Parameters, irParam)
	}

	if returnType := e.checker.GetReturnTypeOfSignature(sig); returnType != nil {
		irSig.ReturnType = e.getOrCreateTypeId(returnType)
	}

	return irSig
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
		Kind:  e.typeKindToString(t),
		Flags: uint32(t.Flags()),
	}

	if t.Symbol() != nil {
		irType.Name = t.Symbol().Name
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
	case flags&ast.SymbolFlagsVariable != 0, flags&ast.SymbolFlagsBlockScopedVariable != 0:
		return "variable"
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

func (e *Emitter) MarshalJSON() ([]byte, error) {
	return json.MarshalIndent(e.irProgram, "", "  ")
}
