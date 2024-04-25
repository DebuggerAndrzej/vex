package editor

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
)

const (
	ARROW_UP = 1000 + iota
	ARROW_DOWN
	ARROW_LEFT
	ARROW_RIGHT
	DEL_KEY
	HOME_KEY
	END_KEY
	PAGE_UP
	PAGE_DOWN
)

type CursorPosition struct {
	X, Y int
}

type Editor struct {
	Rows     int
	Columns  int
	Contents strings.Builder
	Cursor   CursorPosition
}

func (editor *Editor) SetWindowSize() {
	columns, rows, err := term.GetSize(0)
	if err != nil {
		ExitWithMessage("Couldn't get terminal size")
	}

	editor.Columns, editor.Rows = columns, rows
}

func (editor *Editor) EnterReaderLoop() {
	reader := bufio.NewReader(os.Stdin)
	for {
		editor.refreshScreen()
		switch char := readKey(*reader); char {
		case rune(ctrlKey & byte('q')):
			fmt.Print(clearEntireScreen)
			fmt.Print(placeCursorAtBegining)
			os.Exit(0)
		case PAGE_UP:
			for i := 0; i < editor.Rows; i++ {
				editor.moveCursor(ARROW_UP)
			}
		case PAGE_DOWN:
			for i := 0; i < editor.Rows; i++ {
				editor.moveCursor(ARROW_DOWN)
			}
		case HOME_KEY:
			editor.Cursor.X = 0
		case END_KEY:
			editor.Cursor.X = editor.Columns - 1
		case ARROW_UP, ARROW_DOWN, ARROW_LEFT, ARROW_RIGHT:
			editor.moveCursor(char)
		}
	}
}

func (editor *Editor) editorDrawRows() {
	for y := 0; y < editor.Rows; y++ {
		if y == editor.Rows/3 {
			editorTitleMsg := "Vex editor - pre alpha"
			pading := editor.Columns / 2
			editor.Contents.WriteString(fmt.Sprintf("%*s", pading, editorTitleMsg))
		} else {
			editor.Contents.WriteString("~")
		}
		editor.Contents.WriteString(eraseRestOfTheLine)
		if y < editor.Rows-1 {
			editor.Contents.WriteString("\r\n")
		}
	}
}

func (editor *Editor) refreshScreen() {
	editor.Contents.WriteString(hideCursor)
	editor.Contents.WriteString(placeCursorAtBegining)
	editor.editorDrawRows()
	editor.Contents.WriteString(fmt.Sprintf("\x1b[%d;%dH", editor.Cursor.Y, editor.Cursor.X))
	editor.Contents.WriteString(showCursor)
	fmt.Print(editor.Contents.String())
	editor.Contents.Reset()
}

func (editor *Editor) moveCursor(char rune) {
	switch char {
	case ARROW_UP:
		if editor.Cursor.Y != 0 {
			editor.Cursor.Y--
		}
	case ARROW_DOWN:
		if editor.Cursor.Y != editor.Rows-1 {
			editor.Cursor.Y++
		}
	case ARROW_RIGHT:
		if editor.Cursor.X != editor.Columns-1 {
			editor.Cursor.X++
		}
	case ARROW_LEFT:
		if editor.Cursor.X != 0 {
			editor.Cursor.X--
		}
	}
}

func readKey(reader bufio.Reader) rune {
	char, _, err := reader.ReadRune()
	if err != nil {
		ExitWithMessage("Couldn't read inserted character")
	}

	if doesStartWithEscapeCharacter(byte(char)) {
		nextChar, _, _ := reader.ReadRune()
		if nextChar == '[' {
			nextNextChar, _, _ := reader.ReadRune()
			if nextNextChar >= '0' && nextNextChar <= '9' {
				nextNextNextChar, _, _ := reader.ReadRune()
				if nextNextNextChar == '~' {
					switch nextNextChar {
					case '1', '7':
						return HOME_KEY
					case '3':
						return DEL_KEY
					case '4', '8':
						return END_KEY
					case '5':
						return PAGE_UP
					case '6':
						return PAGE_DOWN
					}
				}
			} else {

				switch nextNextChar {
				case 'A':
					return ARROW_UP
				case 'B':
					return ARROW_DOWN
				case 'C':
					return ARROW_RIGHT
				case 'D':
					return ARROW_LEFT
				case 'H':
					return HOME_KEY
				case 'F':
					return END_KEY
				}
			}

		}
	}
	if char == '0' {
		nextChar, _, _ := reader.ReadRune()
		switch nextChar {
		case 'H':
			return HOME_KEY
		case 'F':
			return END_KEY
		}
	}
	return char
}

func doesStartWithEscapeCharacter(sequence byte) bool {
	if sequence == 27 {
		return true
	}
	return false
}
