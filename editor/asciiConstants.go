package editor

const ctrlKey = byte(0b00011111)

const (
	clearEntireScreen     = "\x1b[2J"
	placeCursorAtBegining = "\x1b[H"
)
