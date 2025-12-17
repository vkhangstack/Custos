#!/bin/bash

# Build the app
wails build

# Package the app
nfpm pkg --packager rpm --config nfpm.yaml --target .