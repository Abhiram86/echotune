# Installation Guide

To use EchoTune, you need to install its required dependencies and the `etune` binary.

## 1. Install Dependencies (`mpv` and `yt-dlp`)

EchoTune relies on `mpv` to play audio streams and `yt-dlp` to search and fetch audio from YouTube. You must install these first based on your operating system. We have provided automated scripts for convenience.

### Linux (Debian/Ubuntu)
Run the provided Linux script to install `mpv` via `apt` and the latest `yt-dlp` via `curl`:
```bash
bash installation/deps_linux_install.sh
```

### macOS
Run the provided macOS script to install the dependencies via Homebrew:
```bash
bash installation/deps_mac_install.sh
```

### Windows
Run the provided PowerShell script to install the dependencies via Scoop:
```powershell
.\installation\deps_windows_install.ps1
```

---

## 2. Install EchoTune (`etune`)

Once the dependencies are installed, you have three options to install the EchoTune binary.

### Option A: Download Pre-built Binaries (Recommended)
1. Go to the [Releases](../../releases) page of this GitHub repository.
2. Download the appropriate `.zip` or `.tar.gz` file for your operating system (Linux, macOS, or Windows).
3. Extract the archive and place the `etune` binary in a directory that is in your system's `PATH` (e.g., `/usr/local/bin` for Linux/Mac, or any configured folder for Windows).
4. Run `etune` in your terminal!

### Option B: Install via `go install`
If you have Go installed on your system, you can easily install the latest version directly:
```bash
go install github.com/Abhiram86/echotune@latest
```
*(Note: If installed this way, the binary will default to the name `echotune`, though you can rename it to `etune` if preferred).*

### Option C: Build from Source
If you want to clone the repository and build it manually:
```bash
git clone https://github.com/Abhiram86/echotune.git
cd echotune
go build -o etune main.go
```
Then run `./etune`!
