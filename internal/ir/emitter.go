package ir

import (
	"context"
	"encoding/json"
	"fmt"

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
}

type File struct {
	Path   string `json:"path"`
	Source string `json:"source,omitempty"`
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
	Name string `json:"name"`
	Type string `json:"type,omitempty"`
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
	Signature    string  `json:"signature,omitempty"`
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
	Name       string   `json:"name"`
	Symbol     string   `json:"symbol"`
	Parameters []Param  `json:"parameters,omitempty"`
	ReturnType string   `json:"returnType,omitempty"`
	TypeParams []string `json:"typeParams,omitempty"`
	Signature  string   `json:"signature,omitempty"`
	IsAsync    bool     `json:"isAsync,omitempty"`
	IsGenerator bool    `json:"isGenerator,omitempty"`
}

type Variable struct {
	Name    string `json:"name"`
	Symbol  string `json:"symbol"`
	Type    string `json:"type,omitempty"`
	IsConst bool   `json:"isConst,omitempty"`
}

type Emitter struct {
	program     *compiler.Program
	checker     *checker.Checker
	checkerData *checker.CheckerData
	irProgram   *Program
	typeMap     map[*checker.Type]string
	symbolMap   map[*ast.Symbol]string
	sigMap      map[*checker.Signature]string
	typeIdGen   int
	symbolIdGen int
	sigIdGen    int
}

func NewEmitter(program *compiler.Program) *Emitter {
	return &Emitter{
		program:   program,
		irProgram: &Program{Version: Version},
		typeMap:   make(map[*checker.Type]string),
		symbolMap: make(map[*ast.Symbol]string),
		sigMap:    make(map[*checker.Signature]string),
	}
}

func (e *Emitter) Emit() (*Program, error) {
	ctx := context.Background()
	chk, release := e.program.GetTypeChecker(ctx)
	defer release()
	e.checker = chk

	e.checkerData = checker.CollectAllCheckerData(e.checker)

	e.emitAllTypes(e.checkerData.TypeMaps)
	for i := range e.checkerData.SymbolArena {
		sym := &e.checkerData.SymbolArena[i]
		e.getOrCreateSymbolId(sym)
	}
	for i := range e.checkerData.SignatureArena {
		sig := &e.checkerData.SignatureArena[i]
		e.getOrCreateSignatureId(sig)
	}
	if e.checkerData.Globals != nil {
		for name, sym := range e.checkerData.Globals {
			e.irProgram.Globals = append(e.irProgram.Globals, &Global{
				Name:   name,
				Symbol: e.getOrCreateSymbolId(sym),
			})
		}
	}

	e.emitSourceFiles()

	return e.irProgram, nil
}

func (e *Emitter) emitAllTypes(typeMaps map[string]any) {
	for name, m := range typeMaps {
		e.emitTypeMap(name, m)
	}
}

func (e *Emitter) emitTypeMap(kind string, m any) {
	switch v := m.(type) {
	case map[string]*checker.Type:
		for _, t := range v {
			e.getOrCreateTypeId(t)
		}
	case map[jsnum.Number]*checker.Type:
		for _, t := range v {
			e.getOrCreateTypeId(t)
		}
	case map[jsnum.PseudoBigInt]*checker.Type:
		for _, t := range v {
			e.getOrCreateTypeId(t)
		}
	case map[checker.EnumLiteralKey]*checker.Type:
		for _, t := range v {
			e.getOrCreateTypeId(t)
		}
	case map[checker.CacheHashKey]*checker.Type:
		for _, t := range v {
			e.getOrCreateTypeId(t)
		}
	case map[checker.StringMappingKey]*checker.Type:
		for _, t := range v {
			e.getOrCreateTypeId(t)
		}
	case map[checker.CachedTypeKey]*checker.Type:
		for _, t := range v {
			e.getOrCreateTypeId(t)
		}
	case map[checker.NarrowedTypeKey]*checker.Type:
		for _, t := range v {
			e.getOrCreateTypeId(t)
		}
	case map[checker.InstantiationExpressionKey]*checker.Type:
		for _, t := range v {
			e.getOrCreateTypeId(t)
		}
	case map[checker.SubstitutionTypeKey]*checker.Type:
		for _, t := range v {
			e.getOrCreateTypeId(t)
		}
	case map[checker.UnionOfUnionKey]*checker.Type:
		for _, t := range v {
			e.getOrCreateTypeId(t)
		}
	case map[checker.PropertiesTypesKey]*checker.Type:
		for _, t := range v {
			e.getOrCreateTypeId(t)
		}
	case map[checker.AssignmentReducedKey]*checker.Type:
		for _, t := range v {
			e.getOrCreateTypeId(t)
		}
	case map[checker.ReverseMappedTypeKey]*checker.Type:
		for _, t := range v {
			e.getOrCreateTypeId(t)
		}
	case map[*ast.Symbol]*checker.Type:
		for _, t := range v {
			e.getOrCreateTypeId(t)
		}
	case map[*ast.Node]*checker.Type:
		for _, t := range v {
			e.getOrCreateTypeId(t)
		}
	case map[checker.CachedSignatureKey]*checker.Signature:
		for _, s := range v {
			e.getOrCreateSignatureId(s)
		}
	case map[checker.CacheHashKey][]*checker.Type:
		for _, types := range v {
			for _, t := range types {
				e.getOrCreateTypeId(t)
			}
		}
	}
}

