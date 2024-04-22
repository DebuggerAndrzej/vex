package main

import (
	"bufio"
	"fmt"
	"os"
	"unicode"

	"golang.org/x/term"
)

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

func enterReaderLoop() {
	var char rune
	var err error
	reader := bufio.NewReader(os.Stdin)
	for string(char) != "q" {
		char, _, err = reader.ReadRune()
		if err != nil {
			exitWithMessage("Couldn't read inserted character")
		}
		if unicode.IsControl(char) {
			fmt.Printf("%v\r\n", char)
		} else {
			fmt.Printf("%c\r\n", char)
		}
	}
}
