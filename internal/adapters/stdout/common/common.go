package common

import (
	"os"

	screen "github.com/aditya43/clear-shell-screen-golang"
)

func WriteStringLine(text string) {
	_, _ = os.Stdout.WriteString(text + "\n")
}

func WrapScreenSuccess[T any](writer func(T)) func(T) {
	return func(value T) {
		screen.Clear()
		screen.MoveTopLeft()
		writer(value)
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
