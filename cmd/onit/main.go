package main

import (
	_ "embed"
	"fmt"
	"log"
	"slices"
	"strings"
	"time"

	"github.com/chrispritchard/onit/internal/api"
	"github.com/chrispritchard/onit/internal/bigtext"
	"github.com/chrispritchard/onit/internal/terminal"
)

func main() {
	ba := terminal.BufferedArea{}
	defer ba.Close()

	go api.StartApiServer()()
	render_display(&ba)
}

func render_display(ba *terminal.BufferedArea) {
	hkt, err := time.LoadLocation("Asia/Hong_Kong")
	if err != nil {
		log.Fatal(err)
	}

	time_display := func(time time.Time, show_colon bool) []string {
		zone, _ := time.Zone()
		time_string := time.Format("03:04 PM")
		if !show_colon {
			time_string = strings.Replace(time_string, ":", " ", 1)
		}
		big_string := bigtext.GetBigString(time_string)
		return append([]string{fmt.Sprintf("%s time:", zone)}, big_string...)
	}

	tick := true
	for {
		current := time.Now()
		message := api.GetDisplayState(current, tick)

		to_display := slices.Concat(
			[]string{fmt.Sprintf("Epoch time: %d", current.Unix())},
			time_display(current, tick),
			time_display(current.In(hkt), tick),
			message)

		ba.Update(to_display)

		time.Sleep(time.Second)
		tick = !tick
	}
}
