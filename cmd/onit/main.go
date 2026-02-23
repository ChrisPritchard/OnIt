package main

import (
	_ "embed"
	"strconv"
	"strings"
	"time"

	"github.com/chrispritchard/onit/internal/terminal"
)

var char_set = "0123456789:APM"

//go:embed font_string.txt
var font string

func main() {

	ba := terminal.BufferedArea{}
	defer ba.Close()

	font_lines := strings.Split(font, "\n")
	char_lengths := []int{}
	for s := range strings.SplitSeq(font_lines[1], " ") {
		n, _ := strconv.Atoi(s)
		char_lengths = append(char_lengths, n)
	}

	font_lines = font_lines[2:] // remove header and length

	positions := []int{}
	for i := range char_lengths {
		sum := 0
		for j := range i {
			sum += char_lengths[j]
		}
		positions = append(positions, sum)
	}

	show_colon := true
	for {
		now := time.Now()
		time_string := now.Format("03:04 PM")

		to_display := make([]string, len(font_lines))
		for i, c := range time_string {
			if c == ' ' || (i == 0 && c == '0') || (i == 2 && !show_colon) {
				for j := range to_display {
					to_display[j] += "     "
				}
			} else {
				index := -1
				for k, o := range char_set {
					if o == c {
						index = k
						break
					}
				}
				if index == -1 {
					panic("could not find pattern for rune: " + string(c))
				}

				for j := range to_display {
					to_display[j] += font_lines[j][positions[index] : positions[index]+char_lengths[index]]
				}
			}
		}

		ba.Update(to_display)

		time.Sleep(time.Second)
		show_colon = !show_colon
	}
}
