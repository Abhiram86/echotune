package internal

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/Abhiram86/echotune/internal/models"
)

func SearchQuery(query string) (*models.SearchList, error) {
	cmd := exec.Command(
		"yt-dlp",
		"--flat-playlist",
		"--dump-json",
		"ytsearch10:"+query,
	)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	var results []models.SearchResult

	scanner := bufio.NewScanner(stdout)

	for scanner.Scan() {
		item := models.SearchResult{}

		if err := json.Unmarshal(scanner.Bytes(), &item); err != nil {
			return nil, fmt.Errorf("failed to unmarshal yt-dlp output: %w", err)
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
		return nil, err
	}

	if err := cmd.Wait(); err != nil {
		return nil, err
	}

	return &models.SearchList{
		Results: results,
	}, nil
}
