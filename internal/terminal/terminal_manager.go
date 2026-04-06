package terminal

import (
	"fmt"
)

const (
	escape         = "\x1b"
	cursorNextLine = escape + "[1E"
	cursorPrevLine = escape + "[1F"
	clearToEnd     = escape + "[0K"
	clearWholeLine = escape + "[2K"
	cursorHide     = escape + "[?25l"
	cursorShow     = escape + "[?25h"
)

type BufferedArea struct {
	last          []string
	cursor_hidden bool
}

func (ba *BufferedArea) Update(lines []string) {

	if !ba.cursor_hidden {
		fmt.Print(cursorHide)
		ba.cursor_hidden = true
	}

	if ba.last == nil {
		ba.last = lines
		for _, l := range lines {
			fmt.Println(l)
		}
		return
	}

	for range len(ba.last) {
		fmt.Print(cursorPrevLine)
	}

	for i, l := range lines {
		if i >= len(ba.last) {
			fmt.Println(l)
		} else if l == ba.last[i] {
			fmt.Printf(cursorNextLine)
		} else {
			fmt.Println(l + clearToEnd)
		}
	}

	if len(lines) < len(ba.last) {
		for range len(ba.last) - len(lines) {
			fmt.Println(clearWholeLine)
		}
	}

	ba.last = lines
}

func (ba *BufferedArea) Close() {
	fmt.Print(cursorShow)
}
