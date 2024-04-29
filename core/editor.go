package core

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"golang.org/x/term"
)

const TAB_LENGHT = 4

type CursorPosition struct {
	X, Y int
}

type Offset struct {
	X, Y int
}

type TerminalSize struct {
	Rows, Columns int
}

type StatusMessage struct {
	Message   string
	Timestamp time.Time
}

type Editor struct {
	contents      strings.Builder
	terminalSize  TerminalSize
	cursor        CursorPosition
	offset        Offset
	fileLines     []string
	fileName      string
	statusMessage StatusMessage
}

func (editor *Editor) SetWindowSize() {
	columns, rows, err := term.GetSize(0)
	if err != nil {
		ExitWithMessage("Couldn't get terminal size")
	}

	editor.terminalSize = TerminalSize{Rows: rows - 2, Columns: columns}
}
func (editor *Editor) OpenFile(filePath string) {
	fileData, err := os.Open(filePath)
	if err != nil {
		ExitWithMessage("Couldn't load file")
	}
	editor.fileName = filePath
	fileScanner := bufio.NewScanner(fileData)
	for fileScanner.Scan() {
		editor.fileLines = append(
			editor.fileLines,
			strings.ReplaceAll(fileScanner.Text(), "\t", strings.Repeat(" ", TAB_LENGHT)),
		)
	}
}

func (editor *Editor) SetStatusMessage(message string) {
	editor.statusMessage = StatusMessage{Message: message, Timestamp: time.Now()}
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
			editor.cursor.Y = editor.offset.Y
			for i := 0; i < editor.terminalSize.Rows; i++ {
				editor.moveCursor(ARROW_UP)
			}
		case PAGE_DOWN:
			editor.cursor.Y = editor.offset.Y + editor.terminalSize.Rows - 1
			for i := 0; i < editor.terminalSize.Rows; i++ {
				editor.moveCursor(ARROW_DOWN)
			}
		case HOME_KEY:
			editor.cursor.X = 0
		case END_KEY:
			if editor.cursor.Y < len(editor.fileLines) {
				editor.cursor.X = len(editor.fileLines[editor.cursor.Y])
			}
		case ARROW_UP, ARROW_DOWN, ARROW_LEFT, ARROW_RIGHT:
			editor.moveCursor(char)
		}
	}
}

func (editor *Editor) refreshScreen() {
	editor.updateOffsets()
	editor.contents.WriteString(hideCursor)
	editor.contents.WriteString(placeCursorAtBegining)
	editor.drawRows()
	editor.drawStatusBar()
	editor.drawMessageBar()
	editor.contents.WriteString(
		fmt.Sprintf("\x1b[%d;%dH", editor.cursor.Y-editor.offset.Y+1, editor.cursor.X-editor.offset.X+1),
	)
	editor.contents.WriteString(showCursor)
	fmt.Print(editor.contents.String())
	editor.contents.Reset()
}

func (editor *Editor) drawRows() {
	for y := 0; y < editor.terminalSize.Rows; y++ {
		rowNumberAfterOffset := y + editor.offset.Y
		if rowNumberAfterOffset >= len(editor.fileLines) {
			if len(editor.fileLines) == 0 && y == editor.terminalSize.Rows/3 {
				editorTitleMsg := "Vex editor - pre alpha"
				padding := editor.terminalSize.Columns / 2
				editor.contents.WriteString(fmt.Sprintf("%*s", padding, editorTitleMsg))
			} else {
				editor.contents.WriteString("~")
			}

		} else {
			currentLine := editor.fileLines[rowNumberAfterOffset]
			if len(currentLine)-editor.offset.X > 0 {
				editor.contents.WriteString(currentLine[editor.offset.X:min(editor.terminalSize.Columns+editor.offset.X, len(currentLine))])
			}
		}

		editor.contents.WriteString(eraseRestOfTheLine)
		editor.contents.WriteString("\r\n")
	}
}

func (editor *Editor) moveCursor(char rune) {
	var rowUnderCursorLen int

	if editor.cursor.Y < len(editor.fileLines) {
		rowUnderCursorLen = len(editor.fileLines[editor.cursor.Y])
	}

	switch char {
	case ARROW_UP:
		if editor.cursor.Y != 0 {
			editor.cursor.Y--
		}
	case ARROW_DOWN:
		if editor.cursor.Y < len(editor.fileLines) {
			editor.cursor.Y++
		}
	case ARROW_RIGHT:
		if editor.cursor.X < rowUnderCursorLen {
			editor.cursor.X++
		} else if editor.cursor.X == rowUnderCursorLen && editor.cursor.Y < len(editor.fileLines) {
			editor.cursor.Y++
			editor.cursor.X = 0
		}
	case ARROW_LEFT:
		if editor.cursor.X != 0 {
			editor.cursor.X--
		} else if editor.cursor.Y > 0 {
			editor.cursor.Y--
			editor.cursor.X = len(editor.fileLines[editor.cursor.Y])
		}
	}
	rowUnderCursorLen = 0

	if editor.cursor.Y < len(editor.fileLines) {
		rowUnderCursorLen = len(editor.fileLines[editor.cursor.Y])
	}
	if editor.cursor.X > rowUnderCursorLen {
		editor.cursor.X = rowUnderCursorLen
	}
}

func (editor *Editor) updateOffsets() {
	if editor.cursor.Y < editor.offset.Y {
		editor.offset.Y = editor.cursor.Y
	}
	if editor.cursor.Y >= editor.offset.Y+editor.terminalSize.Rows {
		editor.offset.Y = editor.cursor.Y - editor.terminalSize.Rows + 1
	}
	if editor.cursor.X < editor.offset.X {
		editor.offset.X = editor.cursor.X
	}
	if editor.cursor.X >= editor.offset.X+editor.terminalSize.Columns {
		editor.offset.X = editor.cursor.X - editor.terminalSize.Columns + 1
	}
}

func (editor *Editor) drawStatusBar() {
	fileInfo := fmt.Sprintf("%s - %d", editor.fileName[0:min(len(editor.fileName), 20)], len(editor.fileLines))
	cursorInfo := fmt.Sprintf("%d:%d", editor.cursor.Y+1, editor.cursor.X)
	editor.contents.WriteString(
		fmt.Sprintf(
			"\x1b[7m%s%*s%s\x1b[m\r\n",
			fileInfo,
			editor.terminalSize.Columns-len(fileInfo)-len(cursorInfo),
			"",
			cursorInfo,
		),
	)
}

func (editor *Editor) drawMessageBar() {
	editor.contents.WriteString(eraseRestOfTheLine)
	if editor.statusMessage.Message != "" && time.Now().Sub(editor.statusMessage.Timestamp).Seconds() < 5 {
		editor.contents.WriteString(
			editor.statusMessage.Message[0:min(len(editor.statusMessage.Message), editor.terminalSize.Columns)],
		)
	}
}
