package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	project          string
	output           string
	help             bool
	prune            bool
	compact          bool
	noSource         bool
	treeShake        bool
	runtimeDescriptor string
}

func runEmitIR(args []string) int {
	fs := flag.NewFlagSet("emit-ir", flag.ExitOnError)
	project := fs.String("p", "", "path to tsconfig.json")
	output := fs.String("o", "", "output file path (default: stdout)")
	prune := fs.Bool("prune", false, "remove unreferenced internal types")
	compact := fs.Bool("compact", false, "compact JSON output (no indentation)")
	noSource := fs.Bool("no-source", false, "omit source text from output")
	treeShake := fs.Bool("tree-shake", false, "only emit types/symbols reachable from user code")
	runtimeDescriptor := fs.String("runtime-descriptor", "", "path to a RuntimeDescriptor JSON file; injects it as the WIR 'runtime' field so the backend honors frontend-authored runtime bindings (no post-injection)")
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
		} else if args[i] == "--runtime-descriptor" {
			flagArgs = append(flagArgs, args[i])
			if i+1 < len(args) {
				i++
				flagArgs = append(flagArgs, args[i])
			}
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
  --runtime-descriptor <path>  Inject a RuntimeDescriptor JSON as the WIR 'runtime' field

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

	// Load a frontend-authored runtime descriptor (noStdRoadMap N.3 / miss-wts-002).
	var runtimeDesc *ir.RuntimeDescriptor
	if *runtimeDescriptor != "" {
		rdData, err := os.ReadFile(*runtimeDescriptor)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading runtime descriptor %q: %v\n", *runtimeDescriptor, err)
			return 1
		}
		runtimeDesc = &ir.RuntimeDescriptor{}
		if err := json.Unmarshal(rdData, runtimeDesc); err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing runtime descriptor %q: %v\n", *runtimeDescriptor, err)
			return 1
		}
		opts.Runtime = runtimeDesc
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
	// Create the host up-front: it is needed both to read the tsconfig text and to
	// act as the ParseConfigHost that resolves `files`/`include` globs.
	fs := osvfs.FS()
	host := compiler.NewCachedFSCompilerHost(
		cwd,
		bundled.WrapFS(fs),
		bundled.LibPath(),
		nil,
		nil,
	)

	baseOptions := &core.CompilerOptions{
		Target: core.ScriptTargetESNext,
		Module: core.ModuleKindESNext,
		Strict: core.TSTrue,
	}

	compareOpts := tspath.ComparePathsOptions{
		CurrentDirectory:          cwd,
		UseCaseSensitiveFileNames: host.FS().UseCaseSensitiveFileNames(),
	}

	var parsedConfig *tsoptions.ParsedCommandLine
	var err error

	switch {
	case project != "":
		configPath, errResolve := resolveTsConfigPath(project, cwd)
		if errResolve != nil {
			return nil, errResolve
		}
		parsedConfig, err = parseTsConfig(configPath, host, cwd, baseOptions)
		if err != nil {
			return nil, err
		}
	case len(files) > 0:
		var fileNames []string
		for _, f := range files {
			abs, absErr := filepath.Abs(f)
			if absErr != nil {
				abs = f
			}
			fileNames = append(fileNames, abs)
		}
		parsedConfig = tsoptions.NewParsedCommandLine(baseOptions, fileNames, compareOpts)
	default:
		// No input files or project specified; default to tsconfig.json in cwd if present.
		tsconfigPath := filepath.Join(cwd, "tsconfig.json")
		if _, statErr := os.Stat(tsconfigPath); statErr == nil {
			parsedConfig, err = parseTsConfig(tsconfigPath, host, cwd, baseOptions)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("no input files or tsconfig.json found in %s", cwd)
		}
	}

	opts := compiler.ProgramOptions{
		Host:   host,
		Config: parsedConfig,
	}

	return compiler.NewProgram(opts), nil
}

// resolveTsConfigPath resolves a `-p` project argument to a tsconfig.json path.
// If the argument points at a directory, it is joined with "tsconfig.json".
func resolveTsConfigPath(project, cwd string) (string, error) {
	abs, err := filepath.Abs(project)
	if err != nil {
		return "", fmt.Errorf("invalid project path %q: %w", project, err)
	}
	if info, statErr := os.Stat(abs); statErr == nil && info.IsDir() {
		abs = filepath.Join(abs, "tsconfig.json")
	}
	// ParseConfigFileTextToJson requires a normalized (forward-slash) absolute path.
	return tspath.GetNormalizedAbsolutePath(abs, cwd), nil
}

// parseTsConfig reads and parses a tsconfig.json, returning a ParsedCommandLine
// whose FileNames (files/include/exclude) and CompilerOptions (e.g. noLib, lib)
// are genuinely honored by the program. Previously `-p <path>` merely stuffed
// the config path into the root file list, so the config was compiled as a source
// file and its options (including noLib) were ignored.
func parseTsConfig(configPath string, host compiler.CompilerHost, cwd string, baseOptions *core.CompilerOptions) (*tsoptions.ParsedCommandLine, error) {
	text, ok := host.FS().ReadFile(configPath)
	if !ok {
		return nil, fmt.Errorf("cannot read tsconfig %q", configPath)
	}
	jsonVal, jsonErrs := tsoptions.ParseConfigFileTextToJson(
		configPath,
		tspath.ToPath(configPath, cwd, host.FS().UseCaseSensitiveFileNames()),
		text,
	)
	if len(jsonErrs) > 0 {
		return nil, fmt.Errorf("failed to parse tsconfig %q: %s", configPath, formatDiagnostics(jsonErrs))
	}
	basePath := tspath.GetDirectoryPath(configPath)
	parsedConfig := tsoptions.ParseJsonConfigFileContent(jsonVal, host, basePath, baseOptions, configPath, nil, nil, nil)
	if diags := parsedConfig.GetConfigFileParsingDiagnostics(); len(diags) > 0 {
		return nil, fmt.Errorf("invalid tsconfig %q: %s", configPath, formatDiagnostics(diags))
	}
	return parsedConfig, nil
}

func formatDiagnostics(diags []*ast.Diagnostic) string {
	msgs := make([]string, 0, len(diags))
	for _, d := range diags {
		msgs = append(msgs, d.String())
	}
	return strings.Join(msgs, "; ")
}
