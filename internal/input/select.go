package input

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
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
