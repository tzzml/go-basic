package main

import (
	"flag"
	"fmt"
	"os"
	"zork-basic/internal/bytecode"
	"zork-basic/internal/vm"
)

func main() {
	var dump bool
	flag.BoolVar(&dump, "d", false, "Disassemble bytecode file instead of running")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <bytecode_file>\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(1)
	}

	filename := flag.Arg(0)
	f, err := os.Open(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	chunk, err := bytecode.ReadChunk(f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading bytecode: %v\n", err)
		os.Exit(1)
	}

	if dump {
		fmt.Printf("Disassembly of %s:\n", filename)
		fmt.Println(chunk.Disassemble(filename))
		return
	}

	// Create and run VM
	machine := vm.New(chunk)
	if err := machine.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Runtime error: %v\n", err)
		os.Exit(1)
	}
}
