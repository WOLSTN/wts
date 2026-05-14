package main

import (
	"fmt"
	"os"

	"github.com/wolstn/wts/internal/core"
)

func main() {
	os.Exit(runMain())
}

func runMain() int {
	core.ApplyDebugStackLimit()
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
  wts emit-ir main.ts -o main.wir
  wts emit-ir --project tsconfig.json -o program.wir`)
}

func printVersion() {
	fmt.Println("wts version 0.1.0 (Wolstn TypeScript Frontend)")
}

func runCheck(args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Error: no input files specified")
		return 1
	}

	fmt.Println("Type checking not yet implemented")
	fmt.Printf("Would check: %v\n", args)
	return 0
}

func runEmitIR(args []string) int {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Error: no input files specified")
		return 1
	}

	fmt.Println("IR emission not yet implemented")
	fmt.Printf("Would emit IR for: %v\n", args)
	return 0
}
