package models

import (
	"encoding/json"
	"log"
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
			s.Downloads.Songs = make([]Download, 0)
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
	if _, err := os.Stat(d.DownloadsPath); err == nil {
		err := os.Remove(d.DownloadsPath)
		if err != nil {
			return err
		}
	}
	d.Songs = make([]Download, 0)
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
