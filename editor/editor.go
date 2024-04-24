package editor

import (
	"bufio"
	"fmt"
	"os"

	"golang.org/x/term"
)

type Editor struct {
	Rows    int
	Columns int
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
		}
	}
}

func (editor Editor) editorDrawRows() {
	for i := 0; i < editor.Rows-1; i++ {
		fmt.Print("~\r\n")
	}
	fmt.Print("~")
}

func (editor Editor) refreshScreen() {
	fmt.Print(clearEntireScreen)
	fmt.Print(placeCursorAtBegining)
	editor.editorDrawRows()
	fmt.Print(placeCursorAtBegining)
}

func readKey(reader bufio.Reader) rune {
	char, _, err := reader.ReadRune()
	if err != nil {
		ExitWithMessage("Couldn't read inserted character")
	}
	return char
}
