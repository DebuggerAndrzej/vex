package main

import (
	"github.com/DebuggerAndrzej/vex/editor"
	"golang.org/x/term"
)

func main() {
	previousState, err := term.MakeRaw(0)
	if err != nil {
		editor.ExitWithMessage("Failed to init raw terminal mode")
	}
	defer term.Restore(0, previousState)

	editor := editor.Editor{}
	editor.SetWindowSize()
	editor.EnterReaderLoop()
}
