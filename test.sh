#!/bin/bash

# Cleanup function
cleanup() {
  echo "Cleaning up..."
  rm -f /tmp/tempfile
}

# Set up traps
trap 'cleanup; exit' EXIT
trap 'echo "Interrupted!"; cleanup; exit 1' INT

# Create a temporary file
touch /tmp/tempfile

echo "Script running. Press Ctrl+C to interrupt."

# Simulate a long-running process
sleep 5

echo "Script completed."
