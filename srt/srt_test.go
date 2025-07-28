package srt

import (
	"errors"
	"testing"

	"github.com/OpenMediaStack/subtitle"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expected      []subtitle.Cue
		expectedError error
	}{
		{
			name: "Valid SRT",
			input: `1
00:00:01,000 --> 00:00:02,000
Hello, World!

2
00:00:03,000 --> 00:00:04,000
Goodbye, World!`,
			expected: []subtitle.Cue{
				{
					Index: 1,
					Start: subtitle.Timecode{Hour: 0, Minute: 0, Second: 1, Millisecond: 0},
					End:   subtitle.Timecode{Hour: 0, Minute: 0, Second: 2, Millisecond: 0},
					Text:  "Hello, World!",
				},
				{
					Index: 2,
					Start: subtitle.Timecode{Hour: 0, Minute: 0, Second: 3, Millisecond: 0},
					End:   subtitle.Timecode{Hour: 0, Minute: 0, Second: 4, Millisecond: 0},
					Text:  "Goodbye, World!",
				},
			},
			expectedError: nil,
		},
		{
			name: "Invalid index",
			input: `X
00:00:01,000 --> 00:00:02,000
Bad index`,
			expectedError: subtitle.ErrInvalidIndexLine,
		},
		{
			name: "Missing timecode separator",
			input: `1
00:00:01,000 - 00:00:02,000
Bad timecode line`,
			expectedError: subtitle.ErrMissingTimecodeArrow,
		},
		{
			name: "Malformed timecode",
			input: `1
00:00:01 --> 00:00:02,000
Missing milliseconds`,
			expectedError: subtitle.ErrInvalidTimecodeFormat,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := []byte(tt.input)
			cues, err := parse(data)

			if tt.expectedError != nil {
				if !errors.Is(err, tt.expectedError) {
					t.Fatalf("expected error '%v', got '%v'", tt.expectedError, err)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(cues) != len(tt.expected) {
				t.Errorf("got %d cues, want %d", len(cues), len(tt.expected))
			}

			for i, cue := range cues {
				if cue != tt.expected[i] {
					t.Errorf("cue %d = %+v, want %+v", i+1, cue, tt.expected[i])
				}
			}
		})
	}
}
