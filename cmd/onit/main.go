package main

import (
	_ "embed"
	"fmt"
	"log"
	"slices"
	"strings"
	"time"

	"github.com/chrispritchard/onit/internal/bigtext"
	"github.com/chrispritchard/onit/internal/terminal"
)

func main() {
	ba := terminal.BufferedArea{}
	defer ba.Close()

	zones := get_zones()

	large_display := func(time time.Time, loc *time.Location, tick bool) []string {
		zone, time_string := get_display_time(time, loc, tick)
		big_string := bigtext.GetBigString(time_string)
		return append([]string{fmt.Sprintf("%s time:", zone)}, big_string...)
	}

	newline := []string{""}
	tick := true
	for {
		now := time.Now()
		to_display := slices.Concat(
			newline,
			large_display(now, zones.nzt, tick),
			newline,
			large_display(now, zones.hkt, tick),
			newline,
			[]string{
				fmt.Sprintf("Epoch Time: %d", now.Unix()),
				"",
				"UTC\t\t" + get_combined_time(now, zones.utc, tick) + "\t\t\t\t" + "US East\t\t" + get_combined_time(now, zones.use, tick),
				"Central Europe\t" + get_combined_time(now, zones.cet, tick) + "\t\t\t\t" + "US West\t\t" + get_combined_time(now, zones.usp, tick),
				"",
				now.Format("Monday, 02 Jan 2006"),
			})

		ba.Update(to_display)

		time.Sleep(time.Second)
		tick = !tick
	}
}

func get_display_time(time time.Time, loc *time.Location, tick bool) (zone string, display string) {
	rel := time.In(loc)
	zone, _ = rel.Zone()
	display = rel.Format("03:04 PM")
	if !tick {
		display = strings.Replace(display, ":", " ", 1)
	}
	return
}

func get_combined_time(time time.Time, loc *time.Location, tick bool) string {
	name, display := get_display_time(time, loc, tick)
	return fmt.Sprintf("%s %s", display, name)
}

type zones struct {
	nzt *time.Location
	hkt *time.Location
	cet *time.Location
	utc *time.Location
	usp *time.Location
	use *time.Location
}

func get_zones() zones {
	load_loc := func(iana string) *time.Location {
		loc, err := time.LoadLocation(iana)
		if err != nil {
			log.Fatal(err)
		}
		return loc
	}

	return zones{
		nzt: load_loc("Pacific/Auckland"),
		hkt: load_loc("Asia/Hong_Kong"),
		cet: load_loc("Europe/Berlin"),
		utc: load_loc("Europe/London"),
		usp: load_loc("America/Los_Angeles"),
		use: load_loc("America/New_York"),
	}
}
