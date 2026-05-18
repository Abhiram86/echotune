# This script installs mpv and yt-dlp on Windows using Scoop.
# If you don't have Scoop installed, get it from: https://scoop.sh/

if (!(Get-Command scoop -ErrorAction SilentlyContinue)) {
    Write-Host "Scoop is not installed. Please install it first from https://scoop.sh/" -ForegroundColor Red
    exit 1
}

Write-Host "Installing mpv and yt-dlp via Scoop..." -ForegroundColor Cyan
scoop install mpv yt-dlp

Write-Host "Dependencies installed successfully!" -ForegroundColor Green
