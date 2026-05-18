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

### Windows

Ensure you have [Scoop](https://scoop.sh/) installed, then run:
```powershell
scoop install mpv yt-dlp
```
*(Alternatively, you can download `yt-dlp.exe` and `mpv.exe` manually and add them to your System PATH).*

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
