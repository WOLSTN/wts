package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/wolstn/wts/internal/ast"
	"github.com/wolstn/wts/internal/bundled"
	"github.com/wolstn/wts/internal/compiler"
	"github.com/wolstn/wts/internal/core"
	"github.com/wolstn/wts/internal/diagnosticwriter"
	"github.com/wolstn/wts/internal/ir"
	"github.com/wolstn/wts/internal/tsoptions"
	"github.com/wolstn/wts/internal/tspath"
	"github.com/wolstn/wts/internal/vfs/osvfs"
)

func main() {
	os.Exit(runMain())
}

func runMain() int {
	args := os.Args[1:]

	if len(args) == 0 {
		printUsage()
		return 0
	}

	switch args[0] {
	case "--help", "-h":
		printUsage()
		return 0
	case "--version", "-v":
		printVersion()
		return 0
	case "check":
		return runCheck(args[1:])
	case "emit-ir":
		return runEmitIR(args[1:])
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", args[0])
		printUsage()
		return 1
	}
}

func printUsage() {
	fmt.Println(`wts - TypeScript Frontend for Wolstn Native Compiler

Usage:
  wts <command> [options]

Commands:
  check       Type-check TypeScript files
  emit-ir     Emit Wolstn IR from TypeScript files

Options:
  -h, --help      Show this help message
  -v, --version   Show version information

Examples:
  wts check main.ts
  wts check -p tsconfig.json
  wts emit-ir main.ts -o main.wir
  wts emit-ir -p tsconfig.json -o program.wir`)
}

func printVersion() {
	fmt.Println("wts version 0.1.0 (Wolstn TypeScript Frontend)")
}

type checkFlags struct {
	project string
	help    bool
}

func runCheck(args []string) int {
	fs := flag.NewFlagSet("check", flag.ExitOnError)
	project := fs.String("p", "", "path to tsconfig.json")
	help := fs.Bool("h", false, "show help")

	var flagArgs, posArgs []string
	for i := 0; i < len(args); i++ {
		if args[i] == "-h" || args[i] == "--help" {
			flagArgs = append(flagArgs, args[i])
		} else if args[i] == "-p" {
			flagArgs = append(flagArgs, args[i])
			if i+1 < len(args) {
				i++
				flagArgs = append(flagArgs, args[i])
			}
		} else {
			posArgs = append(posArgs, args[i])
		}
	}

	if err := fs.Parse(flagArgs); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		return 1
	}

	if *help {
		fmt.Println(`Usage: wts check [options] [files...]

Options:
  -p <path>   Path to tsconfig.json (default: auto-detect)
  -h          Show this help message

If no files or project are specified, checks all .ts files in the current directory.`)
		return 0
	}

	remainingArgs := posArgs
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
		return 1
	}

	program, err := createProgramFromArgs(*project, remainingArgs, cwd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	ctx := context.Background()
	var diags []*ast.Diagnostic
	for _, sf := range program.GetSourceFiles() {
		if !sf.IsDeclarationFile {
			diags = append(diags, program.GetSyntacticDiagnostics(ctx, sf)...)
			diags = append(diags, program.GetSemanticDiagnostics(ctx, sf)...)
		}
	}
	diags = append(diags, program.GetGlobalDiagnostics(ctx)...)

	if len(diags) > 0 {
		formatOpts := &diagnosticwriter.FormattingOptions{
			ComparePathsOptions: tspath.ComparePathsOptions{
				CurrentDirectory:          cwd,
				UseCaseSensitiveFileNames: false,
			},
			NewLine: "\n",
		}
		diagnosticwriter.FormatDiagnosticsWithColorAndContext(os.Stderr, diagnosticwriter.FromASTDiagnostics(diags), formatOpts)
		fmt.Fprintf(os.Stderr, "\nFound %d error(s)\n", len(diags))
		return 1
	}

	fmt.Println("No errors found.")
	return 0
}

type emitIRFlags struct {
	project   string
	output    string
	help      bool
	prune     bool
	compact   bool
	noSource  bool
	treeShake bool
}

