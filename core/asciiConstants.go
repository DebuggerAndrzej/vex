package core

const (
	ctrlKey = byte(0b00011111)
	enter   = '\r'
)

const (
	ansiEscape            = '\x1b'
	clearEntireScreen     = "\x1b[2J"
	placeCursorAtBegining = "\x1b[H"
	hideCursor            = "\x1b[?25l"
	showCursor            = "\x1b[?25h"
	eraseRestOfTheLine    = "\x1b[K"
)
