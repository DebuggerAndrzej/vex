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
	TerminalRowCount    int
	TerminalColumnCount int
	Contents            strings.Builder
	Cursor              CursorPosition
	FileLines           []string
	NumberOfFileRows    int
	YOffset             int
	XOffset             int
}

func (editor *Editor) SetWindowSize() {
	columns, rows, err := term.GetSize(0)
	if err != nil {
		ExitWithMessage("Couldn't get terminal size")
	}

	editor.TerminalColumnCount, editor.TerminalRowCount = columns, rows
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
			for i := 0; i < editor.TerminalRowCount; i++ {
				editor.moveCursor(ARROW_UP)
			}
		case PAGE_DOWN:
			for i := 0; i < editor.TerminalRowCount; i++ {
				editor.moveCursor(ARROW_DOWN)
			}
		case HOME_KEY:
			editor.Cursor.X = 0
		case END_KEY:
			editor.Cursor.X = editor.TerminalColumnCount - 1
		case ARROW_UP, ARROW_DOWN, ARROW_LEFT, ARROW_RIGHT:
			editor.moveCursor(char)
		}
	}
}

func (editor *Editor) drawRows() {
	for y := 0; y < editor.TerminalRowCount; y++ {
		rowNumberAfterOffset := y + editor.YOffset
		if rowNumberAfterOffset >= editor.NumberOfFileRows {
			if editor.NumberOfFileRows == 0 && y == editor.TerminalRowCount/3 {
				editorTitleMsg := "Vex editor - pre alpha"
				pading := editor.TerminalColumnCount / 2
				editor.Contents.WriteString(fmt.Sprintf("%*s", pading, editorTitleMsg))
			} else {
				editor.Contents.WriteString("~")
			}

		} else {
			currentLine := editor.FileLines[rowNumberAfterOffset]
			if len(currentLine)-editor.XOffset > 0 {
				editor.Contents.WriteString(currentLine[editor.XOffset:min(editor.TerminalColumnCount+editor.XOffset, len(currentLine))])
			}
		}

		editor.Contents.WriteString(eraseRestOfTheLine)
		if y < editor.TerminalRowCount-1 {
			editor.Contents.WriteString("\r\n")
		}
	}
}

func (editor *Editor) refreshScreen() {
	editor.updateOffsets()
	editor.Contents.WriteString(hideCursor)
	editor.Contents.WriteString(placeCursorAtBegining)
	editor.drawRows()
	editor.Contents.WriteString(
		fmt.Sprintf("\x1b[%d;%dH", editor.Cursor.Y-editor.YOffset+1, editor.Cursor.X-editor.XOffset+1),
	)
	editor.Contents.WriteString(showCursor)
	fmt.Print(editor.Contents.String())
	editor.Contents.Reset()
}

func (editor *Editor) moveCursor(char rune) {
	var rowUnderCursorLen int

	if editor.Cursor.Y < editor.NumberOfFileRows {
		rowUnderCursorLen = len(editor.FileLines[editor.Cursor.Y])
	}

	switch char {
	case ARROW_UP:
		if editor.Cursor.Y != 0 {
			editor.Cursor.Y--
		}
	case ARROW_DOWN:
		if editor.Cursor.Y < editor.NumberOfFileRows {
			editor.Cursor.Y++
		}
	case ARROW_RIGHT:
		if editor.Cursor.X < rowUnderCursorLen {
			editor.Cursor.X++
		} else if editor.Cursor.X == rowUnderCursorLen && editor.Cursor.Y < editor.NumberOfFileRows {
			editor.Cursor.Y++
			editor.Cursor.X = 0
		}
	case ARROW_LEFT:
		if editor.Cursor.X != 0 {
			editor.Cursor.X--
		} else if editor.Cursor.Y > 0 {
			editor.Cursor.Y--
			editor.Cursor.X = len(editor.FileLines[editor.Cursor.Y])
		}
	}
	rowUnderCursorLen = 0

	if editor.Cursor.Y < editor.NumberOfFileRows {
		rowUnderCursorLen = len(editor.FileLines[editor.Cursor.Y])
	}
	if editor.Cursor.X > rowUnderCursorLen {
		editor.Cursor.X = rowUnderCursorLen
	}
}

func (editor *Editor) OpenFile(filePath string) {
	fileData, err := os.Open(filePath)
	if err != nil {
		ExitWithMessage("Couldn't load file")
	}

	fileScanner := bufio.NewScanner(fileData)
	for fileScanner.Scan() {
		editor.FileLines = append(editor.FileLines, fileScanner.Text())
		editor.NumberOfFileRows++
	}
}

func (editor *Editor) updateOffsets() {
	if editor.Cursor.Y < editor.YOffset {
		editor.YOffset = editor.Cursor.Y
	}

	if editor.Cursor.Y >= editor.YOffset+editor.TerminalRowCount {
		editor.YOffset = editor.Cursor.Y - editor.TerminalRowCount + 1
	}
	if editor.Cursor.X < editor.XOffset {
		editor.XOffset = editor.Cursor.X
	}
	if editor.Cursor.X >= editor.XOffset+editor.TerminalColumnCount {
		editor.XOffset = editor.Cursor.X - editor.TerminalColumnCount + 1
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
