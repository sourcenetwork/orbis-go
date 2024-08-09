package keystore

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

// PromptFunc is a function used to prompt the user for a password.
type PromptFunc func(string) ([]byte, error)

func TerminalPrompt(prompt string) ([]byte, error) {
	fmt.Fprintf(os.Stderr, "%s: ", prompt)
	b, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return nil, err
	}
	fmt.Fprintf(os.Stderr, "\n") // line break for formatting consistency
	return b, nil
}

func FixedStringPrompt(value string) PromptFunc {
	return func(_ string) ([]byte, error) {
		return []byte(value), nil
	}
}
