package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/wolstn/wts/internal/compiler"
	"github.com/wolstn/wts/internal/core"
	"github.com/wolstn/wts/internal/diagnostics"
	"github.com/wolstn/wts/internal/diagnosticwriter"
	"github.com/wolstn/wts/internal/ir"
	"github.com/wolstn/wts/internal/tsoptions"
	"github.com/wolstn/wts/internal/tspath"
	"github.com/wolstn/wts/internal/vfs"
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
  wts check --project tsconfig.json
  wts emit-ir main.ts -o main.wir
  wts emit-ir --project tsconfig.json -o program.wir`)
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

	if err := fs.Parse(args); err != nil {
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

	remainingArgs := fs.Args()
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
		return 1
	}

	configPath := *project
	if configPath == "" && len(remainingArgs) == 0 {
		configPath = findTsConfig(cwd)
	}

	var parsedConfig *tsoptions.ParsedCommandLine
	if configPath != "" {
		parsedConfig, err = parseTsConfig(configPath, cwd)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing tsconfig: %v\n", err)
			return 1
		}
	} else if len(remainingArgs) > 0 {
		parsedConfig = createConfigFromFiles(remainingArgs, cwd)
	} else {
		fmt.Fprintln(os.Stderr, "No input files or tsconfig.json found")
		return 1
	}

	program, err := createProgram(parsedConfig, cwd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating program: %v\n", err)
		return 1
	}

	diags := program.GetSyntacticDiagnostics(context.Background())
	diags = append(diags, program.GetSemanticDiagnostics(context.Background())...)
	diags = append(diags, program.GetGlobalDiagnostics(context.Background())...)

	if len(diags) > 0 {
		writer := diagnosticwriter.New(os.Stderr, program.GetSourceFiles(), false)
		for _, d := range diags {
			writer.WriteDiagnostic(d)
		}
		fmt.Fprintf(os.Stderr, "\nFound %d error(s)\n", len(diags))
		return 1
	}

	fmt.Println("No errors found.")
	return 0
}

type emitIRFlags struct {
	project string
	output  string
	help    bool
}

func runEmitIR(args []string) int {
	fs := flag.NewFlagSet("emit-ir", flag.ExitOnError)
	project := fs.String("p", "", "path to tsconfig.json")
	output := fs.String("o", "", "output file path (default: stdout)")
	help := fs.Bool("h", false, "show help")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		return 1
	}

	if *help {
		fmt.Println(`Usage: wts emit-ir [options] [files...]

Options:
  -p <path>   Path to tsconfig.json (default: auto-detect)
  -o <path>   Output file path (default: stdout)
  -h          Show this help message

If no files or project are specified, processes all .ts files in the current directory.`)
		return 0
	}

	remainingArgs := fs.Args()
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current directory: %v\n", err)
		return 1
	}

	configPath := *project
	if configPath == "" && len(remainingArgs) == 0 {
		configPath = findTsConfig(cwd)
	}

	var parsedConfig *tsoptions.ParsedCommandLine
	if configPath != "" {
		parsedConfig, err = parseTsConfig(configPath, cwd)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing tsconfig: %v\n", err)
			return 1
		}
	} else if len(remainingArgs) > 0 {
		parsedConfig = createConfigFromFiles(remainingArgs, cwd)
	} else {
		fmt.Fprintln(os.Stderr, "No input files or tsconfig.json found")
		return 1
	}

	program, err := createProgram(parsedConfig, cwd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating program: %v\n", err)
		return 1
	}

	diags := program.GetSyntacticDiagnostics(context.Background())
	if len(diags) > 0 {
		writer := diagnosticwriter.New(os.Stderr, program.GetSourceFiles(), false)
		for _, d := range diags {
			writer.WriteDiagnostic(d)
		}
		fmt.Fprintf(os.Stderr, "\nFound %d syntax error(s)\n", len(diags))
		return 1
	}

	emitter := ir.NewEmitter(program)
	irProgram, err := emitter.Emit()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error emitting IR: %v\n", err)
		return 1
	}

	jsonData, err := irProgram.ToJSON()
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

func findTsConfig(dir string) string {
	tsconfigPath := filepath.Join(dir, "tsconfig.json")
	if _, err := os.Stat(tsconfigPath); err == nil {
		return tsconfigPath
	}
	return ""
}

func parseTsConfig(configPath string, cwd string) (*tsoptions.ParsedCommandLine, error) {
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, fmt.Errorf("resolving config path: %w", err)
	}

	fs := osvfs.FS()
	parsedConfig, diags := tsoptions.ParseTsconfigFile(
		context.Background(),
		fs,
		tspath.NormalizePath(absPath),
		tsoptions.NewExtendedConfigCache(),
	)

	if len(diags) > 0 {
		for _, d := range diags {
			fmt.Fprintf(os.Stderr, "Config error: %s\n", d.Message())
		}
		if parsedConfig == nil {
			return nil, fmt.Errorf("failed to parse tsconfig")
		}
	}

	return parsedConfig, nil
}

func createConfigFromFiles(files []string, cwd string) *tsoptions.ParsedCommandLine {
	compilerOptions := core.NewCompilerOptions()
	compilerOptions.SetTarget(core.ScriptTargetESNext)
	compilerOptions.SetModule(core.ModuleKindESNext)
	compilerOptions.SetStrict(true)

	var normalizedFiles []string
	for _, f := range files {
		abs, err := filepath.Abs(f)
		if err != nil {
			abs = f
		}
		normalizedFiles = append(normalizedFiles, tspath.NormalizePath(abs))
	}

	return tsoptions.NewParsedCommandLine(
		compilerOptions,
		normalizedFiles,
		nil,
	)
}

func createProgram(parsedConfig *tsoptions.ParsedCommandLine, cwd string) (*compiler.Program, error) {
	fs := osvfs.FS()
	host := compiler.NewCachedFSCompilerHost(
		cwd,
		fs,
		"",
		tsoptions.NewExtendedConfigCache(),
		func(msg *diagnostics.Message, args ...any) {},
	)

	opts := compiler.ProgramOptions{
		Host:   host,
		Config: parsedConfig,
	}

	return compiler.NewProgram(opts), nil
}

func parseCommandLine(args []string) (files []string, options map[string]string) {
	options = make(map[string]string)
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "-") || strings.HasPrefix(arg, "--") {
			key := strings.TrimLeft(arg, "-")
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				options[key] = args[i+1]
				i++
			} else {
				options[key] = "true"
			}
		} else {
			files = append(files, arg)
		}
	}
	return
}
