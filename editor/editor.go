package editor

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"golang.org/x/term"
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

func (editor Editor) EnterReaderLoop() {
	reader := bufio.NewReader(os.Stdin)
	for {
		editor.refreshScreen()
		switch char := readKey(*reader); char {
		case rune(ctrlKey & byte('q')):
			fmt.Print(clearEntireScreen)
			fmt.Print(placeCursorAtBegining)
			os.Exit(0)
		case 'w', 's', 'a', 'd':
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
	case 'w':
		editor.Cursor.Y--
	case 's':
		editor.Cursor.Y++
	case 'a':
		editor.Cursor.X--
	case 'd':
		editor.Cursor.X++
	}
}

func readKey(reader bufio.Reader) rune {
	char, _, err := reader.ReadRune()
	if err != nil {
		ExitWithMessage("Couldn't read inserted character")
	}
	return char
}