func runEmitIR(args []string) int {
	fs := flag.NewFlagSet("emit-ir", flag.ExitOnError)
	project := fs.String("p", "", "path to tsconfig.json")
	output := fs.String("o", "", "output file path (default: stdout)")
	prune := fs.Bool("prune", false, "remove unreferenced internal types")
	compact := fs.Bool("compact", false, "compact JSON output (no indentation)")
	noSource := fs.Bool("no-source", false, "omit source text from output")
	treeShake := fs.Bool("tree-shake", false, "only emit types/symbols reachable from user code")
	help := fs.Bool("h", false, "show help")

	var flagArgs, posArgs []string
	for i := 0; i < len(args); i++ {
		if args[i] == "-h" || args[i] == "--help" {
			flagArgs = append(flagArgs, args[i])
		} else if args[i] == "-o" || args[i] == "-p" {
			flagArgs = append(flagArgs, args[i])
			if i+1 < len(args) {
				i++
				flagArgs = append(flagArgs, args[i])
			}
		} else if args[i] == "--prune" || args[i] == "--compact" || args[i] == "--no-source" || args[i] == "--tree-shake" {
			flagArgs = append(flagArgs, args[i])
		} else if len(args[i]) > 2 && args[i][0] == '-' && args[i][1] == '-' {
			flagArgs = append(flagArgs, args[i])
		} else if len(args[i]) > 1 && args[i][0] == '-' && args[i][1] != '-' {
			flagArgs = append(flagArgs, args[i])
		} else {
			posArgs = append(posArgs, args[i])
		}
	}

	if err := fs.Parse(flagArgs); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		return 1
	}

	if *help {
		fmt.Println(`Usage: wts emit-ir [options] [files...]

Options:
  -p <path>      Path to tsconfig.json (default: auto-detect)
  -o <path>      Output file path (default: stdout)
  -h             Show this help message
  --prune        Remove unreferenced internal noise types
  --compact      Compact JSON output (no indentation)
  --no-source    Omit source text from output
  --tree-shake   Only emit types/symbols reachable from user code (recommended)

If no files or project are specified, processes all .ts files in the current directory.`)
		return 0
	}

	remainingArgs := posArgs
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
		return 1
	}

	program, err := createProgramFromArgs(*project, remainingArgs, cwd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	ctx := context.Background()
	var diags []*ast.Diagnostic
	for _, sf := range program.GetSourceFiles() {
		if !sf.IsDeclarationFile {
			diags = append(diags, program.GetSyntacticDiagnostics(ctx, sf)...)
		}
	}

	if len(diags) > 0 {
		formatOpts := &diagnosticwriter.FormattingOptions{
			ComparePathsOptions: tspath.ComparePathsOptions{
				CurrentDirectory:          cwd,
				UseCaseSensitiveFileNames: false,
			},
			NewLine: "\n",
		}
		diagnosticwriter.FormatDiagnosticsWithColorAndContext(os.Stderr, diagnosticwriter.FromASTDiagnostics(diags), formatOpts)
		fmt.Fprintf(os.Stderr, "\nFound %d syntax error(s)\n", len(diags))
		return 1
	}

	opts := ir.EmitOptions{
		Prune:     *prune,
		Compact:   *compact,
		NoSource:  *noSource,
		TreeShake: *treeShake,
	}
	emitter := ir.NewEmitter(program, opts)
	_, err = emitter.Emit()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error emitting IR: %v\n", err)
		return 1
	}

	jsonData, err := emitter.Serialize()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error serializing IR: %v\n", err)
		return 1
	}

	outputPath := *output
	if outputPath != "" {
		if err := os.WriteFile(outputPath, jsonData, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing output: %v\n", err)
			return 1
		}
		fmt.Printf("IR written to %s\n", outputPath)
	} else {
		fmt.Println(string(jsonData))
	}

	return 0
}

func createProgramFromArgs(project string, files []string, cwd string) (*compiler.Program, error) {
	var fileNames []string

	if project != "" {
		absPath, err := filepath.Abs(project)
		if err != nil {
			return nil, fmt.Errorf("resolving project path: %w", err)
		}
		fileNames = []string{absPath}
	} else if len(files) > 0 {
		for _, f := range files {
			abs, err := filepath.Abs(f)
			if err != nil {
				abs = f
			}
			fileNames = append(fileNames, abs)
		}
	} else {
		tsconfigPath := filepath.Join(cwd, "tsconfig.json")
		if _, err := os.Stat(tsconfigPath); err == nil {
			fileNames = []string{tsconfigPath}
		} else {
			return nil, fmt.Errorf("no input files or tsconfig.json found")
		}
	}

	compilerOptions := &core.CompilerOptions{
		Target: core.ScriptTargetESNext,
		Module: core.ModuleKindESNext,
		Strict: core.TSTrue,
	}

	compareOpts := tspath.ComparePathsOptions{
		CurrentDirectory:          cwd,
		UseCaseSensitiveFileNames: false,
	}

	parsedConfig := tsoptions.NewParsedCommandLine(compilerOptions, fileNames, compareOpts)

	fs := osvfs.FS()
	host := compiler.NewCachedFSCompilerHost(
		cwd,
		bundled.WrapFS(fs),
		bundled.LibPath(),
		nil,
		nil,
	)

	opts := compiler.ProgramOptions{
		Host:   host,
		Config: parsedConfig,
	}

	return compiler.NewProgram(opts), nil
}
