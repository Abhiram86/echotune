package input

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"

	"github.com/Abhiram86/echotune/internal/models"
)

func ReadSelection(
	reader *bufio.Reader,
	max int,
) (int, error) {
	fmt.Print("\nEnter the number of the song you want to play: ")

	input, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}

	idx, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil {
		return 0, fmt.Errorf("invalid selection")
	}

	if idx < 1 || idx > max {
		return 0, fmt.Errorf("selection out of range")
	}

	return idx - 1, nil
}

const views_weight float32 = 0.5
const likes_weight float32 = 0.5

func balancedScore(viewCount, likeCount int) float32 {
	if viewCount+likeCount == 0 {
		return 0
	}
	return (float32(viewCount)*views_weight + float32(likeCount)*likes_weight) / float32(viewCount+likeCount)
}

func SelectBestSong(songs *models.SearchList) (int, error) {
	if len(songs.Results) == 0 {
		return 0, fmt.Errorf("no results found")
	}

	bestScore := float32(0)
	bestIdx := 0

	for idx, song := range songs.Results {
		score := balancedScore(song.ViewCount, song.LikeCount)
		if score > bestScore {
			bestScore = score
			bestIdx = idx
		}
	}

	return bestIdx, nil
}
