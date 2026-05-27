# EchoTune

EchoTune is a fast, terminal-based music player and downloader. It allows you to search for songs on YouTube, play them directly in your terminal without any video overhead, and even download them for offline playback. Built with Go, it provides a clean, distraction-free music experience right from your command line.

<video src="https://github.com/Abhiram86/echotune/raw/main/media/demo.mp4" width="100%" controls></video>

## Prerequisites & Installation

For full details on installing EchoTune and its required dependencies (`mpv` and `yt-dlp`) for Linux, macOS, or Windows, please see our detailed **[Installation Guide](installation/install.md)**.

## How to Use / Commands Documentation

EchoTune uses a simple CLI interface. Below is a detailed explanation of each command.

### `search`
Search for any song on YouTube and instantly play it in the terminal.
```bash
./etune search "song name"
```
- **Interactive TUI:** By default, displays a beautiful terminal UI where you can navigate search results with arrow keys/`n`/`b` and press `Enter` to play.
- **Auto-Play:** Automatically pick and play the most relevant result by passing `--auto` or `-a`.
- **Limit & Repeat:** Control the maximum number of results to fetch (`--limit` or `-l`) or repeat the playback (`--repeat` or `-r`).

### `history`
Show your recently played songs.
```bash
./etune history
```
- **Limit:** Use `--limit` or `-l` to restrict the number of history items shown.

### `downloads`
Manage and play your downloaded songs.
- **`list`:** View all downloaded songs.
  ```bash
  ./etune downloads list
  ```
  - Options: `--sort` or `-s` (sort by download date instead of title), `--limit` or `-l`.
- **`play`:** Play a specific downloaded song by its name or index, or play all.
  ```bash
  ./etune downloads play "song name"
  ./etune downloads play 1
  ./etune downloads play       # Plays all downloads
  ```
  - Options: `--shuffle` or `-sh` (play in random order), `--limit` or `-l`, `--repeat` or `-r`.
- **`remove`:** Remove a downloaded song by index or name.
  ```bash
  ./etune downloads remove "song name"
  ```

### `playlist`
Create, manage, and play custom playlists of downloaded songs. You can add songs to playlists while they are playing using the `a` hotkey in the player.
- **`list`:** View all saved playlists.
  ```bash
  ./etune playlist list
  ```
- **`play`:** Play all songs in a specific playlist.
  ```bash
  ./etune playlist play "playlist name"
  ```
  - Options: `--shuffle` or `-sh`, `--limit` or `-l`, `--repeat` or `-r`.
- **`remove`:** Remove an entire playlist, or a specific song from a playlist.
  ```bash
  ./etune playlist remove "playlist name"
  ./etune playlist remove "playlist name" "song name"
  ```
- **`clear`:** Delete all playlists with a confirmation prompt.
  ```bash
  ./etune playlist clear
  ```

### `clear`
Easily manage your local application data.
```bash
./etune clear cache     # Clear search cache (speeds up repeat searches)
./etune clear history   # Clear your playback history
./etune clear all       # Wipes everything: cache, history, downloads, and playlists
```

## Playback Controls

While a song is actively playing, EchoTune uses a dynamic Bubble Tea TUI. You can control playback instantly with the following hotkeys:

- `Space` or `p`: Toggle Play/Pause
- `k` or `Right Arrow`: Seek forward 5 seconds
- `j` or `Left Arrow`: Seek backward 5 seconds
- `n` or `Up Arrow`: Skip to the next track (in a playlist/queue)
- `b` or `Down Arrow`: Skip to the previous track
- `d`: Download the currently playing song in the background (saved as high-quality Opus)
- `a`: Add the current song to a playlist
- `x`: Remove the current song from a playlist
- `q` or `Esc`: Quit the player

## Architecture

EchoTune uses a `PlaybackSession` model to unify the player state and song queue across all commands. This allows features like "Next", "Previous", and repeat to work consistently whether you are playing a single searched song, a downloaded track, or an entire playlist.

Playlists and application data are saved in standard OS-specific directories (e.g., `~/.local/share/echotune` on Linux, `AppData\Roaming\echotune` on Windows, and `Library/Application Support/echotune` on macOS). The queue resolves downloaded songs to local files automatically (O(1) lookup) before falling back to streaming.

## Potential Improvements

EchoTune works great for its current scope. While not intended to be a massive "full" application, here are some *potential* improvements for the future:

- **Playlist Queue Management:** Reordering songs within a playlist queue during playback.
- **Equalizer & Audio Settings:** Basic audio filtering options via mpv integration.
