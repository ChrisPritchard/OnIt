package bigtext

import (
	_ "embed"
	"strconv"
	"strings"
)

//go:embed font_string.txt
var font string

type font_info struct {
	lines     []string
	positions []int
	lengths   []int
}

func (fi *font_info) segment(c rune, line int) string {
	var char_set = "0123456789:APM"

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

	return fi.lines[line][fi.positions[index] : fi.positions[index]+fi.lengths[index]]
}

var fi font_info

func init() {
	fi = get_font_info()
}

func get_font_info() font_info {
	font_lines := strings.Split(font, "\n")
	char_lengths := []int{}
	for s := range strings.SplitSeq(font_lines[1], " ") {
		n, _ := strconv.Atoi(s)
		char_lengths = append(char_lengths, n)
	}

	positions := []int{}
	for i := range char_lengths {
		sum := 0
		for j := range i {
			sum += char_lengths[j]
		}
		positions = append(positions, sum)
	}

	return font_info{
		lines:     font_lines[2:], // remove header and length
		positions: positions,
		lengths:   char_lengths,
	}
}

func GetBigString(to_render string) []string {

	to_display := make([]string, len(fi.lines))
	for i, c := range to_render {
		if c == ' ' || (i == 0 && c == '0') {
			for j := range to_display {
				to_display[j] += "     "
			}
		} else {
			for j := range to_display {
				to_display[j] += fi.segment(c, j)
			}
		}
	}
	return to_display
}
