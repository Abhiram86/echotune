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

### Search & Play
Search for any song on YouTube and instantly play it in the terminal.
```bash
./et search "song name"
```
- **Auto-Play:** Automatically pick and play the most relevant result (`--auto` or `-a`).
- **Limit & Repeat:** Control the number of results (`--limit`) or repeat playback (`--repeat`).

### Downloads Management
- **Download:** Press `d` while a song is playing to download it in high-quality Opus format.
- **List:** View all downloaded songs with sorting by date or title, and limit output.
- **Play:** Play a downloaded song by index or title. Supports pipeline arguments for full playlist control.
  ```bash
  ./et downloads play -l 5 -sh
  ```
- **Remove:** Remove downloaded songs by index or title.

### Playlists
Create, manage, and play custom playlists of downloaded songs.
- **Add to Playlist:** While a song is playing, press `a` to add it to a playlist.
- **List:** View all saved playlists with sorting and limit options.
- **Play:** Play songs in a playlist with shuffle, limit, and repeat support.
- **Remove:** Remove a whole playlist or a specific song from a playlist.
- **Clear:** Delete all playlists with confirmation.

### History & Caching
- **History:** View recently played songs with reverse chronological order and limit.
- **Cache:** Search results are cached with LRU eviction for faster repeat lookups.

### Playback Controls
- **Pause/Resume:** Toggle playback with `p`.
- **Seek:** Skip forward/backward by 5 seconds with `f` and `b`.
- **Skip Tracks:** Next (`n`) and Previous (`v`) support for queue and playlist playback.
- **Quit:** Stop playback and return to the shell with `q`.

### Data Management
Easily clear your cache, history, downloads, or all saved data via the `clear` command.
```bash
./et clear all       # Clear everything with confirmation
./et clear cache     # Clear search cache
./et clear history   # Clear playback history
```

## Architecture

EchoTune uses a `PlaybackSession` model to unify the player state and song queue across all commands. This allows features like "Next", "Previous", and repeat to work consistently whether you are playing a single searched song, a downloaded track, or an entire playlist.

Playlists are saved as individual JSON files under `~/.local/share/echotune/playlists/`, and the queue resolves downloaded songs to local files automatically (O(1) lookup) before falling back to streaming.

## Further Improvements

EchoTune is still in active development. Here are some planned improvements for the future:

- **UI Overhaul:** Transitioning from a raw CLI interface to a rich, interactive terminal UI using the [Bubble Tea](https://github.com/charmbracelet/bubbletea) framework (once the core CLI features are finalized and polished).
- **Cross-Platform Support:** Expanding official compatibility to macOS and Windows.
- **Playlist Queue Management:** Reordering songs within a playlist queue during playback.
- **Equalizer & Audio Settings:** Basic audio filtering options via mpv integration.
