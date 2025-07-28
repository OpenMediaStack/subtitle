package subtitle

type Cue struct {
	Index int
	Start Timecode
	End   Timecode
	Text  string
}
