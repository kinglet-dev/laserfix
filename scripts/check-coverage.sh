#!/usr/bin/env bash
# Fails unless every statement is covered. Usage: scripts/check-coverage.sh coverage.out
set -euo pipefail
profile="${1:-coverage.out}"
total=$(go tool cover -func="$profile" | awk '/^total:/ {print $NF}')
echo "coverage: $total of statements"
if [[ "$total" != "100.0%" ]]; then
  echo "Coverage must be 100%. Uncovered functions:"
  go tool cover -func="$profile" | awk '$NF != "100.0%" && !/^total:/'
  exit 1
fi
