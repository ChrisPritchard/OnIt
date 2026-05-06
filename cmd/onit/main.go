package main

import (
	_ "embed"
	"fmt"
	"slices"
	"time"

	"github.com/chrispritchard/onit/internal/bigtext"
	"github.com/chrispritchard/onit/internal/terminal"
	"github.com/chrispritchard/onit/internal/zonetimes"
)

func main() {
	ba := terminal.BufferedArea{}
	defer ba.Close()

	zones := zonetimes.GetZones()

	large_display := func(time time.Time, loc *time.Location, tick bool) []string {
		zone, time_string := zonetimes.GetDisplayTime(time, loc, tick)
		big_string := bigtext.GetBigString(time_string)
		return append([]string{fmt.Sprintf("%s time:", zone)}, big_string...)
	}

	newline := []string{""}
	tick := true
	for {
		now := time.Now()
		to_display := slices.Concat(
			newline,
			large_display(now, zones.NZT, tick),
			newline,
			large_display(now, zones.HKT, tick),
			newline,
			[]string{
				fmt.Sprintf("Epoch Time: %d", now.Unix()),
				"",
				"UTC\t\t" + zonetimes.GetTimeWithLoc(now, zones.UTC, tick) + "\t\t\t\t" + "US East\t\t" + zonetimes.GetTimeWithLoc(now, zones.USE, tick),
				"Central Europe\t" + zonetimes.GetTimeWithLoc(now, zones.CET, tick) + "\t\t\t\t" + "US West\t\t" + zonetimes.GetTimeWithLoc(now, zones.USP, tick),
				"Tokyo\t\t" + zonetimes.GetTimeWithLoc(now, zones.JPN, tick) + "\t\t\t\t" + "Brazil SP\t" + zonetimes.GetTimeWithLoc(now, zones.BRZ, tick),
				"",
				now.Format("Monday, 02 Jan 2006"),
			})

		ba.Update(to_display)

		time.Sleep(time.Second)
		tick = !tick
	}
}
