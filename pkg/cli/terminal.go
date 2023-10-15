package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mattn/go-isatty"
	"golang.org/x/crypto/ssh/terminal"
)

func ReadFromPipedStdin() (string, error) {
	if isatty.IsTerminal(os.Stdin.Fd()) {
		return "", fmt.Errorf("not terminal")
	}
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		return "", fmt.Errorf("failed to read input from stdin: %w", err)
	}
	return strings.Replace(string(input), "\n", "", 1), nil
}

func AskInput(prompt string) (string, error) {
	fmt.Print(prompt)
	password, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", fmt.Errorf("failed to read input : %w", err)
	}
	fmt.Println()
	return string(password), nil
}
