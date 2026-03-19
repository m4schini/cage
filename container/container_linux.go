package container

import (
	"os"

	"golang.org/x/term"
)

func PrepareTTY() (cleanup func(), err error) {
	// Set terminal to raw mode for proper TTY interaction
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		return func() {}, err
	}
	return func() {
		term.Restore(int(os.Stdin.Fd()), oldState)
	}, nil
}
