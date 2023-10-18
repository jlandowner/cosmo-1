package cli

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mattn/go-isatty"
	"golang.org/x/term"
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

func AskInput(prompt string, silent bool) (string, error) {
	fmt.Print(prompt)

	var input []byte
	var err error
	if silent {
		input, err = term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
	} else {
		r := bufio.NewReader(os.Stdin)
		input, err = r.ReadBytes('\n')
	}
	if err != nil {
		return "", fmt.Errorf("failed to read input : %w", err)
	}
	return string(input), nil
}
