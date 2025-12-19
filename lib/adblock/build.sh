#!/bin/bash
# Build script for libadblock on Linux/macOS
echo "Building libadblock Rust library..."

cd "$(dirname "$0")"
cargo build --release

if [ $? -eq 0 ]; then
    echo "libadblock build successful!"
else
    echo "libadblock build failed."
    exit 1
fi
