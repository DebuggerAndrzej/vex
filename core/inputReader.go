package core

import (
	"bufio"
)

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
