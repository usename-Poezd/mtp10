#!/bin/bash
set -e

GO_URL="${GO_URL:-http://localhost:8080}"
PYTHON_URL="${PYTHON_URL:-http://localhost:8000}"
THREADS=4
CONNECTIONS=100
DURATION=30s

echo "=== Waiting for services ==="
for url in "$GO_URL/todos" "$PYTHON_URL/todos"; do
  echo -n "Waiting for $url..."
  for i in $(seq 1 30); do
    if curl -sf "$url" > /dev/null 2>&1; then
      echo " ready"
      break
    fi
    sleep 1
    if [ "$i" -eq 30 ]; then
      echo " TIMEOUT"
      exit 1
    fi
  done
done

echo ""
echo "=== Seeding 10 TODO items in Go service ==="
for i in $(seq 1 10); do
  curl -sf -X POST "$GO_URL/todos" \
    -H "Content-Type: application/json" \
    -d "{\"title\": \"Benchmark item $i\"}" > /dev/null
done
echo "Seeded 10 items"

echo ""
echo "=== Benchmarking Go service (Gin) ==="
GO_RESULT=$(wrk -t$THREADS -c$CONNECTIONS -d$DURATION "$GO_URL/todos")
echo "$GO_RESULT"

echo ""
echo "=== Benchmarking Python service (FastAPI) ==="
PYTHON_RESULT=$(wrk -t$THREADS -c$CONNECTIONS -d$DURATION "$PYTHON_URL/todos")
echo "$PYTHON_RESULT"

echo ""
echo "=== Results saved ==="
