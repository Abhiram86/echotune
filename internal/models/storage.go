package models

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const MaxCacheSize = 100
const MaxHistory = 100

type Storage struct {
	History   History
	Cache     Cache
	Downloads Downloads
	Playlists Playlists
}

type Cache struct {
	CachePath string
	Songs     map[string]CachedSong
}

type CachedSong struct {
	SelectedIndex int
	Results       []SearchResult
	Timestamp     time.Time
}

type History struct {
	HistoryPath string
	Songs       []SearchResult
}

type Downloads struct {
	mu            sync.Mutex
	DownloadsPath string
	MediaPath     string
	Songs         map[string]Download
}

type Download struct {
	Title     string
	Path      string
	Metadata  SearchResult
	Timestamp time.Time
}

type DownloadManager struct {
	Mu            sync.Mutex
	IsDownloading bool
}

type Playlist struct {
	Title     string
	Songs     map[string]Download
	Timestamp time.Time
}

type Playlists struct {
	PlaylistsPath string
	Playlists     map[string]Playlist
}

func loadJSON[T any](path string, target *T, fallback func()) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	if err := json.Unmarshal(data, target); err != nil {
		log.Println("error unmarshalling", path+":", err)
		fallback()
	}

	return nil
}

func (s *Storage) mountPaths() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	s.Cache.CachePath = filepath.Join(
		home,
		".cache",
		"echotune",
		"cache.json",
	)

	s.History.HistoryPath = filepath.Join(
		home,
		".local",
		"share",
		"echotune",
		"history.json",
	)

	s.Downloads.DownloadsPath = filepath.Join(
		home,
		".local",
		"share",
		"echotune",
		"downloads.json",
	)

	s.Downloads.MediaPath = filepath.Join(
		home,
		"Music",
		"echotune",
	)

	s.Playlists.PlaylistsPath = filepath.Join(
		home,
		".local",
		"share",
		"echotune",
		"playlists",
	)

	return nil
}

func (s *Storage) ensureDirectories() error {
	dirs := []string{
		filepath.Dir(s.Cache.CachePath),
		filepath.Dir(s.History.HistoryPath),
		filepath.Dir(s.Downloads.DownloadsPath),
		s.Downloads.MediaPath,
		s.Playlists.PlaylistsPath,
	}

	for _, dir := range dirs {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Storage) initDefaults() {
	s.Cache.Songs = make(map[string]CachedSong)
	s.History.Songs = make([]SearchResult, 0)
	s.Downloads.Songs = make(map[string]Download)
	s.Playlists.Playlists = make(map[string]Playlist)
}

func (s *Storage) Mount() error {
	err := s.mountPaths()
	if err != nil {
		return err
	}

	err = s.ensureDirectories()
	if err != nil {
		return err
	}

	s.initDefaults()

	return nil
}

func (s *Storage) LoadCache() error {
	return loadJSON(
		s.Cache.CachePath,
		&s.Cache,
		func() {
			s.Cache.Songs = make(map[string]CachedSong)
		},
	)
}

func (s *Storage) LoadHistory() error {
	return loadJSON(
		s.History.HistoryPath,
		&s.History,
		func() {
			s.History.Songs = make([]SearchResult, 0)
		},
	)
}

func (s *Storage) LoadDownloads() error {
	return loadJSON(
		s.Downloads.DownloadsPath,
		&s.Downloads,
		func() {
			s.Downloads.Songs = make(map[string]Download)
		},
	)
}

func (s *Storage) LoadPlaylists() error {
	files, err := os.ReadDir(s.Playlists.PlaylistsPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		if filepath.Ext(file.Name()) != ".json" {
			continue
		}

		fullPath := filepath.Join(
			s.Playlists.PlaylistsPath,
			file.Name(),
		)

		var playlist Playlist

		err = loadJSON(
			fullPath,
			&playlist,
			func() {
				playlist.Songs = make(map[string]Download, 0)
			},
		)
		if err != nil {
			return err
		}

		name := strings.TrimSuffix(
			file.Name(),
			".json",
		)

		s.Playlists.Playlists[name] = playlist
	}

	return nil
}

func (c *Cache) Clear() error {
	if _, err := os.Stat(c.CachePath); err == nil {
		err := os.Remove(c.CachePath)
		if err != nil {
			return err
		}
	}
	c.Songs = make(map[string]CachedSong)
	return nil
}

func (h *History) Clear() error {
	if _, err := os.Stat(h.HistoryPath); err == nil {
		err := os.Remove(h.HistoryPath)
		if err != nil {
			return err
		}
	}
	h.Songs = make([]SearchResult, 0)
	return nil
}

func (d *Downloads) Clear() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if _, err := os.Stat(d.DownloadsPath); err == nil {
		err := os.Remove(d.DownloadsPath)
		if err != nil {
			return err
		}
	}
	if _, err := os.Stat(d.MediaPath); err == nil {
		err := os.RemoveAll(d.MediaPath)
		if err != nil {
			return err
		}
	}
	d.Songs = make(map[string]Download)
	return nil
}

