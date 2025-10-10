#!/bin/bash
# Quick script to build and run UniCLI

echo "🔨 Building UniCLI..."
go build -o unicli ./cmd/unicli

if [ $? -eq 0 ]; then
    echo "✅ Build successful!"
    echo "🚀 Starting UniCLI..."
    echo ""
    ./unicli
else
    echo "❌ Build failed!"
    exit 1
fi
