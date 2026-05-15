package checker

import (
	"reflect"
	"unsafe"

	"github.com/wolstn/wts/internal/ast"
)

type CheckerData struct {
	SymbolArena   []ast.Symbol
	SignatureArena []Signature
	IndexInfoArena []IndexInfo
	Globals       ast.SymbolTable
	TypeMaps      map[string]any
	Symbols       map[string]*ast.Symbol
	Signatures    map[string]*Signature
	LinkStores    map[string]any
}

func CollectAllCheckerData(c *Checker) *CheckerData {
	d := &CheckerData{
		TypeMaps:   make(map[string]any),
		Symbols:    make(map[string]*ast.Symbol),
		Signatures: make(map[string]*Signature),
		LinkStores: make(map[string]any),
	}

	d.SymbolArena = c.symbolArena.Data()
	d.SignatureArena = c.signatureArena.Data()
	d.IndexInfoArena = c.indexInfoArena.Data()
	d.Globals = c.globals

	v := reflect.ValueOf(c).Elem()
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		name := t.Field(i).Name
		field := v.Field(i)

		switch field.Kind() {
		case reflect.Map:
			if field.IsNil() {
				continue
			}
			ft := field.Type()
			elemKind := ft.Elem().Kind()
			if elemKind == reflect.Ptr || elemKind == reflect.Slice {
				val := reflect.NewAt(ft, unsafe.Pointer(field.UnsafeAddr())).Elem()
				d.TypeMaps[name] = val.Interface()
			}

		case reflect.Ptr:
			if field.IsNil() {
				continue
			}
			etype := field.Type().Elem()
			switch etype.String() {
			case "ast.Symbol":
				val := reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
				if s, ok := val.Interface().(*ast.Symbol); ok && s != nil {
					d.Symbols[name] = s
				}
			case "Signature":
				val := reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
				if s, ok := val.Interface().(*Signature); ok && s != nil {
					d.Signatures[name] = s
				}
			}

		case reflect.Struct:
			typeName := t.Field(i).Type.String()
			if len(typeName) > 14 && typeName[:14] == "core.LinkStore" {
				val := reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
				d.LinkStores[name] = val.Addr().Interface()
			}
		}
	}

	return d
}
