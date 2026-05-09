package models

import (
	"encoding/json"
	"os"
	"path/filepath"
)

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
}

type History struct {
	HistoryPath string
	Songs       []SearchResult
}

type Downloads struct {
	DownloadsPath string
	Songs         []Download
}

type Download struct {
	Title    string
	Path     string
	Metadata SearchResult
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

	// initialize safe defaults FIRST
	s.Cache.Songs = make(map[string]CachedSong)
	s.History.Songs = make([]SearchResult, 0)
	s.Downloads.Songs = make([]Download, 0)

	// load cache
	if data, err := os.ReadFile(s.Cache.CachePath); err == nil {
		json.Unmarshal(data, &s.Cache)
	}

	// load history
	if data, err := os.ReadFile(s.History.HistoryPath); err == nil {
		json.Unmarshal(data, &s.History)
	}

	// load downloads
	if data, err := os.ReadFile(s.Downloads.DownloadsPath); err == nil {
		json.Unmarshal(data, &s.Downloads)
	}

	return nil
}

func (s *Storage) Clear() error {
	err := os.Remove(s.Cache.CachePath)
	if err != nil {
		return err
	}

	err = os.Remove(s.History.HistoryPath)
	if err != nil {
		return err
	}

	err = os.Remove(s.Downloads.DownloadsPath)
	if err != nil {
		return err
	}

	s.Cache.Songs = make(map[string]CachedSong)
	s.History.Songs = make([]SearchResult, 0)
	s.Downloads.Songs = make([]Download, 0)

	return nil
}

func saveToFile(path string, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// func (s *Storage) Save() error {
// 	err := saveToFile(s.Cache.CachePath, s.Cache)
// 	if err != nil {
// 		return err
// 	}

// 	err = saveToFile(s.History.HistoryPath, s.History)
// 	if err != nil {
// 		return err
// 	}

// 	return saveToFile(s.Downloads.DownloadsPath, s.Downloads)
// }

func (c *Cache) Add(songs SearchList, selectedIndex int) error {
	c.Songs[songs.Query] = CachedSong{
		SelectedIndex: selectedIndex,
		Results:       songs.Results,
	}
	return saveToFile(c.CachePath, c)
}

func (c *Cache) Get(query string) (*CachedSong, bool) {
	song, ok := c.Songs[query]
	return &song, ok
}

func (c *Cache) Remove(query string) {
	delete(c.Songs, query)
}

func (h *History) Add(song SearchResult) error {
	h.Songs = append(h.Songs, song)

	if len(h.Songs) > MaxHistory {
		h.Songs = h.Songs[1:]
	}

	return saveToFile(h.HistoryPath, h)
}

func (h *History) Get(idx int) (SearchResult, bool) {
	if idx < 0 || idx >= len(h.Songs) {
		return SearchResult{}, false
	}

	return h.Songs[idx], true
}
