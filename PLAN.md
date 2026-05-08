# gosify — plan.md

## Goal

A tiny Spotify-ish CLI music player in Go.

Features:

* search songs
* play audio
* favourites
* history
* download audio

Uses:

* `yt-dlp` for searching/downloading
* `mpv` for playback

No accounts.
No recommendations.
No cloud bullshit.

---

# Tech Stack

## Core

* Go
* `urfave/cli`
* `yt-dlp`
* `mpv`

## Later

* Bubble Tea (interactive UI)
* SQLite (maybe)

---

# MVP First

## Commands

```bash
gosify play "song"
gosify search "song"
gosify history
gosify fav add "song"
gosify fav list
```

Keep it stupid simple initially.

---

# Project Structure

```txt
cmd/
internal/
    player/
    search/
    storage/
    models/
main.go
```

No overengineering.
No “services/repositories/managers/providers/factories”.
You're not deploying to NASA.

---

# Phase 1 — Basic Playback

## Target

Make this work:

```bash
gosify play "daft punk"
```

## Flow

1. Run:

   ```bash
   yt-dlp "ytsearch1:<query>"
   ```

2. Parse result

3. Extract audio/video URL

4. Run:

   ```bash
   mpv --no-video <url>
   ```

## Learn

* `os/exec`
* structs
* JSON parsing
* error handling

---

# Phase 2 — Search Command

## Target

```bash
gosify search "kendrick"
```

Show:

* title
* duration
* uploader

Later:

* choose result interactively

---

# Phase 3 — Storage

Use JSON first.

Example:

```json
{
  "history": [],
  "favourites": []
}
```

Store in:

```txt
~/.gosify/data.json
```

## Learn

* file handling
* encoding/json

---

# Phase 4 — Flags

Examples:

```bash
gosify search "song" -a
gosify history -s listens
```

## Learn

* CLI flags
* sorting
* slices

---

# Phase 5 — Concurrency

Only after MVP works.

Ideas:

* save history asynchronously
* preload metadata
* parallel downloads

## Learn

* goroutines
* channels
* sync

Do NOT add concurrency “because Go”.
That’s how you invent race conditions for sport.

---

# Phase 6 — Bubble Tea

Only for:

* search result selection
* loading states
* keyboard navigation

Not entire app rewrite.

---

# Nice To Have

* downloads
* playlists
* repeat/shuffle
* cache metadata
* mpv IPC control
* notifications

---

# Rules

* working > clever
* small commits
* CLI first, UI later
* avoid abstraction until repeated 3 times
* finish features before optimizing

---

# Definition of Done (MVP)

This command works reliably:

```bash
gosify play "song name"
```

If that works:
you already built something cooler than 90% of tutorial repos.

