#!/bin/bash

# This script installs mpv and yt-dlp on Debian/Ubuntu-based Linux systems.
# yt-dlp is installed directly from its official repository to ensure it is the latest version.

echo "Installing mpv..."
sudo apt-get update
sudo apt-get install -y mpv

echo "Installing yt-dlp..."
sudo curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /usr/local/bin/yt-dlp
sudo chmod a+rx /usr/local/bin/yt-dlp

echo "Dependencies installed successfully!"
