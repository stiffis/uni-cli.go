#!/bin/bash
# Quick script to build and run UniCLI

echo "🔨 Building UniCLI..."
go build -o app ./cmd/unicli

if [ $? -eq 0 ]; then
    echo "✅ Build successful!"
    echo "🚀 Starting UniCLI..."
    echo ""
    ./app
else
    echo "❌ Build failed!"
    exit 1
fi
