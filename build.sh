#!/bin/bash

echo "🏗️ Building Venotes..."

# Build frontend
echo "📦 Building frontend..."
cd frontend
npm install
npm run build
cd ..

# Build Go backend with embedded frontend
echo "🔧 Building Go backend with embedded frontend..."
go mod tidy
go build -o venotes main.go

echo "✅ Build complete!"
echo "🚀 Run with: ./venotes"
echo "📦 Single standalone executable with embedded frontend!"
echo "💾 No external files needed - everything is embedded!"
