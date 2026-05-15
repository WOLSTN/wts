package compiler

import (
	"github.com/wolstn/wts/internal/ast"
	"github.com/wolstn/wts/internal/outputpaths"
	"github.com/wolstn/wts/internal/printer"
	"github.com/wolstn/wts/internal/tracing"
)

type EmitOnly byte

const (
	EmitOnlyNone EmitOnly = iota
	EmitOnlyDeclaration
	EmitOnlyJS
	EmitOnlyForcedDts
)

type emitter struct {
	host               EmitHost
	emitOnly           EmitOnly
	emitterDiagnostics ast.DiagnosticsCollection
	writer             printer.EmitTextWriter
	paths              *outputpaths.OutputPaths
	sourceFile         *ast.SourceFile
	emitResult         EmitResult
	writeFile          func(fileName string, text string, data *WriteFileData) error
	tr                 *tracing.Tracing
}

func (e *emitter) emit() {
	e.emitResult.EmitSkipped = true
}
