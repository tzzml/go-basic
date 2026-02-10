#!/bin/bash
# Build script for zork-basic suite (zb, zbc, zvm)

echo "Building zork-basic tools..."

# Create bin directory if it doesn't exist
mkdir -p bin

# Stripped build for better size
LDFLAGS="-s -w"

echo "1. Building Zork BASIC (zb)..."
go build -ldflags="$LDFLAGS" -trimpath -o bin/zb ./cmd
ls -lh bin/zb

echo ""
echo "Build complete!"
echo "The 'zb' binary is located in the 'bin' directory."
