#!/usr/bin/env bash
set -e

# Generate shell completions for ason
echo "Generating shell completions..."

# Create completions directory
mkdir -p completions

# Generate completions
go run . completion bash > completions/ason.bash
go run . completion zsh > completions/_ason
go run . completion fish > completions/ason.fish

echo "✅ Completions generated in completions/"
ls -lh completions/
