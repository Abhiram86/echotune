#!/bin/bash

# This script installs mpv and yt-dlp on macOS using Homebrew.

if ! command -v brew &> /dev/null
then
    echo "Homebrew not found. Please install Homebrew first: https://brew.sh/"
    exit 1
fi

echo "Installing mpv and yt-dlp via Homebrew..."
brew install mpv yt-dlp

echo "Dependencies installed successfully!"