func (s *Storage) ClearAll() error {
	if err := s.Cache.Clear(); err != nil {
		return err
	}
	if err := s.History.Clear(); err != nil {
		return err
	}
	if err := s.Downloads.Clear(); err != nil {
		return err
	}

	return s.Playlists.ClearAll()
}

func saveToFile(path string, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func (c *Cache) Add(songs SearchList, selectedIndex int) error {
	if _, exists := c.Songs[songs.Query]; !exists && len(c.Songs) >= MaxCacheSize {
		c.evictOldest()
	}
	c.Songs[songs.Query] = CachedSong{
		SelectedIndex: selectedIndex,
		Results:       songs.Results,
		Timestamp:     time.Now(),
	}
	return saveToFile(c.CachePath, c)
}

func (c *Cache) Get(query string) (*CachedSong, bool) {
	song, ok := c.Songs[query]
	if ok {
		song.Timestamp = time.Now()
		c.Songs[query] = song
		_ = saveToFile(c.CachePath, c)
	}
	return &song, ok
}

func (c *Cache) evictOldest() {
	var oldestKey string
	var oldestTime time.Time

	// Iterate to find the entry with the minimum (oldest) timestamp
	for key, song := range c.Songs {
		if oldestKey == "" || song.Timestamp.Before(oldestTime) {
			oldestTime = song.Timestamp
			oldestKey = key
		}
	}

	if oldestKey != "" {
		delete(c.Songs, oldestKey)
	}
}

func (h *History) Add(song SearchResult) error {
	h.Songs = append(h.Songs, song)

	if len(h.Songs) > MaxHistory {
		h.Songs = h.Songs[len(h.Songs)-MaxHistory:]
	}

	return saveToFile(h.HistoryPath, h)
}

func (h *History) Get(idx int) (SearchResult, bool) {
	if idx < 0 || idx >= len(h.Songs) {
		return SearchResult{}, false
	}

	return h.Songs[idx], true
}

func (d *Downloads) Add(song Download) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.Songs[song.Metadata.ID] = song
	return saveToFile(d.DownloadsPath, d)
}

func (d *Downloads) Remove(song Download) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	id := song.Metadata.ID
	if s, exists := d.Songs[id]; exists {
		delete(d.Songs, id)
		err := os.RemoveAll(s.Path)
		if err != nil {
			return err
		}
		return saveToFile(d.DownloadsPath, d)
	}
	return nil
}

func (p *Playlists) ClearAll() error {
	for _, playlist := range p.Playlists {
		err := p.RemoveOne(playlist.Title)
		if err != nil {
			return err
		}
	}
	p.Playlists = make(map[string]Playlist)
	return nil
}

func (p *Playlists) SaveOne(title string) error {
	playlist, exists := p.Playlists[title]
	if !exists {
		return fmt.Errorf("playlist not found")
	}
	path := filepath.Join(p.PlaylistsPath, title+".json")
	return saveToFile(path, playlist)
}

func (p *Playlists) RemoveOne(title string) error {
	_, exists := p.Playlists[title]
	if !exists {
		return fmt.Errorf("playlist not found")
	}
	delete(p.Playlists, title)
	return os.Remove(filepath.Join(p.PlaylistsPath, title+".json"))
}

func (p *Playlists) AddPlayList(playlist Playlist) error {
	_, exists := p.Playlists[playlist.Title]
	if exists {
		return fmt.Errorf("playlist with title '%s' already exists", playlist.Title)
	}
	p.Playlists[playlist.Title] = playlist
	return nil
}

func (p *Playlists) RemovePlayList(playlist Playlist) error {
	_, exists := p.Playlists[playlist.Title]
	if !exists {
		return fmt.Errorf("playlist with title '%s' does not exist", playlist.Title)
	}
	return p.RemoveOne(playlist.Title)
}

func (p *Playlists) Get(title string) (*Playlist, bool) {
	playlist, exists := p.Playlists[title]
	return &playlist, exists
}

func (p *Playlists) AddSong(title string, song Download) error {
	playlist, exists := p.Playlists[title]
	if !exists {
		return fmt.Errorf("playlist not found")
	}
	_, exists = playlist.add(song)
	if exists {
		return fmt.Errorf("song with id '%s' already exists in playlist '%s'", song.Metadata.ID, title)
	}
	p.Playlists[title] = playlist
	p.SaveOne(title)
	return nil
}

func (p *Playlists) RemoveSong(title string, song Download) error {
	playlist, exists := p.Playlists[title]
	if !exists {
		return fmt.Errorf("playlist not found")
	}
	_, exists = playlist.remove(song)
	if !exists {
		return fmt.Errorf("song with id '%s' does not exist in playlist '%s'", song.Metadata.ID, title)
	}
	p.Playlists[title] = playlist
	p.SaveOne(title)
	return nil
}

func (p *Playlist) add(song Download) (*Playlist, bool) {
	_, exists := p.Songs[song.Metadata.ID]
	if exists {
		return p, true
	}
	p.Songs[song.Metadata.ID] = song
	return p, false
}

func (p *Playlist) remove(song Download) (*Playlist, bool) {
	_, exists := p.Songs[song.Metadata.ID]
	if !exists {
		return p, false
	}
	delete(p.Songs, song.Metadata.ID)
	return p, true
}
