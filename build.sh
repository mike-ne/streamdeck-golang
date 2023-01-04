#!/usr/bin/env sh
# Note, this requires that the Elgato StreamDeck DistributionTool is
# installed on your computer. This tool can be downloaded here:
# https://developer.elgato.com/documentation/stream-deck/sdk/packaging/

# Clean up
rm releases/*
rm src/com.mnelis.godonothing.sdPlugin/bin/*

# Build our golang executable
startDir="$(pwd)"
cd src/com.mnelis.godonothing.sdPlugin
go build -o bin/main
cd "$startDir"

# Create a StreamDeck Plugin distribution
DistributionTool -b -i src/com.mnelis.godonothing.sdPlugin -o ./releases
