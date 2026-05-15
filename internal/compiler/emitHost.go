package compiler

import (
	"context"

	"github.com/wolstn/wts/internal/ast"
	"github.com/wolstn/wts/internal/core"
)

type EmitHost interface {
	SourceFileMayBeEmittedHost
	SourceFiles() []*ast.SourceFile
	UseCaseSensitiveFileNames() bool
	GetCurrentDirectory() string
	CommonSourceDirectory() string
	IsEmitBlocked(file string) bool
}

type SourceFileMayBeEmittedHost interface {
	Options() *core.CompilerOptions
}

type emitHost struct {
	program *Program
}

func newEmitHost(ctx context.Context, program *Program, file *ast.SourceFile) (*emitHost, func()) {
	return &emitHost{program: program}, func() {}
}

func (host *emitHost) Options() *core.CompilerOptions {
	return host.program.Options()
}

func (host *emitHost) SourceFiles() []*ast.SourceFile {
	return host.program.GetSourceFiles()
}

func (host *emitHost) UseCaseSensitiveFileNames() bool {
	return host.program.UseCaseSensitiveFileNames()
}

func (host *emitHost) GetCurrentDirectory() string {
	return host.program.GetCurrentDirectory()
}

func (host *emitHost) CommonSourceDirectory() string {
	return host.program.CommonSourceDirectory()
}

func (host *emitHost) IsEmitBlocked(file string) bool {
	return false
}

func sourceFileMayBeEmitted(sourceFile *ast.SourceFile, host SourceFileMayBeEmittedHost, forceDtsEmit bool) bool {
	return false
}

func getSourceFilesToEmit(host SourceFileMayBeEmittedHost, targetSourceFile *ast.SourceFile, forceDtsEmit bool) []*ast.SourceFile {
	return nil
}

func getDeclarationDiagnostics(host EmitHost, file *ast.SourceFile) []*ast.Diagnostic {
	return nil
}
