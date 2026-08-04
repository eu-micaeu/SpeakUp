#!/bin/sh
set -e

# Start Ollama server in background
ollama serve &
SERVER_PID=$!

MODEL="${OLLAMA_MODEL:-llama3.2:latest}"

echo "Waiting for Ollama service to start..."
until ollama list > /dev/null 2>&1; do
  sleep 1
done

echo "Ollama server is running."
echo "Checking model status for: ${MODEL}..."

# Check if model already exists
if ! ollama list | grep -q "${MODEL%%:*}"; then
  echo "Model '${MODEL}' not found locally. Pulling model automatically..."
  ollama pull "${MODEL}"
  echo "Model '${MODEL}' pulled successfully!"
else
  echo "Model '${MODEL}' is already available."
fi

# Wait for server process
wait $SERVER_PID
