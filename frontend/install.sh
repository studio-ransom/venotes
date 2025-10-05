#!/bin/bash

echo "🗒️ Setting up Svelte Notes..."

# Check if npm is installed
if ! command -v npm &> /dev/null; then
    echo "❌ npm is not installed. Please install Node.js and npm first."
    exit 1
fi

# Install dependencies
echo "📦 Installing dependencies..."
npm install

if [ $? -eq 0 ]; then
    echo "✅ Dependencies installed successfully!"
    echo ""
    echo "🚀 To start development:"
    echo "   npm run dev"
    echo ""
    echo "🏗️ To build for production:"
    echo "   npm run build"
    echo ""
    echo "📝 Then open dist/index.html in your browser!"
else
    echo "❌ Failed to install dependencies"
    exit 1
fi
