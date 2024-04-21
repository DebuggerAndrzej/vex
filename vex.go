package main

import (
	"bufio"
	"fmt"
	"os"
	"unicode"

	"golang.org/x/term"
)

func main() {
	var err error
	var char rune
	previousState, err := term.MakeRaw(0)
	defer term.Restore(0, previousState)
	reader := bufio.NewReader(os.Stdin)
	for err == nil && string(char) != "q" {
		char, _, err = reader.ReadRune()
		if unicode.IsControl(char) {
			fmt.Printf("%v", char)
		}
		fmt.Printf("%c", char)
	}
}
