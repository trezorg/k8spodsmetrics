package common

import (
	"io"
	"os"
	"strings"

	screen "github.com/aditya43/clear-shell-screen-golang"
	"golang.org/x/term"
)

func WriteStringLine(text string) {
	_, _ = os.Stdout.WriteString(text + "\n")
}

// getTerminalHeight returns the terminal height (rows), or default if detection fails
func getTerminalHeight() int {
	const defaultHeight = 24
	fd := os.Stdout.Fd()
	//nolint:gosec // fd is a file descriptor, always a small positive integer
	if _, height, err := term.GetSize(int(fd)); err == nil && height > 0 {
		return height
	}
	return defaultHeight
}

// truncateLines limits output to the first N lines where N = terminal height
func truncateLines(output string, maxLines int) string {
	if output == "" || maxLines <= 0 {
		return ""
	}
	lines := strings.Split(output, "\n")
	if len(lines) <= maxLines {
		return output
	}
	truncated := lines[:maxLines]
	return strings.Join(truncated, "\n")
}

func WrapScreenSuccess[T any](writer func(io.Writer, T)) func(T) {
	return func(value T) {
		var outputBuilder strings.Builder
		writer(&outputBuilder, value)
		truncated := truncateLines(outputBuilder.String(), getTerminalHeight())

		screen.Clear()
		screen.MoveTopLeft()
		_, _ = os.Stdout.WriteString(truncated)
		screen.MoveTopLeft()
	}
}

func WrapScreenError(writer func(error)) func(error) {
	return func(err error) {
		screen.Clear()
		screen.MoveTopLeft()
		writer(err)
		screen.MoveTopLeft()
	}
}
