package main

import (
	_ "embed"
	"strconv"
	"strings"
	"time"

	"github.com/chrispritchard/onit/internal/terminal"
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

func main() {

	ba := terminal.BufferedArea{}
	defer ba.Close()

	fi := get_font_info()

	show_colon := true
	for {
		nz_now := time.Now()
		to_display_nz := get_time_string(nz_now, fi, show_colon)
		hk_now := nz_now.Add(-5 * time.Hour)
		to_display_hk := get_time_string(hk_now, fi, show_colon)

		to_display := append(
			append([]string{"NZDT time:"}, to_display_nz...),
			append([]string{"", "HKT time:"}, to_display_hk...)...)
		ba.Update(to_display)

		time.Sleep(time.Second)
		show_colon = !show_colon
	}
}

func get_time_string(to_render time.Time, font font_info, show_colon bool) []string {
	time_string := to_render.Format("03:04 PM")

	to_display := make([]string, len(font.lines))
	for i, c := range time_string {
		if c == ' ' || (i == 0 && c == '0') || (i == 2 && !show_colon) {
			for j := range to_display {
				to_display[j] += "     "
			}
		} else {
			for j := range to_display {
				to_display[j] += font.segment(c, j)
			}
		}
	}
	return to_display
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
