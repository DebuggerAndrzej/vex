package main

import (
	"flag"

	"github.com/DebuggerAndrzej/vex/core"

	"golang.org/x/term"
)

func main() {
	flag.Parse()
	previousState, err := term.MakeRaw(0)
	if err != nil {
		core.ExitWithMessage("Failed to init raw terminal mode")
	}
	defer term.Restore(0, previousState)

	editor := core.Editor{}
	editor.SetWindowSize()
	editor.OpenFile(flag.Arg(0))
	editor.SetStatusMessage("HELP: Ctrl-S = save | Ctrl-Q = quit")
	editor.EnterReaderLoop()
}
