#!/bin/bash
# Build script for zork-basic suite (zb, zbc, zvm)

echo "Building zork-basic tools..."

# Stripped build for better size
LDFLAGS="-s -w"

echo "1. Building zb (Interpreter/REPL)..."
go build -ldflags="$LDFLAGS" -trimpath -o zb ./cmd/zork-basic
ls -lh zb

echo ""
echo "2. Building zbc (Bytecode Compiler)..."
go build -ldflags="$LDFLAGS" -trimpath -o zbc ./cmd/zbc
ls -lh zbc

echo ""
echo "3. Building zvm (Bytecode VM Runner)..."
go build -ldflags="$LDFLAGS" -trimpath -o zvm ./cmd/zvm
ls -lh zvm

echo ""
echo "Build complete!"
echo "Binaries: zb, zbc, zvm"
