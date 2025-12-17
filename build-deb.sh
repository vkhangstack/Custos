#!/bin/bash

# Build the app
wails build

# Package the app
nfpm pkg --packager deb --config nfpm.yaml --target .