# EchoTune

EchoTune is a fast, terminal-based music player and downloader. It allows you to search for songs on YouTube, play them directly in your terminal without any video overhead, and even download them for offline playback. Built with Go, it provides a clean, distraction-free music experience right from your command line.

## Prerequisites & Installation

Currently, EchoTune is best supported on **Linux**. It relies on the following external dependencies to fetch and play audio:

1. **`yt-dlp`**: Required for searching and downloading audio from YouTube.
2. **`mpv`**: Required for playing the audio streams in the background.

Ensure both are installed and available in your system's `PATH`.

To install EchoTune, you need Go installed on your system. Run the following command:

```bash
go install github.com/Abhiram86/echotune@latest
```

Alternatively, you can clone the repository and build it manually:
```bash
git clone https://github.com/Abhiram86/echotune.git
cd echotune
go build -o et main.go
```

## Features

- **Search & Play:** Search for any song on YouTube and instantly play it in the terminal (`echotune search "song name"`).
- **Auto-Play:** Automatically pick and play the most relevant search result (`--auto` or `-a` flag).
- **Background Audio:** Plays audio using `mpv` with `--no-video` for minimal resource usage.
- **Playback Controls:** Pause, play, seek forward/backward, and quit directly from the terminal.
- **Offline Downloads:** Download the currently playing song in high-quality Opus format (by pressing `d` in the controls).
- **Downloads Management:**
  - **List:** View all downloaded songs, with options to sort by date or title, and limit the output.
  - **Play:** Play downloaded songs by index or title. Supports flexible pipeline arguments like `shuffle`, `limit`, and `repeat` when playing your downloaded library (e.g., `echotune downloads play -l 5 -sh`).
  - **Remove:** Remove downloaded songs easily.
- **History & Caching:** Keeps track of your recently played songs and caches search results (with LRU eviction) for faster lookups.
- **Data Management:** Easily clear your cache, history, or all saved data via the `clear` command.

## Further Improvements

EchoTune is still in active development. Here are some planned improvements for the future:

- **Playlists Support:** Ability to create, manage, and play custom playlists of downloaded songs.
- **UI Overhaul:** Transitioning from a raw CLI interface to a rich, interactive terminal UI using the [Bubble Tea](https://github.com/charmbracelet/bubbletea) framework (once the core CLI features are finalized and polished).
- **Extended Playback Controls:** Fully implementing "Next", "Previous", and "Repeat" controls for uninterrupted listening sessions.
- **Cross-Platform Support:** Expanding official compatibility to macOS and Windows.
