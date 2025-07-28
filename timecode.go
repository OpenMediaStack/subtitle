package subtitle

import (
	"fmt"
	"strconv"
	"strings"
)

type Timecode struct {
	Hour        int
	Minute      int
	Second      int
	Millisecond int
}

func (t Timecode) String() string {
	return fmt.Sprintf("%02d:%02d:%02d,%03d", t.Hour, t.Minute, t.Second, t.Millisecond)
}

// NewTimecode creates a new Timecode from a string in the format "HH:MM:SS,mmm".
func NewTimecode(input string) (Timecode, error) {
	splitted := strings.Split(strings.TrimSpace(input), ",")
	ms, err := strconv.Atoi(splitted[1])
	if err != nil {
		return Timecode{}, err
	}
	times := strings.Split(splitted[0], ":")
	hour, err := strconv.Atoi(times[0])
	if err != nil {
		return Timecode{}, err
	}
	minute, err := strconv.Atoi(times[1])
	if err != nil {
		return Timecode{}, err
	}
	second, err := strconv.Atoi(times[2])
	if err != nil {
		return Timecode{}, err
	}

	return Timecode{
		Hour:        hour,
		Minute:      minute,
		Second:      second,
		Millisecond: ms,
	}, nil
}
