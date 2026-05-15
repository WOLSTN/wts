package emittestutil

import (
	"fmt"

	"github.com/wolstn/wts/internal/ast"
	"github.com/wolstn/wts/internal/core"
	"github.com/wolstn/wts/internal/printer"
)

type stringWriter interface {
	String() string
}

func Emit(sourceFile *ast.SourceFile) string {
	opts := printer.PrinterOptions{
		Target: core.ScriptTargetESNext,
	}
	ctx := printer.NewEmitContext()
	pr := printer.NewPrinter(opts, printer.PrintHandlers{}, ctx)
	writer, release := printer.GetSingleLineStringWriter()
	defer release()
	pr.Write(sourceFile.AsNode(), sourceFile, writer, nil)
	return writer.(stringWriter).String()
}

func EmitNode(node *ast.Node) string {
	opts := printer.PrinterOptions{
		Target: core.ScriptTargetESNext,
	}
	ctx := printer.NewEmitContext()
	pr := printer.NewPrinter(opts, printer.PrintHandlers{}, ctx)
	writer, release := printer.GetSingleLineStringWriter()
	defer release()
	pr.Write(node, nil, writer, nil)
	return writer.(stringWriter).String()
}

func CompareOutput(t interface{ Errorf(format string, args ...any) }, expected, actual string) {
	if expected != actual {
		t.Errorf("output mismatch:\nexpected:\n%s\nactual:\n%s", expected, actual)
	}
}

func FormatOutput(output string) string {
	return fmt.Sprintf("%q", output)
}
