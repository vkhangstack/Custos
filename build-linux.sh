#!/bin/bash

# Exit on error
set -e

echo "Building libadblock Rust library..."
./lib/adblock/build.sh

# Build the app using Wails
wails build -platform linux/amd64

echo "Packaging .deb..."
nfpm pkg --packager deb --config nfpm.yaml --target .

echo "Packaging .rpm..."
nfpm pkg --packager rpm --config nfpm.yaml --target .

echo "Build complete!"
ls -lh *.deb *.rpm
