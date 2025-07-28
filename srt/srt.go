package srt

import (
	"bufio"
	"bytes"
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
	b := bytes.NewReader(data)
	s := bufio.NewScanner(b)

	var state int
	var cues []subtitle.Cue
	var cue subtitle.Cue

	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		switch state {
		case 0:
			if line == "" {
				continue
			}
			index, err := strconv.Atoi(line)
			if err != nil {
				return nil, err
			}
			cue = subtitle.Cue{}
			cue.Index = index
			state = 1
		case 1:
			times := strings.Split(line, "-->")
			start, err := subtitle.NewTimecode(times[0])
			if err != nil {
				return nil, err
			}
			cue.Start = start

			end, err := subtitle.NewTimecode(times[1])
			if err != nil {
				return nil, err
			}
			cue.End = end
			state = 2
			continue
		case 2:
			if line == "" {
				cues = append(cues, cue)
				state = 0
				continue
			}
			if cue.Text != "" {
				cue.Text += "\n"
			}
			cue.Text += line

		}

	}
	if state == 2 {
		cues = append(cues, cue)
	}
	return cues, nil
}
