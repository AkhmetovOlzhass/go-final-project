#!/usr/bin/env bash
set -e

echo "Installing developer tools (macOS/Linux)"

echo "Installing task & air via go generate..."
go generate -tags tools ./tools

GOPATH=$(go env GOPATH)
BIN="$GOPATH/bin"

if [ ! -d "$BIN" ]; then
  echo "ERROR: GOPATH/bin not found: $BIN"
  exit 1
fi

if [[ ":$PATH:" != *":$BIN:"* ]]; then
  echo "Adding Go bin to PATH..."
  echo "export PATH=\$PATH:$BIN" >> ~/.bashrc 2>/dev/null || true
  echo "export PATH=\$PATH:$BIN" >> ~/.zshrc 2>/dev/null || true
else
  echo "Go bin already in PATH."
fi

TASK="$BIN/task"
if [ ! -x "$TASK" ]; then
  echo "ERROR: task not found at $TASK"
  exit 1
fi

echo "Running task init..."
"$TASK" init

echo "Installation complete (macOS/Linux)"
