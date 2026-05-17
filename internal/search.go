package internal

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/Abhiram86/echotune/internal/models"
)

func SearchQuery(ctx context.Context, query string, storage *models.Storage, length string) (*models.CachedSong, *models.SearchList, error) {
	// check cache
	cached, ok := storage.Cache.Get(query)
	if ok {
		return cached, nil, nil
	}

	searchTarget := fmt.Sprintf("ytsearch%s:%s", length, query)

	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx,
		"yt-dlp",
		"--js-runtimes", "node",
		"--flat-playlist",
		"--dump-json",
		searchTarget,
	)

	defer func() {
		if cmd.Process != nil {
			_ = cmd.Process.Kill()
		}
	}()

	cmd.Stderr = os.Stderr

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, nil, err
	}

	var results []models.SearchResult

	scanner := bufio.NewScanner(stdout)

	for scanner.Scan() {
		item := models.SearchResult{}

		if err := json.Unmarshal(scanner.Bytes(), &item); err != nil {
			return nil, nil, fmt.Errorf("failed to unmarshal yt-dlp output: %w", err)
		}

		results = append(results, models.SearchResult{
			ID:         item.ID,
			Title:      item.Title,
			URL:        "https://youtube.com/watch?v=" + item.ID,
			Duration:   item.Duration,
			ViewCount:  item.ViewCount,
			LikeCount:  item.LikeCount,
			Uploader:   item.Uploader,
			Channel:    item.Channel,
			UploadDate: item.UploadDate,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}

	if err := cmd.Wait(); err != nil {
		return nil, nil, err
	}

	return nil, &models.SearchList{
		Query:   query,
		Results: results,
	}, nil
}
