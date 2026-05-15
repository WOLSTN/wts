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
	Id             string          `json:"id"`
	Kind           string          `json:"kind"`
	Flags          uint32          `json:"flags"`
	Name           string          `json:"name,omitempty"`
	Members        []string        `json:"members,omitempty"`
	Properties     []*Property     `json:"properties,omitempty"`
	Signatures     []*TypeSignature `json:"signatures,omitempty"`
	TypeArgs       []string        `json:"typeArgs,omitempty"`
	TypeParams     []string        `json:"typeParams,omitempty"`
	Target         string          `json:"target,omitempty"`
	Value          any             `json:"value,omitempty"`
	ReturnType     string          `json:"returnType,omitempty"`
	Constraint     string          `json:"constraint,omitempty"`
	Default        string          `json:"default,omitempty"`
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

type Param struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type Block struct {
	Statements []Statement `json:"statements"`
}

type Statement struct {
	Kind       string      `json:"kind"`
	Expression *Expression `json:"expression,omitempty"`
	Declaration *Declaration `json:"declaration,omitempty"`
	Return     *ReturnStmt `json:"return,omitempty"`
	If         *IfStmt     `json:"if,omitempty"`
	While      *WhileStmt  `json:"while,omitempty"`
	For        *ForStmt    `json:"for,omitempty"`
	Block      *Block      `json:"block,omitempty"`
}

type Expression struct {
	Kind       string      `json:"kind"`
	Text       string      `json:"text,omitempty"`
	Type       string      `json:"type,omitempty"`
	Symbol     string      `json:"symbol,omitempty"`
	Value      any         `json:"value,omitempty"`
	Operator   string      `json:"operator,omitempty"`
	Left       *Expression `json:"left,omitempty"`
	Right      *Expression `json:"right,omitempty"`
	Operand    *Expression `json:"operand,omitempty"`
	Callee     *Expression `json:"callee,omitempty"`
	Arguments  []Expression `json:"arguments,omitempty"`
	Object     *Expression `json:"object,omitempty"`
	Property   *Expression `json:"property,omitempty"`
	Elements   []Expression `json:"elements,omitempty"`
	Properties []PropertyInit `json:"properties,omitempty"`
}

type PropertyInit struct {
	Name  string      `json:"name"`
	Value *Expression `json:"value,omitempty"`
}

type Declaration struct {
	Name       string      `json:"name"`
	Symbol     string      `json:"symbol,omitempty"`
	Type       string      `json:"type,omitempty"`
	Init       *Expression `json:"init,omitempty"`
	Params     []Param     `json:"params,omitempty"`
	ReturnType string      `json:"returnType,omitempty"`
	Body       *Block      `json:"body,omitempty"`
}

type ReturnStmt struct {
	Value *Expression `json:"value,omitempty"`
}

type IfStmt struct {
	Condition  *Expression `json:"condition"`
	Then       *Block      `json:"then"`
	Else       *Block      `json:"else,omitempty"`
}

type WhileStmt struct {
	Condition *Expression `json:"condition"`
	Body      *Block      `json:"body"`
}

type ForStmt struct {
	Init       *Statement  `json:"init,omitempty"`
	Condition  *Expression `json:"condition,omitempty"`
	Update     *Statement  `json:"update,omitempty"`
	Body       *Block      `json:"body"`
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
		e.collectFunctions(sourceFile)
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
		irNode.Text = node.Text()
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

	flags := t.Flags()

	switch {
	case flags&checker.TypeFlagsUnion != 0:
		e.emitUnionTypeInfo(irType, t)
	case flags&checker.TypeFlagsIntersection != 0:
		e.emitIntersectionTypeInfo(irType, t)
	case flags&checker.TypeFlagsObject != 0:
		e.emitObjectTypeInfo(irType, t)
	case flags&checker.TypeFlagsTypeParameter != 0:
		e.emitTypeParameterInfo(irType, t)
	case flags&checker.TypeFlagsLiteral != 0:
		e.emitLiteralTypeInfo(irType, t)
	case flags&checker.TypeFlagsStringLiteral != 0:
		e.emitStringLiteralTypeInfo(irType, t)
	case flags&checker.TypeFlagsNumberLiteral != 0:
		e.emitNumberLiteralTypeInfo(irType, t)
	}

	e.irProgram.Types = append(e.irProgram.Types, irType)
	return id
}

func (e *Emitter) emitUnionTypeInfo(irType *Type, t *checker.Type) {
	unionType := t.AsUnionType()
	if unionType == nil {
		return
	}
	types := unionType.Types()
	for _, memberType := range types {
		irType.Members = append(irType.Members, e.getOrCreateTypeId(memberType))
	}
}

func (e *Emitter) emitIntersectionTypeInfo(irType *Type, t *checker.Type) {
	intersectionType := t.AsIntersectionType()
	if intersectionType == nil {
		return
	}
	types := intersectionType.Types()
	for _, memberType := range types {
		irType.Members = append(irType.Members, e.getOrCreateTypeId(memberType))
	}
}

