package models

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const MaxCacheSize = 100
const MaxHistory = 100

type Storage struct {
	History   History
	Cache     Cache
	Downloads Downloads
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
	Title    string
	Path     string
	Metadata SearchResult
}

type DownloadManager struct {
	Mu            sync.Mutex
	IsDownloading bool
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

func (s *Storage) Mount() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	// paths
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

	// ensure directories exist
	err = os.MkdirAll(filepath.Dir(s.Cache.CachePath), 0755)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(s.History.HistoryPath), 0755)
	if err != nil {
		return err
	}

	err = os.MkdirAll(filepath.Dir(s.Downloads.DownloadsPath), 0755)
	if err != nil {
		return err
	}

	err = os.MkdirAll(s.Downloads.MediaPath, 0755)
	if err != nil {
		return err
	}

	// initialize safe defaults FIRST
	s.Cache.Songs = make(map[string]CachedSong)
	s.History.Songs = make([]SearchResult, 0)
	s.Downloads.Songs = make(map[string]Download)

	err = loadJSON(
		s.Cache.CachePath,
		&s.Cache,
		func() {
			s.Cache.Songs = make(map[string]CachedSong)
		},
	)
	if err != nil {
		return err
	}

	err = loadJSON(
		s.History.HistoryPath,
		&s.History,
		func() {
			s.History.Songs = make([]SearchResult, 0)
		},
	)
	if err != nil {
		return err
	}

	err = loadJSON(
		s.Downloads.DownloadsPath,
		&s.Downloads,
		func() {
			s.Downloads.Songs = make(map[string]Download)
		},
	)
	if err != nil {
		return err
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
	return nil
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
