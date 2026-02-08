#!/bin/bash
# Build script for zork-basic with size optimizations

echo "Building zork-basic with optimizations..."

# Standard build (with debug symbols)
echo "1. Standard build..."
go build -o zork-basic ./cmd/zork-basic
ls -lh zork-basic

# Optimized build (stripped)
echo ""
echo "2. Optimized build (stripped)..."
go build -ldflags="-s -w" -trimpath -o zork-basic ./cmd/zork-basic
ls -lh zork-basic

# Show size comparison
echo ""
echo "Build complete!"
echo ""
echo "To further reduce size (optional):"
echo "  - Install UPX: brew install upx (macOS) or apt install upx (Linux)"
echo "  - Compress: upx --best --lzma zork-basic"
echo "  - This can reduce size to ~600-800KB but adds decompression overhead at startup"
