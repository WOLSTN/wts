package checker

import "github.com/wolstn/wts/internal/ast"

func (c *Checker) GetTypeOfSymbol(symbol *ast.Symbol) *Type {
	return c.getTypeOfSymbol(symbol)
}

func (c *Checker) GetDeclaredTypeOfSymbol(symbol *ast.Symbol) *Type {
	return c.getDeclaredTypeOfSymbol(symbol)
}

func (c *Checker) GetPropertiesOfType(t *Type) []*ast.Symbol {
	return c.getPropertiesOfType(t)
}

func (c *Checker) GetBaseTypes(t *Type) []*Type {
	return c.getBaseTypes(t)
}

func (c *Checker) GetLocalTypeParametersOfClassOrInterfaceOrTypeAlias(symbol *ast.Symbol) []*Type {
	return c.getLocalTypeParametersOfClassOrInterfaceOrTypeAlias(symbol)
}

func (c *Checker) GetResolvedSignature(node *ast.Node) *Signature {
	return c.getResolvedSignature(node, nil, CheckModeNormal)
}

func (c *Checker) GetReturnTypeOfSignature(sig *Signature) *Type {
	return c.getReturnTypeOfSignature(sig)
}

func (c *Checker) GetTypeFromTypeNode(node *ast.Node) *Type {
	return c.getTypeFromTypeNode(node)
}

func (c *Checker) GetConstraintOfTypeParameter(typeParameter *Type) *Type {
	return c.getConstraintOfTypeParameter(typeParameter)
}

func (c *Checker) GetDefaultFromTypeParameter(typeParameter *Type) *Type {
	return c.getDefaultFromTypeParameter(typeParameter)
}

func (c *Checker) GetTypeArguments(t *Type) []*Type {
	return c.getTypeArguments(t)
}

func (c *Checker) GetIndexInfosOfType(t *Type) []*IndexInfo {
	return c.getIndexInfosOfType(t)
}

func (c *Checker) GetSignaturesOfType(t *Type, kind SignatureKind) []*Signature {
	return c.getSignaturesOfType(t, kind)
}

func (c *Checker) GetMergedSymbol(symbol *ast.Symbol) *ast.Symbol {
	return c.getMergedSymbol(symbol)
}
