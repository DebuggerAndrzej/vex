package core

import (
	"fmt"
	"os"
)

func ExitWithMessage(message string) {
	fmt.Println(message)
	os.Exit(1)
}
