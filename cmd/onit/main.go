package main

import (
	"time"

	"github.com/chrispritchard/onit/internal/terminal"
)

var characters map[rune][]string = map[rune][]string{
	'0': {
		" ## ",
		"#  #",
		"#  #",
		"#  #",
		" ## ",
	},
	'1': {
		" #  ",
		"##  ",
		" #  ",
		" #  ",
		"### ",
	},
	'2': {
		"####",
		"#  #",
		" ## ",
		"#   ",
		"####",
	},
	'3': {
		"####",
		"   #",
		" ###",
		"   #",
		"####",
	},
	'4': {
		" ## ",
		"# # ",
		"####",
		"  # ",
		"  # ",
	},
	'5': {
		"####",
		"#   ",
		"### ",
		"   #",
		"### ",
	},
	'6': {
		" ###",
		"#   ",
		"### ",
		"#  #",
		" ## ",
	},
	'7': {
		"####",
		"   #",
		"  # ",
		" #  ",
		"#   ",
	},
	'8': {
		"####",
		"#  #",
		"####",
		"#  #",
		"####",
	},
	'9': {
		" ## ",
		"#  #",
		"####",
		"  # ",
		" #  ",
	},
	' ': {
		"    ",
		"    ",
		"    ",
		"    ",
		"    ",
	},
	':': {
		"    ",
		" ## ",
		"    ",
		" ## ",
		"    ",
	},
	'A': {
		" ## ",
		"#  #",
		"####",
		"#  #",
		"#  #",
	},
	'P': {
		"####",
		"#  #",
		"####",
		"#   ",
		"#   ",
	},
	'M': {
		" ## ",
		"# ##",
		"# ##",
		"#  #",
		"#  #",
	},
}

func main() {

	ba := terminal.BufferedArea{}
	defer ba.Close()

	show_colon := true
	for {
		now := time.Now()
		time_string := now.Format("03:04 PM")

		to_display := make([]string, 5)
		for i, c := range time_string {
			if (i == 0 && c == '0') || (i == 2 && !show_colon) {
				c = ' '
			}
			pat, exists := characters[c]
			if !exists {
				panic("could not find pattern for rune: " + string(c))
			}
			for j, s := range pat {
				to_display[j] += s + " "
			}
		}

		ba.Update(to_display)

		time.Sleep(time.Second)
		show_colon = !show_colon
	}
}
