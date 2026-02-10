#!/bin/bash
# Build script for zork-basic suite (zb, zbc, zvm)

echo "Building zork-basic tools..."

# Create bin directory if it doesn't exist
mkdir -p bin

# Stripped build for better size
LDFLAGS="-s -w"

echo "1. Building zb (Interpreter/REPL)..."
go build -ldflags="$LDFLAGS" -trimpath -o bin/zb ./cmd/zb
ls -lh bin/zb

echo ""
echo "2. Building zbc (Bytecode Compiler)..."
go build -ldflags="$LDFLAGS" -trimpath -o bin/zbc ./cmd/zbc
ls -lh bin/zbc

echo ""
echo "3. Building zvm (Bytecode VM Runner)..."
go build -ldflags="$LDFLAGS" -trimpath -o bin/zvm ./cmd/zvm
ls -lh bin/zvm

echo ""
echo "Build complete!"
echo "Binaries are located in the 'bin' directory: zb, zbc, zvm"
