package subtitle

import "errors"

var (
	ErrInvalidIndexLine      = errors.New("invalid index line: must be numeric")
	ErrInvalidTimecodeLine   = errors.New("invalid timecode line: must match 'HH:MM:SS,mmm --> HH:MM:SS,mmm'")
	ErrInvalidTimecodeFormat = errors.New("invalid timecode format")
	ErrMissingTimecodeArrow  = errors.New("invalid timecode separator")
)
