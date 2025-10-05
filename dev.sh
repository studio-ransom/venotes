#!/bin/bash

echo "🚀 Starting Svelte Notes in development mode..."

# Check if frontend is built
if [ ! -d "frontend/dist" ]; then
    echo "📦 Building frontend first..."
    cd frontend
    npm install
    npm run build
    cd ..
fi

# Start the Go server
echo "🔧 Starting Go server..."
go run main.go
