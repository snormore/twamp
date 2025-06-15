#!/usr/bin/env bash

targets=(
  FuzzUnmarshalTestPacket
  FuzzRoundTripPacket
  FuzzNTPConversion
)

for target in "${targets[@]}"; do
  echo "=== Running $target ==="
  go test -run=^$ -fuzz="$target" -fuzztime=1m ./...
  echo
done

