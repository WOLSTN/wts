package parsetestutil

import (
	"github.com/wolstn/wts/internal/ast"
	"github.com/wolstn/wts/internal/core"
	"github.com/wolstn/wts/internal/parser"
	"github.com/wolstn/wts/internal/tspath"
)

func ParseTypeScript(code string) *ast.SourceFile {
	return parser.ParseSourceFile(
		ast.SourceFileParseOptions{
			FileName: "test.ts",
			Path:     tspath.Path("test.ts"),
		},
		code,
		core.ScriptKindTS,
	)
}

func ParseTypeScriptFile(fileName, code string) *ast.SourceFile {
	return parser.ParseSourceFile(
		ast.SourceFileParseOptions{
			FileName: fileName,
			Path:     tspath.Path(fileName),
		},
		code,
		core.ScriptKindTS,
	)
}

func ParseJSX(code string) *ast.SourceFile {
	return parser.ParseSourceFile(
		ast.SourceFileParseOptions{
			FileName: "test.tsx",
			Path:     tspath.Path("test.tsx"),
			ExternalModuleIndicatorOptions: ast.ExternalModuleIndicatorOptions{
				JSX: true,
			},
		},
		code,
		core.ScriptKindTSX,
	)
}