func (e *Emitter) emitSourceFiles() {
	for _, sf := range e.program.GetSourceFiles() {
		if sf.IsDeclarationFile {
			continue
		}
		e.irProgram.Files = append(e.irProgram.Files, &File{
			Path:   sf.FileName(),
			Source: string(sf.Text()),
		})
		e.emitFileImports(sf)
		e.emitFileExports(sf)
	}

	for _, sym := range e.checkerData.Globals {
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

func (e *Emitter) emitClassFromSymbol(sym *ast.Symbol) {
	irClass := &Class{
		Name:   sym.Name,
		Symbol: e.getOrCreateSymbolId(sym),
	}

	declaredType := e.checker.GetDeclaredTypeOfSymbol(sym)
	if declaredType != nil {
		typeParams := e.checker.GetLocalTypeParametersOfClassOrInterfaceOrTypeAlias(sym)
		for _, tp := range typeParams {
			irClass.TypeParams = append(irClass.TypeParams, e.getOrCreateTypeId(tp))
		}

		baseTypes := e.checker.GetBaseTypes(declaredType)
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

	e.irProgram.Classes = append(e.irProgram.Classes, irClass)
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
	}

	sig := e.checker.GetResolvedSignature(sym.Declarations[0])
	if sig != nil {
		method.Signature = e.getOrCreateSignatureId(sig)
		for _, param := range sig.Parameters() {
			p := Param{Name: param.Name}
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
		flags := ast.GetFunctionFlags(decl)
		fn.IsAsync = flags&ast.FunctionFlagsAsync != 0
		fn.IsGenerator = flags&ast.FunctionFlagsGenerator != 0
	}

	sig := e.checker.GetResolvedSignature(sym.Declarations[0])
	if sig != nil {
		fn.Signature = e.getOrCreateSignatureId(sig)
		for _, param := range sig.Parameters() {
			p := Param{Name: param.Name}
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

	e.irProgram.Functions = append(e.irProgram.Functions, fn)
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
		Id:      id,
		Flags:   uint32(t.Flags()),
	}

	flags := t.Flags()

	switch {
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
	case flags&checker.TypeFlagsLiteral != 0:
		irType.Kind = "literal"
		e.emitLiteralType(t, irType)
	case flags&checker.TypeFlagsStringLiteral != 0:
		irType.Kind = "stringLiteral"
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
			propType := e.checker.GetTypeOfSymbol(prop)
			irProp := &Property{
				Name:       prop.Name,
				Symbol:     e.getOrCreateSymbolId(prop),
				IsOptional: prop.Flags&ast.SymbolFlagsOptional != 0,
			}
			if propType != nil {
				irProp.Type = e.getOrCreateTypeId(propType)
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
	if id, ok := e.symbolMap[sym]; ok {
		return id
	}

	e.symbolIdGen++
	id := fmt.Sprintf("s%d", e.symbolIdGen)
	e.symbolMap[sym] = id

	irSym := &Symbol{
		Id:         id,
		Name:       sym.Name,
		Flags:      uint64(sym.Flags),
		CheckFlags: uint32(sym.CheckFlags),
		Kind:       e.symbolKindToString(sym),
	}

	if e.checker != nil {
		if t := e.checker.GetTypeOfSymbol(sym); t != nil {
			irSym.Type = e.getOrCreateTypeId(t)
		}
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
	if id, ok := e.sigMap[sig]; ok {
		return id
	}

	e.sigIdGen++
	id := fmt.Sprintf("sig%d", e.sigIdGen)
	e.sigMap[sig] = id

	irSig := &Signature{
		Id:      id,
		Kind:    "call",
		HasRest: sig.HasRestParameter(),
	}

	for _, param := range sig.Parameters() {
		p := Param{Name: param.Name}
		if t := e.checker.GetTypeOfSymbol(param); t != nil {
			p.Type = e.getOrCreateTypeId(t)
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
	return json.MarshalIndent(e.irProgram, "", "  ")
}
