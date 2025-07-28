package srt

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/OpenMediaStack/subtitle"
)

// Read reads SRT subtitles from an io.Reader and returns a slice of subtitle.Cue.
func Read(r io.Reader) ([]subtitle.Cue, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return parse(data)
}

// Open reads SRT subtitles from a file and returns a slice of subtitle.Cue.
func Open(file string) ([]subtitle.Cue, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	return parse(data)
}

// parse parses SRT subtitle data from a byte slice and returns a slice of subtitle.Cue.
func parse(data []byte) ([]subtitle.Cue, error) {
	r := strings.NewReader(string(data))
	s := bufio.NewScanner(r)

	var state int
	lineNb := 1
	var cues []subtitle.Cue
	var cue subtitle.Cue

	for s.Scan() {
		line := strings.TrimSpace(s.Text())

		if line == "" {
			if cue.Text != "" {
				cues = append(cues, cue)
				cue = subtitle.Cue{}
			}
			state = 0
			continue
		}

		switch state {
		case 0:
			if err := validateIndexLine(line); err != nil {
				return nil, fmt.Errorf("%w at line %d: '%s'", err, lineNb, line)
			}
			cue.Index, _ = strconv.Atoi(line)
			state = 1
		case 1:
			if err := validateTimecodeLine(line); err != nil {
				return nil, fmt.Errorf("%w at line %d: '%s'", err, lineNb, line)
			}
			parts := strings.Split(line, " --> ")
			startTimecode, err := subtitle.NewTimecode(parts[0])
			if err != nil {
				return nil, fmt.Errorf("%w at line %d: '%s'", err, lineNb, line)
			}
			cue.Start = startTimecode
			endTimecode, err := subtitle.NewTimecode(parts[1])
			if err != nil {
				return nil, fmt.Errorf("%w at line %d: '%s'", err, lineNb, line)
			}
			cue.End = endTimecode
			state = 2
		case 2:
			if cue.Text != "" {
				cue.Text += "\n"
			}
			cue.Text += line
		}

		lineNb++
	}

	if cue.Text != "" {
		cues = append(cues, cue)
	}

	if err := s.Err(); err != nil {
		return nil, err
	}

	return cues, nil
}

func validateIndexLine(line string) error {
	_, err := strconv.Atoi(line)
	if err != nil {
		return subtitle.ErrInvalidIndexLine
	}
	return nil
}

func validateTimecodeLine(line string) error {
	parts := strings.Split(line, " --> ")
	if len(parts) != 2 {
		return subtitle.ErrMissingTimecodeArrow
	}
	for _, tc := range parts {
		if !isValidTimecode(tc) {
			return subtitle.ErrInvalidTimecodeFormat
		}
	}
	return nil
}

func isValidTimecode(tc string) bool {
	if len(tc) != 12 {
		return false
	}

	parts := strings.Split(tc, ":")
	if len(parts) != 3 {
		return false
	}

	secondsParts := strings.Split(parts[2], ",")
	if len(secondsParts) != 2 {
		return false
	}

	return true
}
