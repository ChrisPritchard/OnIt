package zonetimes

import (
	"fmt"
	"log"
	"strings"
	"time"
)

func GetDisplayTime(time time.Time, loc *time.Location, tick bool) (zone string, display string) {
	rel := time.In(loc)
	zone, _ = rel.Zone()
	display = rel.Format("03:04 PM")
	if !tick {
		display = strings.Replace(display, ":", " ", 1)
	}
	return
}

func GetTimeWithLoc(time time.Time, loc *time.Location, tick bool) string {
	name, display := GetDisplayTime(time, loc, tick)
	return fmt.Sprintf("%s %s", display, name)
}

type Zones struct {
	NZT *time.Location
	HKT *time.Location
	CET *time.Location
	UTC *time.Location
	USP *time.Location
	USE *time.Location
	JPN *time.Location
	BRZ *time.Location
}

func GetZones() Zones {
	load_loc := func(iana string) *time.Location {
		loc, err := time.LoadLocation(iana)
		if err != nil {
			log.Fatal(err)
		}
		return loc
	}

	return Zones{
		NZT: load_loc("Pacific/Auckland"),
		HKT: load_loc("Asia/Hong_Kong"),
		CET: load_loc("Europe/Berlin"),
		UTC: load_loc("Europe/London"),
		USP: load_loc("America/Los_Angeles"),
		USE: load_loc("America/New_York"),
		JPN: load_loc("Asia/Tokyo"),
		BRZ: load_loc("America/Sao_Paulo"),
	}
}
