// zork-basic - BASIC 语言解释器
// 使用 Go 语言编写，基于 PEG 解析器生成器
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"zork-basic/internal/ast"
	"zork-basic/internal/bytecode"
	"zork-basic/internal/compiler"
	"zork-basic/internal/parser"
	"zork-basic/internal/repl"
	"zork-basic/internal/vm"
)

const (
	// Version 版本信息
	Version = "1.1.0"
)

func main() {
	// 定义命令行参数
	interactive := flag.Bool("i", false, "Interactive mode")
	interactiveLong := flag.Bool("interactive", false, "Interactive mode")
	version := flag.Bool("v", false, "Show version")
	versionLong := flag.Bool("version", false, "Show version")
	help := flag.Bool("h", false, "Show help")
	helpLong := flag.Bool("help", false, "Show help")
	modePtr := flag.String("mode", "vm", "Execution mode: ast or vm (for .bas files)")
	outputFile := flag.String("o", "", "Compile to bytecode file (.zbc)")
	disassemble := flag.Bool("d", false, "Disassemble bytecode")

	flag.Parse()

	// 显示帮助
	if *help || *helpLong {
		printHelp()
		return
	}

	// 显示版本
	if *version || *versionLong {
		fmt.Printf("Zork BASIC v%s\n", Version)
		fmt.Println("A high-performance BASIC interpreter and compiler written in Go")
		return
	}

	// 确定运行模式
	isInteractive := *interactive || *interactiveLong
	mode := *modePtr

	// 检查参数
	args := flag.Args()
	if len(args) == 0 && !isInteractive {
		// 无参数且不是交互模式，默认进入交互模式
		isInteractive = true
	}

	if isInteractive {
		if *outputFile != "" {
			fmt.Println("Warning: -o flag is ignored in interactive mode")
		}
		if *disassemble {
			fmt.Println("Warning: -d flag is ignored in interactive mode")
		}
		repl.Run(Version, mode)
	} else {
		// 文件模式
		filename := args[0]

		if *outputFile != "" {
			compileFileToBytecode(filename, *outputFile)
			return
		}

		if *disassemble {
			disassembleFile(filename)
			return
		}

		runFileUnified(filename, mode)
	}
}

// detectFileType 探测文件类型：返回 "bytecode", "source", 或 "unknown"
func detectFileType(filename string) (string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	header := make([]byte, 3)
	n, err := f.Read(header)
	if err != nil && err != io.EOF {
		// handle real error if needed
	}
	if n == 3 && string(header) == "ZBC" {
		return "bytecode", nil
	}

	return "source", nil
}

// disassembleFile 反汇编执行文件
func disassembleFile(filename string) {
	fileType, err := detectFileType(filename)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	var chunk *bytecode.Chunk
	if fileType == "bytecode" {
		f, err := os.Open(filename)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		chunk, err = bytecode.ReadChunk(f)
		if err != nil {
			fmt.Printf("Error reading bytecode: %v\n", err)
			os.Exit(1)
		}
	} else {
		// Compile source first
		data, err := os.ReadFile(filename)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		parsedAST, err := parser.Parse(filename, data)
		if err != nil {
			fmt.Printf("Parse error: %v\n", err)
			os.Exit(1)
		}
		prog := parsedAST.(*ast.Program)
		comp := compiler.New()
		chunk, err = comp.Compile(prog)
		if err != nil {
			fmt.Printf("Compilation error: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Print(chunk.Disassemble(filename))
}

// runFileUnified 统一运行文件（自动识别类型）
func runFileUnified(filename string, mode string) {
	fileType, err := detectFileType(filename)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	if fileType == "bytecode" {
		f, err := os.Open(filename)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
		chunk, err := bytecode.ReadChunk(f)
		if err != nil {
			fmt.Printf("Error reading bytecode: %v\n", err)
			os.Exit(1)
		}
		v := vm.New(chunk)
		if err := v.Run(); err != nil {
			fmt.Printf("Runtime error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("\nProgram complete.")
	} else {
		// 源代码模式
		data, err := os.ReadFile(filename)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		repl.ExecuteProgram(string(data), filename, mode)
	}
}

// compileFileToBytecode 编译文件为字节码并保存
func compileFileToBytecode(inputFile, outputFile string) {
	fmt.Printf("Compiling %s to %s...\n", inputFile, outputFile)
	data, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	parsedAST, err := parser.Parse(inputFile, data)
	if err != nil {
		fmt.Printf("Parse error: %v\n", err)
		os.Exit(1)
	}

	prog, ok := parsedAST.(*ast.Program)
	if !ok {
		fmt.Println("Error: not a valid BASIC program")
		os.Exit(1)
	}

	comp := compiler.New()
	chunk, err := comp.Compile(prog)
	if err != nil {
		fmt.Printf("Compilation error: %v\n", err)
		os.Exit(1)
	}

	f, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	if err := chunk.Write(f); err != nil {
		fmt.Printf("Error writing bytecode: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Compilation successful.")
}

// printHelp 打印帮助信息
func printHelp() {
	fmt.Println("Zork BASIC - A high-performance BASIC interpreter and compiler")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  zb [options] <file.bas>     Run a BASIC source file")
	fmt.Println("  zb [options] <file.zbc>     Run a compiled bytecode file")
	fmt.Println("  zb -i                       Start interactive mode (REPL)")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  -i, --interactive    Run in interactive mode")
	fmt.Println("  -v, --version        Show version information")
	fmt.Println("  -h, --help           Show this help message")
	fmt.Println("  -mode <ast|vm>       Execution mode for source files (default: vm)")
	fmt.Println("  -o <file.zbc>        Compile source to a bytecode file")
	fmt.Println("  -d                   Disassemble bytecode (supports .bas and .zbc)")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  zb hello.bas                Run program using VM")
	fmt.Println("  zb -mode ast hello.bas      Run program using AST interpreter")
	fmt.Println("  zb hello.zbc                Run compiled bytecode")
	fmt.Println("  zb -o hello.zbc hello.bas   Compile to bytecode")
	fmt.Println("  zb -d hello.bas             View bytecode for source file")
}
