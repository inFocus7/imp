#!/bin/bash

dependencies=(
    "shasum"
    "upx"
)

missing_deps=()

for dep in "${dependencies[@]}"; do
    if ! command -v "$dep" &> /dev/null; then
        missing_deps+=("$dep")
    fi
done

if [ ${#missing_deps[@]} -gt 0 ]; then
    echo "Missing dependencies: ${missing_deps[@]}"
    exit 1
fi

echo "All dependencies are installed"
