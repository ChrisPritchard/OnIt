package main

import (
	_ "embed"
	"fmt"
	"log"
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

	show_colon := true
	for {
		message := api.GetDisplayState()

		current := time.Now()
		to_display := append(append(
			time_display(current, show_colon),
			time_display(current.In(hkt), show_colon)...),
			message...)
		ba.Update(to_display)

		time.Sleep(time.Second)
		show_colon = !show_colon
	}
}