func (e *Emitter) emitObjectTypeInfo(irType *Type, t *checker.Type) {
	structuredType := t.AsStructuredType()
	if structuredType == nil {
		return
	}

	for _, prop := range structuredType.Properties() {
		irProp := &Property{
			Name:   prop.Name,
			Symbol: e.getOrCreateSymbolId(prop),
		}
		if e.checker != nil {
			if propType := e.checker.GetTypeOfSymbol(prop); propType != nil {
				irProp.Type = e.getOrCreateTypeId(propType)
			}
		}
		irProp.IsOptional = prop.Flags&ast.SymbolFlagsOptional != 0
		irType.Properties = append(irType.Properties, irProp)
	}

	callSigs := structuredType.CallSignatures()
	for _, sig := range callSigs {
		irSig := e.emitSignatureInfo(sig, "call")
		irType.Signatures = append(irType.Signatures, irSig)
	}

	constructSigs := structuredType.ConstructSignatures()
	for _, sig := range constructSigs {
		irSig := e.emitSignatureInfo(sig, "construct")
		irType.Signatures = append(irType.Signatures, irSig)
	}
}

func (e *Emitter) emitTypeParameterInfo(irType *Type, t *checker.Type) {
}

func (e *Emitter) emitLiteralTypeInfo(irType *Type, t *checker.Type) {
	literalType := t.AsLiteralType()
	if literalType == nil {
		return
	}
	val := literalType.Value()
	switch v := val.(type) {
	case string, bool, nil:
		irType.Value = v
	case float64:
		irType.Value = v
	default:
		irType.Value = fmt.Sprintf("%v", v)
	}
}

func (e *Emitter) emitStringLiteralTypeInfo(irType *Type, t *checker.Type) {
	if t.IsStringLiteral() {
		irType.Kind = "stringLiteral"
	}
}

func (e *Emitter) emitNumberLiteralTypeInfo(irType *Type, t *checker.Type) {
	if t.IsNumberLiteral() {
		irType.Kind = "numberLiteral"
	}
}

func (e *Emitter) emitSignatureInfo(sig *checker.Signature, kind string) *TypeSignature {
	irSig := &TypeSignature{
		Kind: kind,
	}

	for _, param := range sig.Parameters() {
		irParam := Param{
			Name: param.Name,
		}
		if e.checker != nil {
			if paramType := e.checker.GetTypeOfSymbol(param); paramType != nil {
				irParam.Type = e.getOrCreateTypeId(paramType)
			}
		}
		irSig.Parameters = append(irSig.Parameters, irParam)
	}

	if e.checker != nil {
		if returnType := e.checker.GetReturnTypeOfSignature(sig); returnType != nil {
			irSig.ReturnType = e.getOrCreateTypeId(returnType)
		}
	}

	for _, tp := range sig.TypeParameters() {
		irSig.TypeParameters = append(irSig.TypeParameters, e.getOrCreateTypeId(tp))
	}

	return irSig
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

func (p *Program) ToJSON() ([]byte, error) {
	return json.MarshalIndent(p, "", "  ")
}

func (e *Emitter) collectFunctions(file *ast.SourceFile) {
	file.AsNode().ForEachChild(func(child *ast.Node) bool {
		if child.Kind == ast.KindFunctionDeclaration {
			e.emitFunction(child)
		}
		return false
	})
}

func (e *Emitter) emitFunction(node *ast.Node) {
	if node == nil {
		return
	}

	nameNode := node.Name()
	if nameNode == nil {
		return
	}

	sym := e.checker.GetSymbolAtLocation(nameNode)
	if sym == nil {
		return
	}

	irFunc := &Function{
		Name:   sym.Name,
		Symbol: e.getOrCreateSymbolId(sym),
	}

	sig := e.checker.GetSignatureFromDeclaration(node)
	if sig != nil {
		for _, param := range sig.Parameters() {
			irParam := Param{
				Name: param.Name,
			}
			if paramType := e.checker.GetTypeOfSymbol(param); paramType != nil {
				irParam.Type = e.getOrCreateTypeId(paramType)
			}
			irFunc.Parameters = append(irFunc.Parameters, irParam)
		}

		if returnType := e.checker.GetReturnTypeOfSignature(sig); returnType != nil {
			irFunc.ReturnType = e.getOrCreateTypeId(returnType)
		}
	}

	body := node.Body()
	if body != nil {
		irFunc.Body = e.emitBlock(body)
	}

	e.irProgram.Functions = append(e.irProgram.Functions, irFunc)
}

func (e *Emitter) emitBlock(node *ast.Node) *Block {
	if node == nil || node.Kind != ast.KindBlock {
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
		for _, s := range node.Statements() {
			if s.Kind == ast.KindVariableDeclaration {
				stmt.Declaration = e.emitVariableDeclaration(s)
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
			if init := prop.Initializer(); init != nil {
				propInit.Value = e.emitExpression(init)
			} else if prop.Kind == ast.KindPropertyAssignment && name != nil {
				propInit.Value = e.emitExpression(name)
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
	forStmt := &ForStmt{
		Body: e.emitBlock(forNode.Statement),
	}

	if forNode.Initializer != nil {
		forStmt.Init = e.emitStatement(forNode.Initializer)
	}

	if forNode.Condition != nil {
		forStmt.Condition = e.emitExpression(forNode.Condition)
	}

	if forNode.Incrementor != nil {
		forStmt.Update = &Statement{
			Kind:       ast.KindExpressionStatement.String(),
			Expression: e.emitExpression(forNode.Incrementor),
		}
	}

	return forStmt
}
