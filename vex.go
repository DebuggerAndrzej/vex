package main

import (
	"bufio"
	"fmt"
	"os"

	"golang.org/x/term"
)

const ctrlKey = byte(0b00011111)

func main() {
	previousState, err := term.MakeRaw(0)
	if err != nil {
		exitWithMessage("Failed to init raw terminal mode")
	}
	defer term.Restore(0, previousState)

	enterReaderLoop()
}

func exitWithMessage(message string) {
	fmt.Println(message)
	os.Exit(1)
}

func editorReadKey(reader bufio.Reader) rune {
	char, _, err := reader.ReadRune()
	if err != nil {
		exitWithMessage("Couldn't read inserted character")
	}
	return char
}

func editorDrawRows() {
	for i := 0; i < 24; i++ {
		fmt.Print("~\r\n")
	}

}

func editorRefreshScreen() {
	fmt.Print("\x1b[2J")
	fmt.Print("\x1b[H")
	editorDrawRows()
	fmt.Print("\x1b[H")
}

func enterReaderLoop() {
	reader := bufio.NewReader(os.Stdin)
	for {
		editorRefreshScreen()
		switch char := editorReadKey(*reader); char {
		case rune(ctrlKey & byte('q')):
			fmt.Print("\x1b[2J")
			fmt.Print("\x1b[H")
			os.Exit(0)
		}
	}
}
