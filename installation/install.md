# Installation Guide

To use EchoTune, you need to install its required dependencies and the `etune` binary.

## 1. Install Dependencies (`mpv` and `yt-dlp`)

EchoTune relies on `mpv` to play audio streams and `yt-dlp` to search and fetch audio from YouTube. You must install these first. 

Below are the easiest commands to copy and paste into your terminal based on your operating system. We use standard package managers for `mpv` and `curl` for fetching the latest official `yt-dlp` binary.

### Linux

**Ubuntu / Debian / Mint:**
```bash
# Install mpv
sudo apt-get update && sudo apt-get install -y mpv

# Download the latest yt-dlp using curl
sudo curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /usr/local/bin/yt-dlp
sudo chmod a+rx /usr/local/bin/yt-dlp
```

**Fedora / RHEL:**
```bash
# Install mpv
sudo dnf install -y mpv

# Download the latest yt-dlp using curl
sudo curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp -o /usr/local/bin/yt-dlp
sudo chmod a+rx /usr/local/bin/yt-dlp
```

**Arch Linux / Manjaro:**
```bash
# Install both via pacman
sudo pacman -Syu --needed mpv yt-dlp
```

### macOS

Ensure you have [Homebrew](https://brew.sh/) installed, then run:
```bash
brew install mpv yt-dlp
```
> **Apple Silicon (M1/M2/M3) users:** Homebrew installs to `/opt/homebrew/bin`. Make sure this directory is in your `PATH` (usually added automatically by the Homebrew installer, but worth verifying if `etune` can't find the dependencies).

### Windows

**Option A: Scoop (Recommended)**

Ensure you have [Scoop](https://scoop.sh/) installed, then run:
```powershell
# Install yt-dlp (available in the default main bucket)
scoop install yt-dlp

# Install mpv (requires the extras bucket)
scoop bucket add extras
scoop install extras/mpv
```

**Option B: Winget (Built-in)**

If you prefer not to use Scoop, you can use Windows' built-in package manager:
```powershell
winget install yt-dlp
winget install mpv
```

**Option C: Manual Download**

Download the standalone binaries and add them to your System PATH:
- **yt-dlp:** Download `yt-dlp.exe` from the [latest release](https://github.com/yt-dlp/yt-dlp/releases/latest)
- **mpv:** Download the latest Windows build from [mpv-winbuild](https://github.com/shinchiro/mpv-winbuild-cmake/releases)

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
